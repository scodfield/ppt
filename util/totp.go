package util

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"ppt/logger"
)

func GenerateTOTPKey(issuer string, account string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	if err != nil {
		logger.Error("GenerateTOTPKey totp Generate error", zap.String("issuer", issuer), zap.String("account", account), zap.Error(err))
		return nil, err
	}
	return key, nil
}

func ValidateTOTP(code, secret string) bool {
	return totp.Validate(code, secret)
}
