package api

import (
	"sync"

	"my-blog-backend/internal/ssh"
)

// SessionStore provides a threadsafe session registry shared by SSH handlers.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*ssh.Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*ssh.Session),
	}
}

func (s *SessionStore) Get(id string) (*ssh.Session, bool) {
	s.mu.RLock()
	session, ok := s.sessions[id]
	s.mu.RUnlock()
	return session, ok
}

func (s *SessionStore) Set(id string, session *ssh.Session) {
	s.mu.Lock()
	s.sessions[id] = session
	s.mu.Unlock()
}

func (s *SessionStore) Delete(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}

// List returns a snapshot copy for safe iteration.
func (s *SessionStore) List() map[string]*ssh.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make(map[string]*ssh.Session, len(s.sessions))
	for id, session := range s.sessions {
		snapshot[id] = session
	}
	return snapshot
}
