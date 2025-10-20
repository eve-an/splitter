package session

import (
	"sync"
	"time"
)

type Session struct {
	ID        string
	ExpiresAt time.Time
}

type Service struct {
	mu       sync.RWMutex
	sessions map[string]Session
	ttl      time.Duration
}

func NewService(ttl time.Duration) *Service {
	return &Service{
		sessions: make(map[string]Session),
		ttl:      ttl,
	}
}

func (s *Service) Create(username string) Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateSessionID()
	sess := Session{
		ID:        id,
		ExpiresAt: time.Now().Add(s.ttl),
	}

	s.sessions[id] = sess
	return sess
}

func (s *Service) Get(id string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[id]
	if !ok || sess.ExpiresAt.Before(time.Now()) {
		return Session{}, false
	}
	return sess, true
}

func (s *Service) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}
