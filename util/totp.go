package util

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"ppt/log"
	"sync"
	"time"
)

func GenerateTOTPKey(issuer string, account string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	if err != nil {
		log.Error("GenerateTOTPKey totp Generate error", zap.String("issuer", issuer), zap.String("account", account), zap.Error(err))
		return nil, err
	}
	return key, nil
}

func ValidateTOTP(code, secret string) bool {
	return totp.Validate(code, secret)
}

type RotatingTOTP struct {
	issuer         string
	account        string
	currentSecret  string
	nextSecret     string
	rotationPeriod time.Duration
	lastRotation   time.Time
	mutex          sync.RWMutex
}

// 是否需要轮询密钥
func (r *RotatingTOTP) checkRotation() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if time.Since(r.lastRotation) > r.rotationPeriod {
		newNextSecret, _ := GenerateTOTPKey(r.issuer, r.account)
		r.currentSecret = r.nextSecret
		r.nextSecret = newNextSecret.Secret()
		r.lastRotation = time.Now()
	}
}

// Validate 验证TOTP(同时轮询密钥)
func (r *RotatingTOTP) Validate(code string) bool {
	r.checkRotation()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 检查当前密钥
	if totp.Validate(code, r.currentSecret) {
		return true
	}

	// 检查下一个密钥
	if totp.Validate(code, r.nextSecret) {
		return true
	}
	return false
}

// CurrentSecret 获取当前密钥
func (r *RotatingTOTP) CurrentSecret() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.currentSecret
}

func NewRotatingTOTP(secret string, issuer string, account string, rotationPeriod time.Duration) (*RotatingTOTP, error) {
	nextSecret, err := GenerateTOTPKey(issuer, account)
	if err != nil {
		return nil, err
	}
	return &RotatingTOTP{
		issuer:         issuer,
		account:        account,
		currentSecret:  secret,
		nextSecret:     nextSecret.Secret(),
		rotationPeriod: rotationPeriod,
		lastRotation:   time.Now(),
	}, nil
}
