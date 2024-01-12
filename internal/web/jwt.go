package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var JWTKey = []byte("RrRqvf7sVUhBwm0hTl9Umu1vu1unNkp6")

type jwtHandler struct {
}

type UsersClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (c *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UsersClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			//30分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {

		ctx.String(http.StatusOK, "tokensStr系统错误")

	}
	ctx.Header("x-jwt-token", tokenStr)
}
