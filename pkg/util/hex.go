package util

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const prefix = "0x"

func Hex(b []byte) string {
	return fmt.Sprintf("%s%x", prefix, b)
}

func FromHex(h string) ([]byte, error) {
	if strings.HasPrefix(h, prefix) {
		return hex.DecodeString(strings.TrimPrefix(h, prefix))
	}
	return hex.DecodeString(h)
}
