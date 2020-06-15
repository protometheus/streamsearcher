package main

import (
	"io/ioutil"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateInput creates some randomized input with the searchTerm
// sprinkled in pseudo-randomly
func GenerateInput(filename, searchTerm string, length int64) {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	// pick two random numbers and combine them
	// add 53 to make sure it's not too small
	p := int64(randomizer.Intn(1500) + 53)
	q := int64(randomizer.Intn(1500) + 53)
	pq := p * q

	needle := []byte(searchTerm)
	needleLen := int64(len(needle))

	b := make([]byte, length)
	var i int64
	for i = 0; i < length; i++ {
		if i != 0 && i%(pq) == 0 && i+needleLen < length {
			var j int64
			for j = 0; j < needleLen; {
				b[i] = needle[j]
				i++
				j++
			}
			i = i + needleLen
		} else if i != 0 && i%500 == 0 {
			// add some newlines for readability
			b[i] = byte('\n')
		} else {
			b[i] = charset[randomizer.Intn(len(charset))]
		}
	}

	ioutil.WriteFile(filename, b, 0644)
}
