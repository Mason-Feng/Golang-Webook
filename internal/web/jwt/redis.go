package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var JWTKey = []byte("RrRqvf7sVUhBwm0hTl9Umu1vu1unNkp6")
var RCJWTKey = []byte("RrRqvf7sVUhBwm0hTl9Umu1vu1unNkp7")

type RedisJWTHandler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod

	rcExpiration time.Duration
}

func NewRedisJWTHandler(client redis.Cmdable) *RedisJWTHandler {
	return &RedisJWTHandler{
		client:        client,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Hour * 24 * 7,
	}

}
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		//没登录，没有token，Authorization这个头部都没有

		return authCode
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {

		return ""
	}
	return segs[1]
}

var _ Handler = &RedisJWTHandler{}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := h.client.Exists(ctx, fmt.Sprintf("users:ssid%s", ssid)).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token无效")
	}
	return nil
}

type UsersClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}

func (c *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := c.setRefreshToken(ctx, uid, ssid)
	if err != nil {

		return err
	}
	c.SetJWTToken(ctx, uid, ssid)
	return nil
}
func (c *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	//这里表示每次登录都会自动延期token时间
	err := c.setRefreshToken(ctx, uid, ssid)
	if err != nil {

		return err
	}
	uc := UsersClaims{
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			//30分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(c.signingMethod, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {

		return err

	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

func (c *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {

	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.rcExpiration)),
		},
	}
	token := jwt.NewWithClaims(c.signingMethod, rc)
	tokenStr, err := token.SignedString(RCJWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (c *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc := ctx.MustGet("user").(UsersClaims)
	return c.client.Set(ctx, fmt.Sprintf("users:ssid%s", uc.Ssid), "", c.rcExpiration).Err()
}
