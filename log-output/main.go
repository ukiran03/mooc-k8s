package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "01234567890123456789"

func stringWithCharset(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func newRandString() string {
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		stringWithCharset(8),
		stringWithCharset(4),
		stringWithCharset(4),
		stringWithCharset(4),
		stringWithCharset(8),
	)
}

func main() {
	str := newRandString()
	log.Println(str)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println(str)
	}
}
