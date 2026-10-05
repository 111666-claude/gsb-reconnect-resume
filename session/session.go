// Package session 管理断线重连的续局会话。
package session

import "fmt"

// Session 是一次续局会话的快照。
type Session struct {
	Token   string `json:"token"`
	Player  string `json:"player"`
	Version int    `json:"version"`
	State   string `json:"state"`
	LastMs  int64  `json:"lastMs"`
}

// Manager 是会话台账。
type Manager struct {
	GraceMs int64

	sessions []*Session
	byToken  map[string]*Session
	seen     map[string]Session
	Scanned  int
	Resumes  int
	Rejected int
}

// Open 开一局新会话并返回令牌。
func (m *Manager) Open(player string, atMs int64) Session {
	s := &Session{
		Token:   fmt.Sprintf("tok-%s-%d", player, atMs),
		Player:  player,
		Version: 1,
		State:   "active",
		LastMs:  atMs,
	}
	m.sessions = append(m.sessions, s)
	if m.byToken == nil {
		m.byToken = map[string]*Session{}
	}
	m.byToken[s.Token] = s
	return *s
}

func (m *Manager) find(token string) *Session {
	m.Scanned++
	return m.byToken[token]
}

// Resume 用令牌续局，返回会话快照与拒绝原因。
func (m *Manager) Resume(requestID, token string, version int, atMs int64) (Session, string) {
	if m.seen == nil {
		m.seen = map[string]Session{}
	}
	if snap, ok := m.seen[requestID]; ok {
		return snap, ""
	}
	s := m.find(token)
	if s == nil {
		m.Rejected++
		return Session{}, "unknown"
	}
	if atMs-s.LastMs > m.GraceMs {
		m.Rejected++
		return Session{}, "expired"
	}
	if version != s.Version {
		m.Rejected++
		return Session{}, "stale"
	}
	m.Resumes++
	s.Version++
	s.LastMs = atMs
	s.State = "resumed"
	snap := *s
	m.seen[requestID] = snap
	return snap, ""
}
