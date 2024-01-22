package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	ijwt "webook/internal/web/jwt"

	//"github.com/jinzhu/now"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const (
	emailRegexPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, hdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
		Handler:        hdl,
	}
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", c.SignUp)
	//ug.POST("/login", c.Login)
	ug.POST("/login", c.LoginJWT)
	ug.POST("/logout", c.LogoutJWT)
	ug.POST("/edit", c.Edit)
	ug.GET("/profile", c.ProfileJWT)
	ug.GET("/refresh_token", c.RefreshToken)
	//手机验证码登录相关功能
	ug.POST("/login_sms/code/send", c.SendSMSLoginCode)
	ug.POST("/login_sms", c.LoginSMS)
}

func (c *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := c.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		zap.L().Error("手机验证码验证失败",
			//在生产环境绝对不能打
			//开发环境你可以随便打
			zap.String("phone", req.Phone),
			zap.Error((err)))
		return

	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	u, err := c.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = c.SetLoginToken(ctx, u.Id)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, Result{

		Msg: "登录成功",
	})

}

func (c *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号",
		})
		return
	}
	err := c.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{

			Msg: "发送成功",
		})

	case service.ErrCodeSendTooMany:
		zap.L().Warn("频繁发送验证码")
		ctx.JSON(http.StatusOK, Result{

			Msg: "短信发送太频繁，请稍候再试",
		})

	default:
		ctx.JSON(http.StatusOK, Result{

			Msg: "系统错误",
		})

	}
	return

}
func (c *UserHandler) SignUp(ctx *gin.Context) {
	//优先使用内部类
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	//校验邮箱地址
	isEmail, err := c.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	//两次密码输入校验
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	//校验密码
	isPassword, err := c.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "非法密码格式,密码必须包含字母、数字、特殊字符，并且不少于8位")
		return
	}

	err = c.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:

		ctx.String(http.StatusOK, "hello,你在注册")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}

}
func (c *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			//十五分钟
			MaxAge: 30,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "sesion存储错误")
			return
		}

		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}

}

func (c *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := c.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// uc := UsersClaims{
		// 	Uid:       u.Id,
		// 	UserAgent: ctx.GetHeader("User-Agent"),
		// 	RegisteredClaims: jwt.RegisteredClaims{
		// 		//30分钟过期
		// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		// 	},
		// }
		// token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		// tokenStr, err := token.SignedString(JWTKey)
		// if err != nil {

		// 	ctx.String(http.StatusOK, "tokensStr系统错误")

		// }
		// ctx.Header("x-jwt-token", tokenStr)
		err = c.SetLoginToken(ctx, u.Id)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}

		ctx.String(http.StatusOK, "登录成功")
		//sess := sessions.Default(ctx)
		//sess.Set("userId", u.Id)
		//sess.Options(sessions.Options{
		//	//十五分钟
		//	MaxAge: 30,
		//})
		//err = sess.Save()
		//if err != nil {
		//	ctx.String(http.StatusOK, "sesion存储错误")
		//	return
		//}

	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}
}

func (c *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Nickname == "" {
		ctx.String(http.StatusOK, "昵称不能为空")
	}
	if len(req.AboutMe) > 1024 {
		ctx.String(http.StatusOK, "关于我过长")
	}

	sess := sessions.Default(ctx)

	uc := sess.Get("userId")

	birthday, err := time.Parse(time.DateOnly, req.Birthday)

	if err != nil {
		ctx.String(http.StatusOK, "日期格式不对")
		return
	}

	err = c.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.(int64),
		Nickname: req.Nickname,
		Birthday: birthday.String(),
		AboutMe:  req.AboutMe,
	})

	switch err {
	case nil:

		ctx.String(http.StatusOK, "信息编辑完成")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}

}

func (c *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email    string
		Nickname string
		Birthday string
		AboutMe  string
	}

	uc := ctx.MustGet("user").(ijwt.UsersClaims)
	u, err := c.svc.Profile(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday,
	})
}
func (c *UserHandler) Profile(ctx *gin.Context) {
	//us :=ctx.MustGet("user").(UsersClaims)

	//sess := sessions.Default(ctx)
	//uc := sess.Get("userId")
	//
	//u, err := c.svc.FindById(ctx, uc.(int64))
	//if err != nil {
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}

	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutme"`
		Birthday string `json:"birthday"`
	}
	//
	//ctx.JSON(http.StatusOK, User{
	//	Nickname: u.Nickname,
	//	Email:    u.Email,
	//	AboutMe:  u.AboutMe,
	//	Birthday: u.Brithday,
	//})
	//type Profile struct {
	//	Email string
	//}
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	u, err := c.svc.Profile(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	ctx.String(http.StatusOK, "这是Profile")
	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday,
	})

}

func (c *UserHandler) RefreshToken(ctx *gin.Context) {
	//约定，前端在Authorization里面带上这个refresh_token
	tokenStr := c.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil

	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = c.CheckSession(ctx, rc.Ssid)
	//这种情况比较严格，redis崩溃后无法为用户提供服务
	if err != nil {
		//token无效或者redis有问题
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = c.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (c *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := c.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return

	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}
