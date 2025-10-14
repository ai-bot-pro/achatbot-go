package common

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	sessionID := "test-session"
	size := 5
	session := NewSession(sessionID, &size)

	if session == nil {
		t.Error("Expected new Session instance, got nil")
	}

	if session.GetSessionID() != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, session.GetSessionID())
	}

	if session.GetChatRound() != 0 {
		t.Errorf("Expected chat round 0, got %d", session.GetChatRound())
	}
}

func TestSessionInitChatMessage(t *testing.T) {
	session := NewSession("test-session", nil)
	msg := map[string]any{"role": "system", "content": "You are a helpful assistant"}
	session.InitChatMessage(msg)

	// Get the chat history list and check if the message is there
	historyList := session.GetChatHistory().ToList()
	if len(historyList) != 1 {
		t.Errorf("Expected history list length 1, got %d", len(historyList))
	}

	if historyList[0]["content"] != "You are a helpful assistant" {
		t.Errorf("Expected content 'You are a helpful assistant', got '%s'", historyList[0]["content"])
	}
}

func TestSessionReset(t *testing.T) {
	session := NewSession("test-session", nil)

	// Increment chat round and add some history
	session.IncrementChatRound()
	msg := map[string]any{"role": "user", "content": "Hello"}
	session.InitChatMessage(msg)

	// Reset the session
	session.Reset()

	if session.GetChatRound() != 0 {
		t.Errorf("Expected chat round 0 after reset, got %d", session.GetChatRound())
	}

	// Check if history is cleared
	historyList := session.GetChatHistory().ToList()
	if len(historyList) != 0 {
		t.Errorf("Expected empty history after reset, got length %d", len(historyList))
	}
}

func TestSessionSetChatHistorySize(t *testing.T) {
	session := NewSession("test-session", nil)
	size := 3
	session.SetChatHistorySize(&size)

	// Check if the size was set correctly by accessing the chat history directly
	chSize := session.GetChatHistory().size
	if chSize == nil || *chSize != size {
		t.Errorf("Expected chat history size %d, got %v", size, chSize)
	}
}

func TestSessionSetSessionID(t *testing.T) {
	session := NewSession("old-session", nil)
	newSessionID := "new-session"
	session.SetSessionID(newSessionID)

	if session.GetSessionID() != newSessionID {
		t.Errorf("Expected session ID %s, got %s", newSessionID, session.GetSessionID())
	}
}

func TestSessionIncrementChatRound(t *testing.T) {
	session := NewSession("test-session", nil)

	// Increment chat round multiple times
	session.IncrementChatRound()
	if session.GetChatRound() != 1 {
		t.Errorf("Expected chat round 1, got %d", session.GetChatRound())
	}

	session.IncrementChatRound()
	session.IncrementChatRound()
	if session.GetChatRound() != 3 {
		t.Errorf("Expected chat round 3, got %d", session.GetChatRound())
	}
}

func TestSessionGetters(t *testing.T) {
	sessionID := "test-session"
	size := 5
	session := NewSession(sessionID, &size)

	// Test GetSessionID
	if session.GetSessionID() != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, session.GetSessionID())
	}

	// Test GetChatRound
	if session.GetChatRound() != 0 {
		t.Errorf("Expected chat round 0, got %d", session.GetChatRound())
	}

	// Test GetChatHistory
	history := session.GetChatHistory()
	if history == nil {
		t.Error("Expected chat history, got nil")
	}
}
