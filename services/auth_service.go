package services

import (
	"errors"
	"fmt"
	"time"

	"iris-cn-sample-project/config"
	"iris-cn-sample-project/models"
	"iris-cn-sample-project/utils"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT 声明结构体
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokenPair 生成令牌对（访问令牌 + 刷新令牌）
func GenerateTokenPair(user *models.User) (accessToken, refreshToken string, expiresAt time.Time, err error) {
	// 生成访问令牌
	accessToken, expiresAt, err = GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("生成访问令牌失败: %v", err)
	}

	// 生成刷新令牌（有效期更长）
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour) // 7天
	refreshClaims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.GetConfig().JWT.Issuer,
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshTokenObj.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("生成刷新令牌失败: %v", err)
	}

	return accessToken, refreshTokenString, expiresAt, nil
}

// ValidateToken 验证令牌并返回用户信息
func ValidateToken(tokenString string) (*models.User, error) {
	// 验证 JWT 令牌
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return nil, fmt.Errorf("令牌验证失败: %v", err)
	}

	// 根据令牌中的用户ID获取用户信息
	user, err := GetUserByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 将用户信息转换为完整的用户模型
	fullUser := &models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return fullUser, nil
}

// RefreshAccessToken 使用刷新令牌获取新的访问令牌
func RefreshAccessToken(refreshTokenString string) (string, time.Time, error) {
	// 解析刷新令牌（不验证过期时间）
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().JWT.Secret), nil
	})

	if err != nil {
		return "", time.Time{}, fmt.Errorf("解析刷新令牌失败: %v", err)
	}

	// 检查令牌是否有效
	if claims, ok := refreshToken.Claims.(*JWTClaims); ok && refreshToken.Valid {
		// 验证用户是否仍然存在且有效
		user, err := GetUserByID(claims.UserID)
		if err != nil {
			return "", time.Time{}, fmt.Errorf("用户不存在或已被禁用: %v", err)
		}

		// 生成新的访问令牌
		newToken, expiresAt, err := GenerateJWT(user.ID, user.Username, user.Role)
		if err != nil {
			return "", time.Time{}, fmt.Errorf("生成新令牌失败: %v", err)
		}

		return newToken, expiresAt, nil
	}

	return "", time.Time{}, errors.New("无效的刷新令牌")
}

// InvalidateToken 使令牌失效（在实际应用中，可以将令牌加入黑名单）
func InvalidateToken(tokenString string) error {
	// 这里简化处理，实际应用中可以实现令牌黑名单机制
	// 例如将令牌的 JTI（JWT ID）存储在 Redis 或数据库中
	// 并在验证时检查令牌是否在黑名单中
	
	// 模拟将令牌加入黑名单
	_ = tokenString
	
	return nil
}

// GetTokenExpiration 获取令牌过期时间
func GetTokenExpiration(tokenString string) (time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JWT.Secret), nil
	})

	if err != nil {
		return time.Time{}, fmt.Errorf("解析令牌失败: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.ExpiresAt.Time, nil
	}

	return time.Time{}, errors.New("无效的令牌")
}

// IsTokenExpired 检查令牌是否过期
func IsTokenExpired(tokenString string) bool {
	expiration, err := GetTokenExpiration(tokenString)
	if err != nil {
		return true
	}
	return time.Now().After(expiration)
}

// GetTokenInfo 获取令牌信息（不验证签名）
func GetTokenInfo(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 不验证签名，只解析
		return []byte("dummy"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, errors.New("无法解析令牌声明")
}

// ValidateTokenForUser 验证令牌是否属于指定用户
func ValidateTokenForUser(tokenString string, userID uint) error {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("令牌验证失败: %v", err)
	}

	if claims.UserID != userID {
		return errors.New("令牌不属于指定用户")
	}

	return nil
}

// ValidateTokenRole 验证令牌用户角色
func ValidateTokenRole(tokenString string, requiredRole string) error {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("令牌验证失败: %v", err)
	}

	if claims.Role != requiredRole && claims.Role != "admin" {
		return errors.New("用户权限不足")
	}

	return nil
}

// CreateSession 创建用户会话
func CreateSession(user *models.User) (map[string]interface{}, error) {
	// 生成令牌对
	accessToken, refreshToken, expiresAt, err := GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("创建会话失败: %v", err)
	}

	// 构建会话信息
	session := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_at":    expiresAt.Unix(),
		"expires_in":    int(time.Until(expiresAt).Seconds()),
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"status":   user.Status,
		},
	}

	return session, nil
}

// RefreshSession 刷新用户会话
func RefreshSession(refreshTokenString string) (map[string]interface{}, error) {
	// 使用刷新令牌获取新的访问令牌
	newAccessToken, expiresAt, err := RefreshAccessToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("刷新会话失败: %v", err)
	}

	// 构建新的会话信息
	session := map[string]interface{}{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_at":   expiresAt.Unix(),
		"expires_in":   int(time.Until(expiresAt).Seconds()),
	}

	return session, nil
}

// DestroySession 销毁用户会话
func DestroySession(tokenString string) error {
	// 将令牌加入黑名单
	if err := InvalidateToken(tokenString); err != nil {
		return fmt.Errorf("销毁会话失败: %v", err)
	}

	return nil
}