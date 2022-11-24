package generator

import (
	"math"
	"math/rand"
)

type Int8 struct {
}

func NewInt8() *Int8 {
	return &Int8{}
}

func (g *Int8) Generate() interface{} {
	return rand.Intn(math.MaxInt8)
}
