package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"webook/internal/service"
	"webook/internal/service/oauth2/wechat"
	ijwt "webook/internal/web/jwt"
)

type OAuth2WechatHandler struct {
	ijwt.Handler
	oauthsvc        wechat.OAuth2Service
	userService     service.UserService
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.OAuth2Service, userSvc service.UserService, hdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		oauthsvc:        svc,
		userService:     userSvc,
		key:             []byte("RrRqvf7sVUhBwm0hTl9Umu1vu1unNkp7"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
	}
}
func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}
func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	state := uuid.New()
	val, err := o.oauthsvc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "服务器异常",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})

}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 4,
		})
		return
	}
	code := ctx.Query("code")
	//state:=ctx.Query("state")
	wechatInfo, err := o.oauthsvc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		return
	}
	u, err := o.userService.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return

}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {

		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr, 600, "oauth/wechat/callback", "", false, true)
	return nil
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return fmt.Errorf("无法获得cookie %w", err)

	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("解析token失败 %w", err)
	}
	if state != sc.State {
		//state不匹配，有攻击者
		return fmt.Errorf("state不匹配")
	}
	return nil
}
