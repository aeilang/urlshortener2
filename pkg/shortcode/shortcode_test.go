package shortcode

import (
	"testing"

	"github.com/aeilang/urlshortener/config"
)

func TestGenerate(t *testing.T) {
	g := NewShortCodeGenerator(config.ShortCodeConfig{MinLength: 6})

	code := g.GenerateID()
	if len(code) == 6 {
		t.Error(code)
	}
}
