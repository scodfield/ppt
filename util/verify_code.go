package util

import (
	crand "crypto/rand"
	"golang.org/x/exp/rand"
	"math/big"
	"strconv"
	"time"
)

// GenerateDigitVerifyCode 生成6位数字验证码-简易版
func GenerateDigitVerifyCode() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	randNum := rand.Intn(900000) + 100000
	return strconv.Itoa(randNum)
}

// GenerateSecureDigitVerifyCode 生成6位数字验证码-安全版
func GenerateSecureDigitVerifyCode() string {
	max := big.NewInt(900000)
	n, err := crand.Int(crand.Reader, max)
	if err != nil {
		return GenerateDigitVerifyCode()
	}
	randNum := n.Int64() + 100000
	return strconv.Itoa(int(randNum))
}
