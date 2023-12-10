package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	//"github.com/jinzhu/now"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const (
	emailRegexPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct {
	emailRexExp     *regexp.Regexp
	passworedRexExp *regexp.Regexp
	svc             *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:     regexp.MustCompile(emailRegexPattern, regexp.None),
		passworedRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:             svc,
	}
}

func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login", c.Login)
	ug.POST("/edit", c.Edit)
	ug.GET("/profile", c.Profile)
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
	isPassword, err := c.passworedRexExp.MatchString(req.Password)
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
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
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
			MaxAge: 900,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}

		ctx.String(http.StatusOK, "登录成功")
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
		Brithday: birthday.String(),
		AboutMe:  req.AboutMe,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "信息编辑完成")
	default:
		ctx.String(http.StatusOK, "系统错误")

	}

}
func (c *UserHandler) Profile(ctx *gin.Context) {

	sess := sessions.Default(ctx)
	uc := sess.Get("userId")

	u, err := c.svc.FindById(ctx, uc.(int64))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutme"`
		Birthday string `json:"birthday"`
	}

	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Brithday,
	})

}
