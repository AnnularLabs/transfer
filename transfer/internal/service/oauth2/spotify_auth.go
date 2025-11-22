package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyOAuthService interface {
	GetAuthURL(state string) string
	HandleCallback(code, state string) (*UserToken, error)
	RefreshUserToken(refreshToken string) (*UserToken, error)
	GetAuthenticatedClient(userID string) (spotify.Client, error)
	RevokeToken(userID string) error
}

// TokenManager Token 管理接口
type TokenManager interface {
	StoreUserToken(userID string, token *UserToken) error
	GetUserToken(userID string) (*UserToken, error)
	DeleteUserToken(userID string) error
	IsTokenValid(userID string) bool
}

type UserToken struct {
	UserID       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scopes       []string  `json:"scopes"`
}

type SpotifyAuth struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}

// MemoryTokenManager 内存 Token 管理器
type MemoryTokenManager struct {
	tokens map[string]*UserToken
	mutex  sync.RWMutex
}

func NewMemoryTokenManager() TokenManager {
	return &MemoryTokenManager{
		tokens: make(map[string]*UserToken),
	}
}

func (m *MemoryTokenManager) StoreUserToken(userID string, token *UserToken) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tokens[userID] = token
	return nil
}

func (m *MemoryTokenManager) GetUserToken(userID string) (*UserToken, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	token, exists := m.tokens[userID]
	if !exists {
		return nil, fmt.Errorf("token not found for user %s", userID)
	}
	return token, nil
}

func (m *MemoryTokenManager) DeleteUserToken(userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.tokens, userID)
	return nil
}

func (m *MemoryTokenManager) IsTokenValid(userID string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	token, exists := m.tokens[userID]
	if !exists {
		return false
	}
	return time.Now().Before(token.ExpiresAt)
}

// SpotifyOAuth OAuth 服务实现
type SpotifyOAuth struct {
	authenticator spotify.Authenticator
	config        *oauth2.Config // 用于 refresh token 操作
	tokenManager  TokenManager
}

func NewSpotifyOAuth(clientID, clientSecret, redirectURL string, tokenManager TokenManager) SpotifyOAuthService {
	auth := spotify.NewAuthenticator(redirectURL,
		spotify.ScopeUserReadPrivate,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
	)
	auth.SetAuthInfo(clientID, clientSecret)

	// 创建标准 OAuth2 配置用于 refresh token
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user-read-private",
			"playlist-read-private",
			"playlist-modify-public",
			"playlist-modify-private",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotify.AuthURL,
			TokenURL: spotify.TokenURL,
		},
	}

	return &SpotifyOAuth{
		authenticator: auth,
		config:        config,
		tokenManager:  tokenManager,
	}
}

func (s *SpotifyOAuth) GetAuthURL(state string) string {
	return s.authenticator.AuthURL(state)
}

func (s *SpotifyOAuth) HandleCallback(code, state string) (*UserToken, error) {
	// 交换授权码获取 token
	token, err := s.authenticator.Exchange(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	userToken := &UserToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresAt:    token.Expiry,
		Scopes:       []string{"user-read-private", "playlist-read-private", "playlist-modify-public", "playlist-modify-private"},
	}

	return userToken, nil
}

func (s *SpotifyOAuth) RefreshUserToken(refreshToken string) (*UserToken, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := s.config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	userToken := &UserToken{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		TokenType:    newToken.TokenType,
		ExpiresAt:    newToken.Expiry,
		Scopes:       s.config.Scopes,
	}

	return userToken, nil
}

func (s *SpotifyOAuth) GetAuthenticatedClient(userID string) (spotify.Client, error) {
	token, err := s.tokenManager.GetUserToken(userID)
	if err != nil {
		return spotify.Client{}, err
	}

	// 检查 token 是否过期，如需要则刷新
	if time.Now().After(token.ExpiresAt) {
		newToken, err := s.RefreshUserToken(token.RefreshToken)
		if err != nil {
			return spotify.Client{}, fmt.Errorf("failed to refresh expired token: %w", err)
		}
		newToken.UserID = userID
		s.tokenManager.StoreUserToken(userID, newToken)
		token = newToken
	}

	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.ExpiresAt,
	}

	return spotify.Authenticator{}.NewClient(oauthToken), nil
}

func (s *SpotifyOAuth) RevokeToken(userID string) error {
	return s.tokenManager.DeleteUserToken(userID)
}

func NewSpotifyAuth(clientID, clientSecret string) spotify.Client {
	authConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}
	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		panic(err)
	}

	return spotify.Authenticator{}.NewClient(accessToken)
}

// GenerateSecureState 生成安全的 state 参数（添加这个函数）
func GenerateSecureState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return fmt.Sprintf("state_%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)
}
