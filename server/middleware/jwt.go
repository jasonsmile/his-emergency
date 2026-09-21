package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"emergency-his/server/config"
	"emergency-his/server/response"
)

type Claims struct {
	UserID string `json:"user_id"`
	UserName string `json:"user_name"`
	DeptCode string `json:"dept_code"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, userName, deptCode string) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil || strings.TrimSpace(cfg.JWT.Secret) == "" { return "", errors.New("jwt secret is not configured") }
	hours := cfg.JWT.ExpireHours
	if hours <= 0 { hours = 24 }
	claims := Claims{UserID: userID, UserName: userName, DeptCode: deptCode, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hours) * time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now()), Issuer: "emergency-his",
	}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()
		if cfg == nil || strings.TrimSpace(cfg.JWT.Secret) == "" { response.Error(c, http.StatusInternalServerError, 500, "JWT未配置"); c.Abort(); return }
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") { response.Unauthorized(c, "Authorization格式错误"); c.Abort(); return }
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 { return nil, errors.New("unexpected signing method") }
			return []byte(cfg.JWT.Secret), nil
		}, jwt.WithIssuer("emergency-his"))
		if err != nil || !token.Valid { response.Unauthorized(c, "token无效或已过期"); c.Abort(); return }
		c.Set("user_id", claims.UserID); c.Set("user_name", claims.UserName); c.Set("dept_code", claims.DeptCode)
		c.Next()
	}
}

func GetCurrentUser(c *gin.Context) (userID, userName, deptCode string) {
	return c.GetString("user_id"), c.GetString("user_name"), c.GetString("dept_code")
}
