package utils

import (
	"math/rand"
	"time"
)

func CreateRandCode() string {
	rand.Seed(time.Now().UnixNano())
	CODE_CHARSET := []rune("12346789ABCDEFGHJKLMNPQRTUVWXY")

	b := make([]rune, 5)
	for i := range b {
		b[i] = CODE_CHARSET[rand.Intn(len(CODE_CHARSET))]
	}
	return string(b)
}
