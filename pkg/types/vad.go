package types

// VADState 表示VAD状态
type VADState int

const (
	Quiet VADState = iota
	Starting
	Speaking
	Stopping
)

func (s VADState) String() string {
	switch s {
	case Quiet:
		return "QUIET"
	case Starting:
		return "STARTING"
	case Speaking:
		return "SPEAKING"
	case Stopping:
		return "STOPPING"
	default:
		return "UNKNOWN"
	}
}
