package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	input := HexToBytes(os.Args[1])
	output := DetectSingleByteXOR(input)
	fmt.Println(string(output))
}

func HexToBytes(str string) []byte {
	bytes, _ := hex.DecodeString(str)
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

func DetectSingleByteXOR(in []byte) []byte {
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

	return bestGuess
}

func ScoreEnglishPlaintext(plaintext []byte) (score int) {
	for _, r := range string(plaintext) {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			score += 1
		}
	}
	return score
}
