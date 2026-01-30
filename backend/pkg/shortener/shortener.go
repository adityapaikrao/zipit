package shortener

import "fmt"

// Shortener defines methods for encoding and decoding short URLs.
type Shortener interface {
	Encode(id int64) string
	Decode(shortCode string) (int64, error)
}

type base62shortener struct {
	alphabet string
	base     int64
	keyMap   map[rune]int
}

func NewBase62Shortener() Shortener {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	mapKey := make(map[rune]int)
	for i, char := range alphabet {
		mapKey[char] = i
	}
	return &base62shortener{
		alphabet: alphabet,
		base:     62,
		keyMap:   mapKey,
	}
}

func (bsh *base62shortener) Encode(id int64) string {
	if id == 0 {
		return string(bsh.alphabet[0])
	}
	encoded := make([]byte, 0)
	for id > 0 {
		rem := id % bsh.base
		id /= bsh.base
		encoded = append(encoded, bsh.alphabet[int(rem)])
	}
	// reverse to get the correct order
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}
	return string(encoded)
}

func (bsh *base62shortener) Decode(shortCode string) (int64, error) {
	var id int64 = 0
	for _, char := range shortCode {
		val, exists := bsh.keyMap[char]
		if !exists {
			return -1, fmt.Errorf("invalid character %c in short code %s", char, shortCode)
		}
		id = id*bsh.base + int64(val)
	}
	return id, nil
}
