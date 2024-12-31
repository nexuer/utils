package rand

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

var (
	letterIdxBits int64 = 6                    // 6 bits to represent a letter index
	letterIdxMask int64 = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  int64 = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var rng = struct {
	sync.Mutex
	rand        *rand.Rand
	letterBytes string
}{
	rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	letterBytes: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
}

func Letter(charset string) {
	rng.Lock()
	defer rng.Unlock()

	if len(charset) < 10 || len(charset) > 64 {
		panic(fmt.Sprintf("rand: letter must be between 10 and 64 characters, got %d", len(charset)))
	}

	computed := int64(math.Ceil(math.Log2(float64(len(charset)))))

	rng.letterBytes = charset
	letterIdxBits = computed
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax = 63 / letterIdxBits
}

// Seed seeds the rng with the provided seed.
func Seed(seed int64) {
	rng.Lock()
	defer rng.Unlock()

	rng.rand = rand.New(rand.NewSource(seed))
}

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n)
// from the default Source.
func Perm(n int) []int {
	rng.Lock()
	defer rng.Unlock()
	return rng.rand.Perm(n)
}

// Int returns a non-negative pseudo-random int.
func Int() int {
	rng.Lock()
	defer rng.Unlock()
	return rng.rand.Int()
}

// Intn generates an integer in range [0,max).
// By design this should panic if input is invalid, <= 0.
func Intn(max int) int {
	rng.Lock()
	defer rng.Unlock()
	return rng.rand.Intn(max)
}

// IntnRange generates an integer in range [min,max).
// By design this should panic if input is invalid, <= 0.
func IntnRange(min, max int) int {
	rng.Lock()
	defer rng.Unlock()
	return rng.rand.Intn(max-min) + min
}

// Int63nRange generates an int64 integer in range [min,max).
// By design this should panic if input is invalid, <= 0.
func Int63nRange(min, max int64) int64 {
	rng.Lock()
	defer rng.Unlock()
	return rng.rand.Int63n(max-min) + min
}

func String(n int) string {
	b := make([]byte, n)
	rng.Lock()
	defer rng.Unlock()

	randomInt63 := rng.rand.Int63()
	remaining := letterIdxMax
	for i := 0; i < n; {
		if remaining == 0 {
			randomInt63, remaining = rng.rand.Int63(), letterIdxMax
		}
		if idx := int(randomInt63 & letterIdxMask); idx < len(rng.letterBytes) {
			b[i] = rng.letterBytes[idx]
			i++
		}
		randomInt63 >>= letterIdxBits
		remaining--
	}
	return string(b)
}

// SafeEncodeString encodes s using the same characters as rand.String. This reduces the chances of bad words and
// ensures that strings generated from hash functions appear consistent throughout the API.
func SafeEncodeString(s string) string {
	r := make([]byte, len(s))
	for i, b := range []rune(s) {
		r[i] = rng.letterBytes[(int(b) % len(rng.letterBytes))]
	}
	return string(r)
}
