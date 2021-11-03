package names

import (
	"fmt"
	"github.com/dustinkirkland/golang-petname"
	"math/rand"
	"time"
)

const seed = "0123456789"

func GetRandomName() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return fmt.Sprintf("%s-%s", petname.Generate(2, "-"), randSuffix(4))
}

func randSuffix(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = seed[rand.Intn(len(seed))]
	}
	return string(b)
}
