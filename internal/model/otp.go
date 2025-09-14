package model

import "time"

type OTP struct {
	Phone     string
	Code      string
	ExpiresAt time.Time
}
