package hashcash

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
)

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func base64EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func base64EncodeInt(n int) string {
	return base64EncodeBytes([]byte(strconv.Itoa(n)))
}

func sha1Hash(s string) string {
	hash := sha1.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
