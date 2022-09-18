package quotes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	n1 := Next()
	n2 := Next()

	assert.NotEmpty(t, n1)
	assert.NotEmpty(t, n2)
	assert.NotEqual(t, n1, n2)
}
