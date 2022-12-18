package event

const (
	Login            = "Login"
	Registered       = "Registered"
	Message          = "Message"
	AppendRoom       = "AppendRoom"
	UserAppendToRoom = "UserAppendToRoom"
	UpdateRoomUsers  = "UpdateRoomUsers"
	Ping             = "Ping"
)

// TODO переименовать на EventName
type ChatEvent interface {
	Name() string
}
