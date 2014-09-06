package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	key := []byte(os.Args[2])
	out := RepeatingKeyXOR(file, key)
	fmt.Println(BytesToHex(out))
}

func RepeatingKeyXOR(plaintext, key []byte) []byte {
	keyLen := len(key)
	out := make([]byte, len(plaintext))
	for i := range plaintext {
		out[i] = plaintext[i] ^ key[i % keyLen]
	}
	return out
}
