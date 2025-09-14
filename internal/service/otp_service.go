package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

type OTPService struct {
	redis       *redis.Client
	MaxRequests int           // e.g. 3
	Window      time.Duration // e.g. 10 * time.Minute
	TTL         time.Duration // e.g. 2 * time.Minute
}

func NewOTPService(rdb *redis.Client) *OTPService {
	return &OTPService{
		redis:       rdb,
		MaxRequests: 3,
		Window:      10 * time.Minute,
		TTL:         2 * time.Minute,
	}
}

func (s *OTPService) GenerateOTP(phone string) (string, error) {
	ctx := context.Background()
	rlKey := fmt.Sprintf("otp:req:%s", phone)

	count, err := s.redis.Incr(ctx, rlKey).Result()
	if err != nil {
		return "", err
	}
	if count == 1 {
		if err := s.redis.Expire(ctx, rlKey, s.Window).Err(); err != nil {
			return "", err
		}
	}

	if count > int64(s.MaxRequests) {
		return "", fmt.Errorf("rate limit exceeded")
	}

	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	otpKey := fmt.Sprintf("otp:%s", phone)
	if err := s.redis.Set(ctx, otpKey, code, s.TTL).Err(); err != nil {
		return "", err
	}

	return code, nil
}

func (s *OTPService) ValidateOTP(phone, code string) (bool, error) {
	ctx := context.Background()
	otpKey := fmt.Sprintf("otp:%s", phone)

	stored, err := s.redis.Get(ctx, otpKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if stored == code {
		_ = s.redis.Del(ctx, otpKey).Err()
		return true, nil
	}
	return false, nil
}
