package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"unicode"
)

func HexToBytes(str string) []byte {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return bytes
}

func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func Base64ToBytes(str string) []byte {
	bytes, _ := base64.StdEncoding.DecodeString(str)
	return bytes
}

// Both buffers must have the same length.
func FixedLengthXOR(a, b []byte) []byte {
	length := len(a)
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}

func SingleByteXOR(plaintext []byte, key byte) []byte {
	length := len(plaintext)
	repeated := bytes.Repeat([]byte{key}, length)
	return FixedLengthXOR(plaintext, repeated)
}

func DetectSingleByteXOR(in []byte) ([]byte, byte) {
	var maxScore int
	var bestGuess []byte
	var bestGuessKey byte

	for b := byte(0); b < 0xff; b++ {
		out := SingleByteXOR(in, b)

		if score := ScoreEnglishPlaintext(out); score > maxScore {
			maxScore = score
			bestGuess = out
			bestGuessKey = b
		}
	}

	return bestGuess, bestGuessKey
}

var EnglishFrequency = map[rune]int{
	'a': 8,
	'b': 1,
	'c': 3,
	'd': 4,
	'e': 12,
	'f': 2,
	'g': 2,
	'h': 6,
	'i': 7,
	'j': 0,
	'k': 1,
	'l': 4,
	'm': 2,
	'n': 7,
	'o': 8,
	'p': 2,
	'q': 0,
	'r': 6,
	's': 6,
	't': 9,
	'u': 3,
	'v': 1,
	'w': 2,
	'x': 0,
	'y': 2,
	'z': 0,
	' ': 10,
}

func ScoreEnglishPlaintext(plaintext []byte) (total int) {
	for _, r := range string(plaintext) {
		lower := unicode.ToLower(r)
		if score, ok := EnglishFrequency[lower]; ok {
			// Add one because the fact that we matched one
			// of the characters constitutes "success".
			total += score + 1
		}
	}
	return total
}
