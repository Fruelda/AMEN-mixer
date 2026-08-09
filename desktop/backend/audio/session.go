package audio

type SessionManager struct {
	sessions []Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: []Session{},
	}
}

func (s *SessionManager) Sessions() []Session {
	return s.sessions
}
