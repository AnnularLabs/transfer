package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionData struct {
	UserID      string    `json:"user_id"`
	SpotifyID   string    `json:"spotify_id"`
	IsAuthed    bool      `json:"is_authed"`
	AuthTime    time.Time `json:"auth_time"`
	State       string    `json:"state"`        // OAuth state
	StateExpiry time.Time `json:"state_expiry"` // state 过期时间
}

type SessionManager interface {
	SetSession(c *gin.Context, data *SessionData) error
	GetSession(c *gin.Context) *SessionData
	DeleteSession(c *gin.Context) error
	SetState(c *gin.Context, state string) error
	ValidateState(c *gin.Context, state string) bool
}

type MemorySessionManager struct {
	sessions map[string]*SessionData
	mutex    sync.RWMutex
}

func NewMemorySessionManager() SessionManager {
	return &MemorySessionManager{
		sessions: make(map[string]*SessionData),
	}
}

func (m *MemorySessionManager) getSessionKey(c *gin.Context) string {
	return fmt.Sprintf("%s_%s", c.ClientIP(), c.GetHeader("User-Agent"))
}

func (m *MemorySessionManager) SetSession(c *gin.Context, data *SessionData) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := m.getSessionKey(c)
	m.sessions[key] = data

	c.SetCookie("session_id", "authenticated", 3600, "/", "", false, true)

	return nil
}

func (m *MemorySessionManager) GetSession(c *gin.Context) *SessionData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	key := m.getSessionKey(c)
	session, exists := m.sessions[key]
	if !exists {
		return nil
	}

	return session
}

func (m *MemorySessionManager) DeleteSession(c *gin.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := m.getSessionKey(c)
	delete(m.sessions, key)

	// 清除 cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)

	return nil
}

func (m *MemorySessionManager) SetState(c *gin.Context, state string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	key := m.getSessionKey(c)
	session := m.sessions[key]
	if session == nil {
		session = &SessionData{}
		m.sessions[key] = session
	}

	session.State = state
	session.StateExpiry = time.Now().Add(10 * time.Minute) // state 10分钟过期

	return nil
}

func (m *MemorySessionManager) ValidateState(c *gin.Context, state string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	key := m.getSessionKey(c)
	session, exists := m.sessions[key]
	if !exists {
		return false
	}

	// 检查 state 是否匹配且未过期
	if session.State != state || time.Now().After(session.StateExpiry) {
		return false
	}

	return true
}
