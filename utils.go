package main

import (
	"bytes"
	"encoding/hex"
	"unicode"
)

func HexToBytes(str string) []byte {
	bytes, _ := hex.DecodeString(str)
	return bytes
}

func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
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

func DetectSingleByteXOR(in []byte) ([]byte, int) {
	var maxScore int
	var bestGuess []byte
	length := len(in)

	for b := byte(0); b < 0xff; b++ {
		repeated := bytes.Repeat([]byte{b}, length)
		out := FixedLengthXOR(in, repeated)

		if score := ScoreEnglishPlaintext(out); score > maxScore {
			maxScore = score
			bestGuess = out
		}
	}

	return bestGuess, maxScore
}

func ScoreEnglishPlaintext(plaintext []byte) (score int) {
	for _, r := range string(plaintext) {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsDigit(r) {
			score += 1
		}
	}
	return score
}
