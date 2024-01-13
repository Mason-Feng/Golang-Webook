package middleware

import (
	"net/http"

	ijwt "webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(hdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}

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

		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UsersClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})

		if err != nil {
			//token不对，token是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			//token解析出来了，但是token可能是非法的，或者过期的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//if uc.UserAgent != ctx.GetHeader("User-Agent") {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		//expireTime := uc.ExpiresAt
		////剩余过期时间<50s要刷新
		//if expireTime.Sub(time.Now()) < time.Second*50 {
		//	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
		//	tokenStr, err = token.SignedString(web.JWTKey)
		//	ctx.Header("x-jwt-token", tokenStr)
		//	if err != nil {
		//		//这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
		//		log.Panicln(err)
		//	}
		//}
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {

			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		////这种情况比较严格，redis崩溃后无法为用户提供服务
		//if err != nil || cnt > 0 {
		//	//token无效或者redis有问题
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		////这种可以兼容Redis异常的情况，即使redis崩溃后，依然能为用户提供有损对的服务
		////但是这种情况要做好监控，监控有没有error
		//if cnt>0{
		//	//token无效或者redis有问题
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		ctx.Set("user", uc)

	}
}
