package household

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const codeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const codeLength = 8

func GenerateJoinCode() (string, error) {
	code := make([]byte, codeLength)
	max := big.NewInt(int64(len(codeChars)))
	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("generate join code: %w", err)
		}
		code[i] = codeChars[n.Int64()]
	}
	return string(code), nil
}
