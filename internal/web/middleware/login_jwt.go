package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	//注册一下这个类型
	//gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			return
		}
		//根据约定，token在Authorization头部
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			//没登录，没有token，Authorization这个头部都没有
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			//没登录，Authorization中的内容是乱传的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UsersClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})

		if err != nil {
			//token不对，token是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || uc.ExpiresAt.Before(time.Now()) {
			//token解析出来了，但是token可能是非法的，或者过期的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//if uc.UserAgent != ctx.GetHeader("User-Agent") {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		expireTime := uc.ExpiresAt
		//剩余过期时间<50s要刷新
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				//这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
				log.Panicln(err)
			}
		}
		ctx.Set("user", uc)

	}
}
