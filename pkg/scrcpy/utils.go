package scrcpy

func ConvertTouchMessage(msg *TouchMessage) *TouchEvent {
	var action AndroidTouchAction
	switch msg.Action {
	case TouchMessageActionDown:
		action = AndroidTouchActionDown
	case TouchMessageActionUp:
		action = AndroidTouchActionUp
	case TouchMessageActionMove:
		action = AndroidTouchActionMove
	case TouchMessageActionCancel:
		action = AndroidTouchActionMove
	default:
		action = AndroidTouchActionMove
	}

	pressure := msg.Pressure
	if pressure == 0 {
		pressure = 1.0
	}

	return &TouchEvent{
		Action:    action,
		PointerId: msg.PointerId,
		X:         msg.X,
		Y:         msg.Y,
		Pressure:  pressure,
		Buttons:   AndroidTouchButtonPrimary,
	}
}
