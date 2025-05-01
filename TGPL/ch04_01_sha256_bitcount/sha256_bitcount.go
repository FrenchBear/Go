// sha256_bitcount.go
// Exercises from §4 of The Go Programming Language
// Count the number of ≠ bits between sha256 hashes of "x" and "X"
// Personal implem of SHA256, adapted from C# version (VS Proj 519), adapted from Wikipedia algorithms
//
// 2019-10-14	PV

package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	// Go version
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%v -> %x\n%v -> %x\n\n", "x", c1, "X", c2)

	// My version
	pv1 := mySha256("x")
	pv2 := mySha256("X")
	fmt.Printf("%v -> %x\n%v -> %x\n\n", "x", pv1, "X", pv2)

	// Result must be the same
	if c1 != pv1 {
		panic("c1≠pv1")
	}
	if c2 != pv2 {
		panic("c2≠pv2")
	}

	// More tests, panics if test fails
	testSha256()

	sum := 0
	for i := 0; i < len(c1); i++ {
		sum += countBitDiff(c1[i], c2[i])
	}
	fmt.Println("Bits:", 8*len(c1), " bits≠:", sum)
}

func testSha256() {
	// For tests, strings are only considered as composed of simple bytes (ASCII), not unicode characters

	// Empty string
	oneTestSha256("",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

	oneTestSha256("The quick brown fox jumps over the lazy dog",
		"d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592")

	// Just adding a character completely change the output
	oneTestSha256("The quick brown fox jumps over the lazy dog.",
		"ef537f25c895bfa782526529a9b63d97aa631564d5d789c2b765448c8635fb6c")

	// Two blocks processing
	oneTestSha256("abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq",
		"248d6a61d20638b8e5c026930c3e6039a33ce45964ff2167f6ecedd419db06c1")

	// 2019-20-14 Compare with Go
	oneTestSha256("x", "2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881")
	oneTestSha256("X", "4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015")
}

func oneTestSha256(s, hashed string) {
	c := mySha256(s)
	if hashed != fmt.Sprintf("%x", c) {
		panic("err in oneTestSha256")
	}
}

func countBitDiff(x, y byte) int {
	diffBits := ^(x ^ y)
	n := 0
	for diffBits != 0 {
		n++
		diffBits &= diffBits - 1
	}
	return n
}

// Example of implementation of SHA-2 protocol
// Notes:
// - This is just to understand SHA-2, not a performance/reference implementation, use Go provided version for that!
// - Length is only managed in bytes (multiple of 8 bits)
func mySha256(s string) [32]byte {
	// Initialize hash values:
	// First 32 bits of the fractional parts of the square roots of the first 8 primes 2..19:
	// h[0..7]
	h := [8]uint32{
		0x6a09e667,
		0xbb67ae85,
		0x3c6ef372,
		0xa54ff53a,
		0x510e527f,
		0x9b05688c,
		0x1f83d9ab,
		0x5be0cd19,
	}

	return sha256224(s, h)
}

func preprocessing(s string, blocksize int, lengthsize int) ([]byte, int) {
	// Pre-processing:
	// append the bit '1' to the message
	// append k bits '0', where k is the minimum number >= 0 such that the resulting message
	// length (modulo <blocksize> in bits) is <blocksize>-<length>.
	// append length of message (without the '1' bit or padding), _in bits_, as 64-bit big-endian integer
	// (this will make the entire post-processed length a multiple of 512 bits)

	lB := len(s)         // Message length in Bytes (only process ascii here, one byte per character)
	lb := 8 * lB         // Message length in bits
	nb := lb / blocksize // Number of blocks of 512 bits
	if lb == 0 || lb%blocksize != 0 {
		nb++
	}
	if (lb % blocksize) >= (blocksize - lengthsize) {
		nb++
	}

	tb := make([]byte, nb*(blocksize/8))
	for i := 0; i < lB; i++ {
		tb[i] = s[i]
	}
	j := lB
	tb[j] = 0x80
	j++
	for (j % (blocksize / 8)) != (blocksize-lengthsize)/8 {
		tb[j] = 0
		j++
	}
	// Length is always 32-bit in this implementation
	if lengthsize == 128 {
		for i := 0; i < 8; i++ {
			tb[j] = 0
			j++
		}
	}
	tb[j] = 0
	j++
	tb[j] = 0
	j++
	tb[j] = 0
	j++
	tb[j] = 0
	j++
	tb[j] = byte(lb >> 24)
	j++
	tb[j] = byte((lb & 0x00FF0000) >> 16)
	j++
	tb[j] = byte((lb & 0x0000FF00) >> 8)
	j++
	tb[j] = byte((lb & 0x000000FF))
	j++

	// Some assertions
	if !(j%(blocksize/8) == 0) {
		panic("err1")
	}
	if !(j == nb*(blocksize/8)) {
		panic("err2")
	}

	return tb, nb
}

// Note 1: All variables are 32 bit unsigned integers and addition is calculated modulo 2^32
// Note 2: For each round, there is one round constant k[i] and one entry in the message schedule array w[i], 0 ≤ i ≤ 63
// Note 3: The compression function uses 8 working variables, a through h
// Note 4: Big-endian convention is used when expressing the constants in this pseudocode,
// and when parsing message block data from bytes to words, for example,
// the first word of the input message "abc" after padding is 0x61626380
func sha256224(s string, h [8]uint32) [32]byte {
	// Initialize array of round constants:
	// (first 32 bits of the fractional parts of he cube roots of the first 64 primes 2..311):
	// k[0..63]
	k := [...]uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
	}

	//var tb []byte       // Table of bytes
	//var nb int          // Number of blocks
	tb, nb := preprocessing(s, 512, 64)

	// Process the message in successive 512-bit chunks:
	// break message into 512-bit chunks for each chunk

	// br is the bloc rank (number)
	for br := 0; br < nb; br++ {
		// create a 64-entry message schedule array w[0..63] of 32-bit words
		// (The initial values in w[0..63] don't matter, so many implementations zero them here)
		var w [64]uint32

		// copy chunk into first 16 words w[0..15] of the message schedule array
		for i := 0; i < 64; i += 4 {
			w[i>>2] = uint32(tb[(br<<6)+i])<<24 + uint32(tb[(br<<6)+i+1])<<16 + uint32(tb[(br<<6)+i+2])<<8 + uint32(tb[(br<<6)+i+3])
		}

		// Extend the first 16 words into the remaining 48 words w[16..63] of the message schedule array:
		// for i from 16 to 63
		//      s0 := (w[i-15] rightrotate 7) xor (w[i-15] rightrotate 18) xor (w[i-15] rightshift 3)
		//      s1 := (w[i-2] rightrotate 17) xor (w[i-2] rightrotate 19) xor (w[i-2] rightshift 10)
		//      w[i] := w[i-16] + s0 + w[i-7] + s1
		for i := 16; i < 64; i++ {
			s0 := rightrotate(w[i-15], 7) ^ rightrotate(w[i-15], 18) ^ (w[i-15] >> 3)
			s1 := rightrotate(w[i-2], 17) ^ rightrotate(w[i-2], 19) ^ (w[i-2] >> 10)
			w[i] = w[i-16] + s0 + w[i-7] + s1
		}

		// Initialize working variables to current hash value:
		a := h[0]
		b := h[1]
		c := h[2]
		d := h[3]
		e := h[4]
		f := h[5]
		g := h[6]
		hh := h[7]

		// Compression function main loop:
		// for i from 0 to 63
		for i := 0; i < 64; i++ {
			// S1 := (e rightrotate 6) xor (e rightrotate 11) xor (e rightrotate 25)
			// ch := (e and f) xor ((not e) and g)
			// temp1 := h + S1 + ch + k[i] + w[i]
			// S0 := (a rightrotate 2) xor (a rightrotate 13) xor (a rightrotate 22)
			// maj := (a and b) xor (a and c) xor (b and c)
			// temp2 := S0 + maj

			S1 := rightrotate(e, 6) ^ rightrotate(e, 11) ^ rightrotate(e, 25)
			ch := (e & f) ^ ((^e) & g)
			temp1 := hh + S1 + ch + k[i] + w[i]
			S0 := rightrotate(a, 2) ^ rightrotate(a, 13) ^ rightrotate(a, 22)
			maj := (a & b) ^ (a & c) ^ (b & c)
			temp2 := S0 + maj

			hh = g
			g = f
			f = e
			e = d + temp1
			d = c
			c = b
			b = a
			a = temp1 + temp2
		}

		// Add the compressed chunk to the current hash value:
		h[0] = h[0] + a
		h[1] = h[1] + b
		h[2] = h[2] + c
		h[3] = h[3] + d
		h[4] = h[4] + e
		h[5] = h[5] + f
		h[6] = h[6] + g
		h[7] = h[7] + hh
	}

	// Produce the final hash value (big-endian):
	// return h[0].ToString("x8") + h[1].ToString("x8") + h[2].ToString("x8") + h[3].ToString("x8") + h[4].ToString("x8") + h[5].ToString("x8") + h[6].ToString("x8") + h[7].ToString("x8")
	// In Go, we retirn an array of 32 bytes -> break uint32 into 4 bytes
	var res [32]byte
	i := 0
	for _, x := range h {
		res[i+3] = byte(x)
		x >>= 8
		res[i+2] = byte(x)
		x >>= 8
		res[i+1] = byte(x)
		x >>= 8
		res[i] = byte(x)
		i += 4
	}
	return res
}

// equivalent of C++ _rotr
// 32-bit version
func rightrotate(original uint32, bits uint) uint32 {
	return (original >> bits) | (original << (32 - bits))
}
