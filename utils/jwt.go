package utils

import (
	"fmt"
	"time"

	"iris-cn-sample-project/config"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT 声明结构体
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT 令牌
func GenerateJWT(userID uint, username, role string) (string, time.Time, error) {
	cfg := config.GetConfig()
	
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpirationTime) * time.Second)
	
	// 创建 JWT 声明
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.JWT.Issuer,
			Subject:   fmt.Sprintf("user:%d", userID),
			ID:        fmt.Sprintf("%d_%d", userID, time.Now().Unix()),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名令牌
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成JWT令牌失败: %v", err)
	}

	return tokenString, expirationTime, nil
}

// ValidateJWT 验证 JWT 令牌
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	cfg := config.GetConfig()
	
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析JWT令牌失败: %v", err)
	}

	// 验证令牌并提取声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// 检查令牌是否过期
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, fmt.Errorf("令牌已过期")
		}
		
		// 检查令牌是否在生效时间之前
		if time.Now().Before(claims.NotBefore.Time) {
			return nil, fmt.Errorf("令牌尚未生效")
		}
		
		return claims, nil
	}

	return nil, fmt.Errorf("无效的JWT令牌")
}

// ParseJWTWithoutValidation 解析JWT令牌但不验证签名和过期时间
func ParseJWTWithoutValidation(tokenString string) (*JWTClaims, error) {
	// 解析令牌但不验证签名
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return nil, fmt.Errorf("解析JWT令牌失败: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("无法解析JWT声明")
}

// RefreshJWT 刷新JWT令牌
func RefreshJWT(tokenString string) (string, time.Time, error) {
	// 首先解析原始令牌（不验证过期时间）
	claims, err := ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("解析原始令牌失败: %v", err)
	}

	// 生成新的令牌
	newToken, expirationTime, err := GenerateJWT(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成新令牌失败: %v", err)
	}

	return newToken, expirationTime, nil
}

// GetJWTFromHeader 从请求头中提取JWT令牌
func GetJWTFromHeader(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	
	if authHeader == "" {
		return "", fmt.Errorf("缺少Authorization头")
	}
	
	if len(authHeader) <= len(bearerPrefix) {
		return "", fmt.Errorf("Authorization头格式错误")
	}
	
	if authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("Authorization头必须以Bearer开头")
	}
	
	return authHeader[len(bearerPrefix):], nil
}

// IsJWTExpired 检查JWT令牌是否过期
func IsJWTExpired(tokenString string) bool {
	claims, err := ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return true
	}
	
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetJWTExpiration 获取JWT令牌过期时间
func GetJWTExpiration(tokenString string) (time.Time, error) {
	claims, err := ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return time.Time{}, err
	}
	
	return claims.ExpiresAt.Time, nil
}

// GetJWTIssuer 获取JWT令牌发行者
func GetJWTIssuer(tokenString string) (string, error) {
	claims, err := ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return "", err
	}
	
	return claims.Issuer, nil
}

// GetJWTSubject 获取JWT令牌主题
func GetJWTSubject(tokenString string) (string, error) {
	claims, err := ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return "", err
	}
	
	return claims.Subject, nil
}

// GetJWTID 获取JWT令牌ID
func GetJWTID(tokenString string) (string, error) {
	claims, err := ParseJWTWithoutValidation(tokenString)
	if err != nil {
		return "", err
	}
	
	return claims.ID, nil
}

// ValidateJWTWithClaims 验证JWT令牌并返回自定义声明
func ValidateJWTWithClaims(tokenString string, customClaims interface{}) error {
	cfg := config.GetConfig()
	
	token, err := jwt.ParseWithClaims(tokenString, customClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return fmt.Errorf("解析JWT令牌失败: %v", err)
	}

	if !token.Valid {
		return fmt.Errorf("无效的JWT令牌")
	}

	return nil
}

// CreateJWTWithCustomClaims 创建带自定义声明的JWT令牌
func CreateJWTWithCustomClaims(claims jwt.Claims) (string, time.Time, error) {
	cfg := config.GetConfig()
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("生成JWT令牌失败: %v", err)
	}

	// 尝试从声明中获取过期时间
	var expirationTime time.Time
	if expClaims, ok := claims.(jwt.Claims); ok {
		if exp := expClaims.GetExpirationTime(); exp != nil {
			expirationTime = exp.Time
		}
	}

	return tokenString, expirationTime, nil
}

// DecodeJWTWithoutValidation 解码JWT令牌（不验证签名）
func DecodeJWTWithoutValidation(tokenString string) (map[string]interface{}, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("解码JWT令牌失败: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		result := make(map[string]interface{})
		for key, value := range claims {
			result[key] = value
		}
		return result, nil
	}

	return nil, fmt.Errorf("无法解析JWT声明")
}

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	// 注意：实际项目中应该使用 bcrypt 或其他安全的密码哈希算法
	// 这里为了简化示例，使用简单的哈希
	return password + "_hashed", nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) bool {
	// 注意：实际项目中应该使用 bcrypt 或其他安全的密码验证算法
	// 这里为了简化示例，使用简单的比较
	return password+"_hashed" == hashedPassword
}