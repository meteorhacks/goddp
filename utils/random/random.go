package random

import (
	"math/rand"
)

var dict = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Id(n int) string {
	a := make([]rune, n)
	l := len(dict)

	for i := range a {
		a[i] = dict[rand.Intn(l)]
	}

	return string(a)
}
