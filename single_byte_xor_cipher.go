package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	input := HexToBytes(os.Args[1])
	output, _ := DetectSingleByteXOR(input)
	fmt.Println(string(output))
}
