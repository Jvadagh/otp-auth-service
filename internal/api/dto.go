package api

// RequestOTPRequest represents the input body for requesting an OTP
type RequestOTPRequest struct {
	Phone string `json:"phone"`
}

// VerifyOTPRequest represents the input body for verifying an OTP
type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}
