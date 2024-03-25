package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUniqueString() string {
	ctime := time.Now().UnixNano()

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	randomString := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	uniqueString := fmt.Sprintf("%d%s", ctime, randomString)
	return uniqueString
}
