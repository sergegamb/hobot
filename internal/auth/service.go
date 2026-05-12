package auth

import "sync"

type Service struct {
	mu sync.RWMutex

	approved map[int64]bool
	pending  map[int64]bool
}

func NewService() *Service {
	return &Service{
		approved: make(map[int64]bool),
		pending:  make(map[int64]bool),
	}
}

func (s *Service) IsApproved(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.approved[userID]
}

func (s *Service) IsPending(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.pending[userID]
}

func (s *Service) AddPending(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pending[userID] = true
}

func (s *Service) Approve(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pending, userID)

	s.approved[userID] = true
}

func (s *Service) Deny(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.pending, userID)
}
