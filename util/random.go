package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// ** Initialization
func init() {
	rand.Seed(time.Now().UnixMicro())
}

// ** Random Int Generator
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// ** Random String Generator
func RandomString(n int) string {
	var stringBuilder strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		stringBuilder.WriteByte(c)
	}
	return stringBuilder.String()
}

// ** Random Owner Generator
func RandomOwner() string {
	return RandomString(6)
}

// ** Random Money Generator
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// ** Random Currency Generator
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD", "TRY"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// ** Random Amount Generator
func RandomAmount() int64 {
	return RandomInt(-1000, 1000) // can be negative or positive
}
