package generator

import (
	"math/rand"
	"strings"
)

type String struct {
	len int
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewString(len int) (*String, error) {
	return &String{len: len}, nil
}

func (g *String) Generate() interface{} {
	b := make([]rune, g.len)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return strings.ToUpper(string(b))
}
