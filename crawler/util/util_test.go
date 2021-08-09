package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentBlock(t *testing.T) {
	block1 := CurrentBlock()
	assert.Greater(t, block1, int64(0))
	time.Sleep(12 * time.Second)
	block2 := CurrentBlock()
	assert.Equal(t, block1, block2-1)
}
