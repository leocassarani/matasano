package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	in := Base64ToBytes(string(file))
	out := BreakRepeatingKeyXOR(in)
	fmt.Println(string(out))
}

func BreakRepeatingKeyXOR(in []byte) []byte {
	keySizes := DetectKeySizes(in)

	var bestKey []byte
	max := 0

	for _, keySize := range keySizes {
		key := BreakRepeatingKeyXORWithSize(in, keySize)
		plaintext := RepeatingKeyXOR(in, key)
		score := ScoreEnglishPlaintext(plaintext)
		if score > max {
			max = score
			bestKey = key
		}
	}

	return RepeatingKeyXOR(in, bestKey)
}

func BreakRepeatingKeyXORWithSize(in []byte, keySize int) []byte {
	blocks := SplitBytes(in, keySize)
	transposed := TransposeBlocks(blocks)
	return FindRepeatingXORKey(transposed)
}

func DetectKeySizes(in []byte) (guesses []int) {
	distances := map[int]float64{}

	for keySize := 2; keySize <= 40; keySize++ {
		first := in[0:keySize]
		second := in[keySize:(keySize * 2)]
		third := in[(keySize * 2):(keySize * 3)]
		fourth := in[(keySize * 3):(keySize * 4)]

		dist := HammingDistance(first, second) +
			HammingDistance(second, third) +
			HammingDistance(third, fourth)

		distances[keySize] = float64(dist) / float64(3*keySize)
	}

	// Hamming distance can be 8 at most.
	var topScores []float64
	for _, dist := range distances {
		topScores = append(topScores, dist)
	}
	sort.Float64s(topScores)

	// Take the sizes with the top 3 scores.
	for size, dist := range distances {
		for i, score := range topScores {
			if i > 2 {
				break
			}

			if dist <= score {
				guesses = append(guesses, size)
				break
			}
		}
	}
	return guesses
}

// The two bytestrings must be of equal length.
func HammingDistance(a, b []byte) (dist int) {
	if len(a) != len(b) {
		panic("buffers must have equal length")
	}

	for i := 0; i < len(a); i++ {
		xor := a[i] ^ b[i]

		for j := uint(0); j < 8; j++ {
			// Count how many bits are '1'.
			if xor&(1<<j) > 0 {
				dist += 1
			}
		}
	}

	return dist
}

func SplitBytes(in []byte, size int) (out [][]byte) {
	count := len(in) / size

	for i := 0; i < count; i++ {
		j := i * size
		block := in[j:(j + size)]
		out = append(out, block)
	}

	return out
}

// All blocks must have equal length.
func TransposeBlocks(blocks [][]byte) [][]byte {
	if len(blocks) == 0 {
		return nil
	}

	length := len(blocks[0])
	out := make([][]byte, length)

	for i := range blocks {
		block := blocks[i]

		for j := 0; j < length; j++ {
			out[j] = append(out[j], block[j])
		}
	}

	return out
}

func FindRepeatingXORKey(blocks [][]byte) (out []byte) {
	for _, block := range blocks {
		_, key := DetectSingleByteXOR(block)
		out = append(out, key)
	}

	return out
}

func RepeatingKeyXOR(plaintext, key []byte) []byte {
	keyLen := len(key)
	out := make([]byte, len(plaintext))
	for i := range plaintext {
		out[i] = plaintext[i] ^ key[i%keyLen]
	}
	return out
}
