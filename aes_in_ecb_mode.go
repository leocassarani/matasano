package main

import (
	"bytes"
	"crypto/aes"
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

	in := Base64ToBytes(string(file))
	out := DecryptAES_ECB(in, key)
	fmt.Println(string(out))
}

func DecryptAES_ECB(cipher, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	size := block.BlockSize()
	var plain [][]byte

	for i := 0; i < len(cipher)/size; i++ {
		ps := make([]byte, size)
		cs := cipher[(size * i):(size * (i + 1))]
		block.Decrypt(ps, cs)
		plain = append(plain, ps)
	}

	return bytes.Join(plain, nil)
}
