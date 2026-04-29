package utils

import (
	"math/rand"
	"time"

	"github.com/gosimple/slug"
)

func GenerateSlug(name string) string {
	baseSlug := slug.Make(name)

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	randomStr := string(b)

	return baseSlug + "-" + randomStr
}
