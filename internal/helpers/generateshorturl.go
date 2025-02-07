package helpers

import (
	"time"

	"golang.org/x/exp/rand"
)

var SeededRand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

// stringWithCharset generates a random string of a specified length using the provided character set.
// Parameters:
//   - length: The desired length of the generated string.
//   - charset: A string containing the characters to use for generating the random string.
//
// Returns:
//   - A random string composed of characters from the given charset.
func StringWithCharset() (string, error) {
	const charSet string = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

	rangeStart := 2
	rangeEnd := 10
	offset := rangeEnd - rangeStart
	randLength := SeededRand.Intn(offset) + rangeStart

	// if length <= 0 {
	// 	return "", &storage.InvalidLengthError{}
	// }

	// if charset == "" {
	// 	return "", &storage.InvalidCharSetError{}
	// }

	b := make([]byte, randLength)
	for i := range b {
		b[i] = charSet[SeededRand.Intn(len(charSet)-1)]
	}
	return string(b), nil
}
