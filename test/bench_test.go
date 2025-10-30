package test

import (
	"ppt/util"
	"testing"
)

func BenchmarkGenVerifyCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = util.GenerateDigitVerifyCode()
	}
}

func BenchmarkGenVerifyCodeCrypto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = util.GenerateSecureDigitVerifyCode()
	}
}
