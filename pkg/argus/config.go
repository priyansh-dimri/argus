package argus

type SecurityMode string

const (
	LatencyFirst SecurityMode = "LATENCY_FIRST"
	SmartShield  SecurityMode = "SMART_SHIELD"
	Paranoid     SecurityMode = "PARANOID"
)

type Config struct {
	AppID  string
	APIKey string
	Mode   SecurityMode
}
