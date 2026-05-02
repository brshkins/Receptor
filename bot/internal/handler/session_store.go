package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]*UserSession
	path     string
}

func NewSessionStore(path string) (*SessionStore, error) {
	ss := &SessionStore{
		sessions: make(map[int64]*UserSession),
		path:     path,
	}
	if ss.path == "" {
		return ss, nil
	}
	if err := ss.Load(); err != nil {
		return nil, err
	}
	return ss, nil
}

func (s *SessionStore) GetOrCreate(userID int64) *UserSession {
	s.mu.RLock()
	st := s.sessions[userID]
	s.mu.RUnlock()
	if st != nil {
		return st
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	st = s.sessions[userID]
	if st == nil {
		st = &UserSession{State: UserStateIdle}
		s.sessions[userID] = st
	}
	return st
}

// GetCopy returns a copy of the user session.
// It never exposes internal pointers, so callers can't race on session fields.
func (s *SessionStore) GetCopy(userID int64) UserSession {
	s.mu.Lock()
	st := s.sessions[userID]
	if st == nil {
		st = &UserSession{State: UserStateIdle}
		s.sessions[userID] = st
	}
	out := *st
	s.mu.Unlock()
	return out
}

func (s *SessionStore) Set(userID int64, state UserState, email, token string) error {
	s.mu.Lock()
	st := s.sessions[userID]
	if st == nil {
		st = &UserSession{}
		s.sessions[userID] = st
	}
	st.State = state
	st.Email = email
	st.Token = token
	s.mu.Unlock()

	return s.Save()
}

// SetState updates only the FSM state without touching token/user fields.
// Critical: navigation must reset FSM to idle without losing token.
// SetPanelMessage запоминает сообщение, которое нужно редактировать вместо отправки новых.
func (s *SessionStore) SetPanelMessage(userID int64, chatID int64, messageID int) error {
	s.mu.Lock()
	st := s.sessions[userID]
	if st == nil {
		st = &UserSession{State: UserStateIdle}
		s.sessions[userID] = st
	}
	st.PanelChatID = chatID
	st.PanelMessageID = messageID
	s.mu.Unlock()
	return s.Save()
}

func (s *SessionStore) SetState(userID int64, state UserState) error {
	s.mu.Lock()
	st := s.sessions[userID]
	if st == nil {
		st = &UserSession{}
		s.sessions[userID] = st
	}
	st.State = state
	s.mu.Unlock()
	return s.Save()
}

func (s *SessionStore) Load() error {
	if s.path == "" {
		return nil
	}
	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read session store: %w", err)
	}
	if len(b) == 0 {
		return nil
	}

	var raw map[int64]*UserSession
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("decode session store: %w", err)
	}

	s.mu.Lock()
	if raw == nil {
		raw = make(map[int64]*UserSession)
	}
	s.sessions = raw
	s.mu.Unlock()
	return nil
}

func (s *SessionStore) Save() error {
	if s.path == "" {
		return nil
	}

	// Create a snapshot under lock to avoid races on *UserSession fields.
	s.mu.RLock()
	snap := make(map[int64]UserSession, len(s.sessions))
	for id, st := range s.sessions {
		if st == nil {
			continue
		}
		snap[id] = *st
	}
	s.mu.RUnlock()

	b, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("encode session store: %w", err)
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir session store dir: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return fmt.Errorf("write session store: %w", err)
	}
	_ = os.Remove(s.path)
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("rename session store: %w", err)
	}
	return nil
}

