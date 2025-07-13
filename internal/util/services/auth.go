package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type authKey string

type AuthService struct {
	AuthContextKey authKey
	secret         []byte
	tokenExpires   int
}

func NewAuthService(secret string, tokenExpires int) *AuthService {
	return &AuthService{
		AuthContextKey: authKey("authKey"),
		secret:         []byte(secret),
		tokenExpires:   tokenExpires,
	}
}

// VerifyToken verifies token and returns map containing its data
func (s *AuthService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Second * time.Duration(s.tokenExpires)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *AuthService) GetUserID(r *http.Request) (pgtype.UUID, error) {
	userUUID := pgtype.UUID{}
	claims, _ := r.Context().Value(s.AuthContextKey).(jwt.MapClaims)
	userID, ok := claims["user_id"].(string)
	if !ok {
		return userUUID, fmt.Errorf("cant read user_id")
	}
	if err := userUUID.Scan(userID); err != nil {
		return userUUID, fmt.Errorf("cant scan user_id")
	}
	return userUUID, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
