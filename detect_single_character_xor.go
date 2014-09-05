package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"unicode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var maxScore int
	var bestGuess []byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := HexToBytes(scanner.Text())
		out, score := DetectSingleByteXOR(line)
		if score > maxScore {
			maxScore = score
			bestGuess = out
		}
	}

	fmt.Println(string(bestGuess))

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
