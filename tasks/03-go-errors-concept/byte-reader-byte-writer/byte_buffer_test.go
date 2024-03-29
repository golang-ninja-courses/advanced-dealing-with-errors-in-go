package bytebuffer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByteBufferImplementsNecessary(t *testing.T) {
	assert.Implements(t, (*io.ByteWriter)(nil), new(ByteBuffer))
	assert.Implements(t, (*io.ByteReader)(nil), new(ByteBuffer))
}

func TestByteBufferIOReader(t *testing.T) {
	var b ByteBuffer

	expected := "TestByteBuffer_IOReader"
	var actual strings.Builder

	for _, c := range []byte(expected) {
		err := b.WriteByte(c)
		require.NoError(t, err)
	}

	for i := 0; i < len(expected)+1; i++ {
		bb, err := b.ReadByte()
		if err != nil {
			if isEndOfBuffer(err) {
				break
			}
			require.NoError(t, err)
		}

		actual.WriteByte(bb)
	}

	assert.Equal(t, expected, actual.String())
}

func TestByteBufferIOWriter(t *testing.T) {
	var b ByteBuffer

	for i := 0; i < bufferMaxSize+1; i++ {
		err := b.WriteByte('1')
		if err != nil {
			if isMaxSizeExceededError(err) {
				break
			}
			require.NoError(t, err)
		}
	}

	assert.Len(t, b.buffer, bufferMaxSize)
}

func TestByteBufferReadFromEmptyBuffer(t *testing.T) {
	var b ByteBuffer
	n, err := b.ReadByte()
	assert.EqualValues(t, 0, n)
	assert.True(t, isEndOfBuffer(err))
}

func TestEndOfBufferError(t *testing.T) {
	assert.NotEmpty(t, new(EndOfBufferError).Error())
}

func TestMaxSizeExceededError(t *testing.T) {
	assert.NotEmpty(t, new(MaxSizeExceededError).Error())
}

func isEndOfBuffer(err error) bool {
	_, ok := err.(*EndOfBufferError)
	return ok
}

func isMaxSizeExceededError(err error) bool {
	_, ok := err.(*MaxSizeExceededError)
	return ok
}
