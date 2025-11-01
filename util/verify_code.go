package util

import (
	crand "crypto/rand"
	"golang.org/x/exp/rand"
	"math/big"
	"strconv"
	"sync"
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

type VerifyCodeGenerator struct {
	rnd  *rand.Rand
	mu   sync.Mutex
	pool sync.Pool
}

func NewVerifyCodeGenerator() *VerifyCodeGenerator {
	g := &VerifyCodeGenerator{
		rnd: rand.New(rand.NewSource(uint64(time.Now().UnixNano()))),
	}
	g.pool.New = func() interface{} {
		return make([]byte, 6)
	}
	return g
}

func (g *VerifyCodeGenerator) Generate() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	buf := g.pool.Get().([]byte)
	defer g.pool.Put(buf)

	for i := range buf {
		buf[i] = byte(g.rnd.Intn(10)) + '0'
	}
	return string(buf)
}

func (g *VerifyCodeGenerator) GenerateSecure() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	buf := g.pool.Get().([]byte)
	defer g.pool.Put(buf)

	for i := range buf {
		randNum, err := crand.Int(crand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		buf[i] = byte(randNum.Int64()) + '0'
	}
	return string(buf), nil
}
