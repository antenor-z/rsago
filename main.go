package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

func randomPrime() *big.Int {
	keySize := 1024
	prime, err := rand.Prime(rand.Reader, keySize)
	if err != nil {
		panic("Error generating prime")
	}
	return prime
}

// https://www.cs.utexas.edu/~mitra/honors/soln.html

type publicKey struct {
	e *big.Int
	n *big.Int
}

type privateKey struct {
	d *big.Int
	n *big.Int
}

func encodeBlock(block *big.Int, pub publicKey) *big.Int {
	return new(big.Int).Exp(block, pub.e, pub.n)
}

func decodeBlock(block *big.Int, pub privateKey) *big.Int {
	return new(big.Int).Exp(block, pub.d, pub.n)
}

func totient(p *big.Int, q *big.Int) *big.Int {
	one := big.NewInt(1)
	p1 := new(big.Int).Sub(p, one)
	q1 := new(big.Int).Sub(q, one)
	fmt.Printf("p-1=%d\n", p1)
	fmt.Printf("q-1=%d\n", q1)
	return new(big.Int).Mul(p1, q1)
}

func coprime(tot *big.Int) *big.Int {
	one := big.NewInt(1)
	e := big.NewInt(1<<16 + 1)
	gcd := new(big.Int)
	for e.Cmp(tot) < 0 { // e < tot
		gcd.GCD(nil, nil, tot, e)
		if gcd.Cmp(one) == 0 {
			return e
		}
		e.Add(e, one)
	}
	panic("Could not find coprime") // this should never happen
}

func computeD(tot *big.Int, e *big.Int) *big.Int {
	d := new(big.Int).ModInverse(e, tot)
	if d == nil {
		panic("No modular inverse exists; e and totient are not coprime")
	}
	return d
}

func deriveKeys(p *big.Int, q *big.Int) (publicKey, privateKey) {
	n := new(big.Int)
	n.Mul(p, q)
	fmt.Printf("n=pq=%d\n", n)

	tot := totient(p, q)
	fmt.Printf("totient=%d\n", tot)

	e := coprime(tot)
	fmt.Printf("e=%d\n", e)

	d := computeD(tot, e)
	fmt.Printf("d=%d\n", d)

	pub := publicKey{e, n}
	prv := privateKey{d, n}

	eBytes := pub.e.Bytes()
	nBytes := pub.n.Bytes()
	dBytes := prv.d.Bytes()

	eB64 := base64.StdEncoding.EncodeToString(eBytes)
	nB64 := base64.StdEncoding.EncodeToString(nBytes)
	dB64 := base64.StdEncoding.EncodeToString(dBytes)

	fmt.Printf("pub:\ne=%s\nn=%s\n", eB64, nB64)
	fmt.Printf("prv:\nd=%s\nn=%s\n", dB64, nB64)

	return pub, prv
}

func encodeString(s string, pub publicKey) string {
	bs := []byte(s)
	block := new(big.Int).SetBytes(bs)
	encBlock := encodeBlock(block, pub).Bytes()
	return string(encBlock)
}

func decodeString(s string, prv privateKey) string {
	bs := []byte(s)
	block := new(big.Int).SetBytes(bs)
	decBlock := decodeBlock(block, prv).Bytes()
	return string(decBlock)
}

func main() {

	pub, prv := deriveKeys(randomPrime(), randomPrime())

	print("------------------\n")
	encoded := encodeString("As of Go 1.8, ProbablyPrime(0) is allowed and applies only a Baillie-PSW test. Before Go 1.8, ProbablyPrime applied only the Miller-Rabin tests, and ProbablyPrime(0) panicked. ", pub)
	decoded := decodeString(encoded, prv)

	print(decoded)
}
