package context

import (
	"sync"
	"time"
)

// UserTicketState holds the state of a user's ticket list interaction
type UserTicketState struct {
	CurrentPage    int
	SelectedFilter string
	LastUpdated    time.Time
}

// SessionStore manages user sessions in memory
// In production, this could be replaced with Redis
type SessionStore struct {
	sessions map[int64]*UserTicketState
	mu       sync.RWMutex
	ttl      time.Duration
}

// NewSessionStore creates a new session store with default TTL of 1 hour
func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[int64]*UserTicketState),
		ttl:      1 * time.Hour,
	}

	// Start cleanup goroutine to remove expired sessions
	go store.cleanupExpiredSessions()

	return store
}

// GetUserState retrieves the current state for a user
// Returns default state if user has no saved state
func (s *SessionStore) GetUserState(userID int64) *UserTicketState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if state, exists := s.sessions[userID]; exists {
		// Check if session has expired
		if time.Since(state.LastUpdated) < s.ttl {
			return state
		}
		// Session expired, will be cleaned up by cleanup goroutine
	}

	// Return default state
	return &UserTicketState{
		CurrentPage:    1,
		SelectedFilter: "all",
		LastUpdated:    time.Now(),
	}
}

// SaveUserState saves or updates the state for a user
func (s *SessionStore) SaveUserState(userID int64, state *UserTicketState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state.LastUpdated = time.Now()
	s.sessions[userID] = state
}

// UpdateUserPage updates only the current page for a user
func (s *SessionStore) UpdateUserPage(userID int64, page int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.sessions[userID]
	if !exists {
		state = &UserTicketState{
			CurrentPage:    page,
			SelectedFilter: "all",
		}
	}

	state.CurrentPage = page
	state.LastUpdated = time.Now()
	s.sessions[userID] = state
}

// UpdateUserFilter updates only the selected filter for a user
// Also resets the page to 1 when filter changes
func (s *SessionStore) UpdateUserFilter(userID int64, filter string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.sessions[userID]
	if !exists {
		state = &UserTicketState{
			CurrentPage:    1,
			SelectedFilter: filter,
		}
	} else {
		state.CurrentPage = 1 // Reset to first page when filter changes
		state.SelectedFilter = filter
	}

	state.LastUpdated = time.Now()
	s.sessions[userID] = state
}

// DeleteUserState removes the state for a user
func (s *SessionStore) DeleteUserState(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, userID)
}

// cleanupExpiredSessions removes expired sessions periodically
func (s *SessionStore) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		now := time.Now()
		for userID, state := range s.sessions {
			if now.Sub(state.LastUpdated) > s.ttl {
				delete(s.sessions, userID)
			}
		}

		s.mu.Unlock()
	}
}

// SetTTL sets the time-to-live for sessions
func (s *SessionStore) SetTTL(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ttl = ttl
}

// GetSessionCount returns the number of active sessions
func (s *SessionStore) GetSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.sessions)
}
