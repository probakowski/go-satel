package satel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformCode(t *testing.T) {
	assert := assert.New(t)
	assert.Equal([]byte{0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, transformCode("0"))
	assert.Equal([]byte{0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, transformCode("0000"))
	assert.Equal([]byte{0x00, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, transformCode("000"))
	assert.Equal([]byte{0x12, 0x34, 0x56, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, transformCode("123456"))
	assert.Equal([]byte{0x98, 0x12, 0x4F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, transformCode("98124"))
	assert.Equal([]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0xFF}, transformCode("12345678901234"))
}
