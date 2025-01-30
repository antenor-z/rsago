package main

import (
	"bufio"
	"fmt"
	"os"
	"rsago/rsa"
)

func main() {
	pub, prv := rsa.DeriveKeys()

	print("------------------\n")
	original := ""
	fmt.Printf("Original: ")
	in := bufio.NewReader(os.Stdin)
	original, err := in.ReadString('\n')
	if err != nil {
		panic("Error on ReadString()")
	}
	encoded := rsa.EncodeString(original, pub)
	fmt.Printf("Encoded: %s\n", encoded)
	decoded := rsa.DecodeString(encoded, prv)
	fmt.Printf("Decoded: %s\n", decoded)
}
