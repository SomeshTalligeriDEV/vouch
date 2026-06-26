package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 30 * 24 * time.Hour
)

// SubjectType distinguishes whether the token belongs to a user or company.
type SubjectType string

const (
	SubjectTypeUser    SubjectType = "user"
	SubjectTypeCompany SubjectType = "company"
)

// Claims is the JWT payload Vouch issues.
type Claims struct {
	UserID      string      `json:"uid"`
	Username    string      `json:"username"`
	Role        string      `json:"role"`
	SubjectType SubjectType `json:"stype"`
	jwt.RegisteredClaims
}

// Manager issues and verifies access and refresh tokens.
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
}

// NewManager builds a Manager from the configured secrets.
func NewManager(accessSecret, refreshSecret string) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

// TokenPair is an access + refresh token issued together.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Generate issues a fresh access + refresh token pair for a user.
func (m *Manager) Generate(userID, username, role string) (*TokenPair, error) {
	return m.GenerateTyped(userID, username, role, SubjectTypeUser)
}

// GenerateTyped issues tokens with an explicit subject type (user or company).
func (m *Manager) GenerateTyped(id, name, role string, stype SubjectType) (*TokenPair, error) {
	access, err := m.sign(id, name, role, stype, m.accessSecret, accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("jwt.Generate access: %w", err)
	}
	refresh, err := m.sign(id, name, role, stype, m.refreshSecret, refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("jwt.Generate refresh: %w", err)
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
	}, nil
}

func (m *Manager) sign(id, name, role string, stype SubjectType, secret []byte, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:      id,
		Username:    name,
		Role:        role,
		SubjectType: stype,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "vouch",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// VerifyAccess validates an access token and returns its claims.
func (m *Manager) VerifyAccess(tokenString string) (*Claims, error) {
	return m.verify(tokenString, m.accessSecret)
}

// VerifyRefresh validates a refresh token and returns its claims.
func (m *Manager) VerifyRefresh(tokenString string) (*Claims, error) {
	return m.verify(tokenString, m.refreshSecret)
}

func (m *Manager) verify(tokenString string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.verify: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("jwt.verify: invalid token")
	}
	return claims, nil
}
