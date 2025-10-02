package model

type MirroringSetupRequest struct {
	FPS        uint8  `json:"fps" validate:"required,oneof=30 60"`
	Bitrate    uint32 `json:"bitrate" validate:"required,oneof=1000000 2000000 3000000 4000000 5000000 6000000 7000000 8000000"`
	Resolution uint16 `json:"resolution" validate:"required,oneof=360 480 720 1080"`
}
