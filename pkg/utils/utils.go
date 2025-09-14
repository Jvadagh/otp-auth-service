package utils

import (
	"fmt"
	"regexp"
)

// NormalizePhone validates Iranian phone numbers and normalizes them to +98 format.
func NormalizePhone(phone string) (string, error) {
	re1 := regexp.MustCompile(`^\+989\d{9}$`)
	re2 := regexp.MustCompile(`^09\d{9}$`)
	if re1.MatchString(phone) {
		return phone, nil
	}
	if re2.MatchString(phone) {
		return "+98" + phone[1:], nil
	}
	return "", fmt.Errorf("invalid phone number format")
}
