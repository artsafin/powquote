package quotes

import "math/rand"

func init() {
	rand.Seed(42)
}

func Next() string {
	index := rand.Intn(len(quotes))
	return quotes[index]
}
