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

type RedisJWTHandler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod
	rcExpiration  time.Duration
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		client:        client,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Hour * 24 * 7,
	}
}

func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := h.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	// 面试可以用的降级策略
	// 在 Redis 没有崩溃的时候,就会严格执行 ssid 的校验，判定用户有没有主动退出登录
	// 而如果 Redis 崩溃,就不会严格执行 ssid 的校验
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token 无效")
	}
	return nil
}

// ExtractToken 根据约定, token 在 Authorization 头部
// Bearer XXXX
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return authCode
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

var _ Handler = &RedisJWTHandler{}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return h.SetJWTToken(ctx, uid, ssid)
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc := ctx.MustGet("user").(UserClaims)

	return h.client.Set(ctx, fmt.Sprintf("users:ssid:%s", uc.Ssid),
		"", h.rcExpiration).Err()
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:       uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(h.signingMethod, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置为七天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
	}
	refreshToken := jwt.NewWithClaims(h.signingMethod, rc)
	refreshTokenStr, err := refreshToken.SignedString(RCJWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

var (
	JWTKey   = []byte("sL0hO7kV4cT1aW1zN2gN7qG7jU1jL1iK")
	RCJWTKey = []byte("sL0hO7kV4cT1aW1zN2gN7qG7jU1jL1iK")
)

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}
