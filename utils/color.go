package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
)

func GenColorCode(txt string) int {
	txt = strings.ToLower(txt)

	hasher := md5.New()
	hasher.Write([]byte(txt))
	hash := hex.EncodeToString(hasher.Sum(nil))

	colorInt, err := strconv.ParseInt(hash[:6], 16, 64)
	if err != nil {
		return rand.Intn(16777215)
	}

	return int(colorInt) & 0xFFFFFF
}
