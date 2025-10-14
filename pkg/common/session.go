package common

// Session represents a chat session with chat history
type Session struct {
	chatRound   int
	sessionID   string
	chatHistory *ChatHistory
}

// NewSession creates a new Session instance
func NewSession(sessionID string, chatHistorySize *int) *Session {
	// Create a new ChatHistory with the given size and no initial messages
	chatHistory := NewChatHistory(chatHistorySize, nil, nil)

	return &Session{
		chatRound:   0,
		sessionID:   sessionID,
		chatHistory: chatHistory,
	}
}

// InitChatMessage initializes the chat with a message
func (s *Session) InitChatMessage(initChatMessage map[string]any) {
	s.chatHistory.Init(initChatMessage)
}

// Reset resets the session chat round and clears chat history
func (s *Session) Reset() {
	s.chatRound = 0
	s.chatHistory.Clear()
	// Also reset the initial chat message to fully clear the history
	s.chatHistory.Init(nil)
}

// SetChatHistorySize sets the size limit of the chat history
func (s *Session) SetChatHistorySize(chatHistorySize *int) {
	s.chatHistory.SetSize(chatHistorySize)
}

// SetSessionID sets the session ID
func (s *Session) SetSessionID(sessionID string) {
	s.sessionID = sessionID
}

// IncrementChatRound increments the chat round counter
func (s *Session) IncrementChatRound() {
	s.chatRound++
}

// GetChatRound returns the current chat round
func (s *Session) GetChatRound() int {
	return s.chatRound
}

// GetSessionID returns the session ID
func (s *Session) GetSessionID() string {
	return s.sessionID
}

// GetChatHistory returns the chat history
func (s *Session) GetChatHistory() *ChatHistory {
	return s.chatHistory
}
