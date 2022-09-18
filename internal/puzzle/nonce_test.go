package puzzle

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNonceGenerator(t *testing.T) {
	gen := NewNonceGenerator(time.Millisecond * 100)
	go gen.Start(context.Background())

	v1 := gen.Current()
	v2 := gen.Current()

	assert.NotEqual(t, v1, 0, "v1 != 0")
	assert.Equal(t, v1, v2, "v1 == v2")

	time.Sleep(time.Millisecond * 300)

	v4 := gen.Current()
	v5 := gen.Current()

	assert.NotEqual(t, v4, 0, "v4 != 0")
	assert.NotEqual(t, v1, v4, "v1 != v4")
	assert.Equal(t, v4, v5, "v4 == v5")
}

func TestGenerateNonceOnce(t *testing.T) {
	v1 := GenerateNonceOnce()
	assert.NotEqual(t, 0, v1)

	v2 := GenerateNonceOnce()
	assert.NotEqual(t, v1, v2)
}
