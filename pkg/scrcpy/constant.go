package scrcpy

type AndroidControlMsgType byte
type AndroidTouchAction byte
type AndroidKeyAction byte
type AndroidKeyCode int32
type TouchMessageAction string
type MessageType string
type KeyCommand string

const (
	// Scrcpy constants
	SCRCPY_SERVER_VERSION     string = "3.3.1"
	SCRCPY_SERVER_TCP_PORT    uint16 = 1234
	SCRCPY_REMOTE_SERVER_PATH string = "/data/local/tmp/scrcpy-server.jar"

	// Scrcpy device info constants
	DEVICE_NAME_LENGTH  uint8 = 64
	DUMMY_LENGTH        uint8 = 1
	VIDEO_HEADER_LENGTH uint8 = 12

	// Scrcpy control constants
	// Touch event actions
	AndroidTouchActionDown AndroidTouchAction = 0
	AndroidTouchActionUp   AndroidTouchAction = 1
	AndroidTouchActionMove AndroidTouchAction = 2

	// Touch event buttons
	AndroidTouchButtonPrimary   uint32 = 1
	AndroidTouchButtonSecondary uint32 = 2
	AndroidTouchButtonTertiary  uint32 = 4

	// Control message types
	ControlMsgTypeInjectTouchEvent AndroidControlMsgType = 2
	ControlMsgTypeInjectKeyEvent   AndroidControlMsgType = 0

	// Key event actions
	AndroidKeyActionDown AndroidKeyAction = 0
	AndroidKeyActionUp   AndroidKeyAction = 1

	// Android key codes
	AndroidKeyCodeHome   AndroidKeyCode = 3
	AndroidKeyCodeBack   AndroidKeyCode = 4
	AndroidKeyCodePower  AndroidKeyCode = 26
	AndroidKeyCodeRecent AndroidKeyCode = 187

	// Touch message actions
	TouchMessageActionDown   TouchMessageAction = "down"
	TouchMessageActionUp     TouchMessageAction = "up"
	TouchMessageActionMove   TouchMessageAction = "move"
	TouchMessageActionCancel TouchMessageAction = "cancel"

	// Message types
	MessageTypeTouch MessageType = "touch"
	MessageTypeKey   MessageType = "key"

	// Key commands
	KeyCommandBack   KeyCommand = "back"
	KeyCommandHome   KeyCommand = "home"
	KeyCommandRecent KeyCommand = "recent"
	KeyCommandPower  KeyCommand = "power"
)

// Map KeyCommand to Android key codes
var KeyCodeMap = map[KeyCommand]AndroidKeyCode{
	KeyCommandBack:   AndroidKeyCodeBack,
	KeyCommandHome:   AndroidKeyCodeHome,
	KeyCommandRecent: AndroidKeyCodeRecent,
	KeyCommandPower:  AndroidKeyCodePower,
}
