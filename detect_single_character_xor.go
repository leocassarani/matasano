package main

import (
	"bufio"
	"fmt"
	"os"
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
		out, _ := DetectSingleByteXOR(line)
		score := ScoreEnglishPlaintext(out)
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
