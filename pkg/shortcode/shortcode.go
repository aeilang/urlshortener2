package shortcode

import (
	"math/rand"

	"github.com/aeilang/urlshortener/config"
)

type ShortCodeGenerator struct {
	minLength int
}

func NewShortCodeGenerator(cfg config.ShortCodeConfig) *ShortCodeGenerator {
	return &ShortCodeGenerator{
		minLength: cfg.MinLength,
	}
}

const chars = "abcdefghijklmnopqrstuvwsyzABCDEFJHIJKLMNOKPRSTUVWSVZ0123456789"

func (s *ShortCodeGenerator) GenerateID() string {
	length := len(chars)
	shortCode := make([]byte, s.minLength)

	for i := 0; i < s.minLength; i++ {
		shortCode[i] = chars[rand.Intn(length)]
	}

	return string(shortCode)
}
