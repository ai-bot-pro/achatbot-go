package consts

type MessageType int

// Constants for WebSocket message types
const (
	TextMessage   MessageType = 1
	BinaryMessage MessageType = 2
)

func (m MessageType) String() string {
	switch m {
	case TextMessage:
		return "Text"
	case BinaryMessage:
		return "Binary"
	default:
		return "Unknown"
	}
}
