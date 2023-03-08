package stringstream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamCreation(t *testing.T) {
	tmp := MakeStringStream("test")
	require.Equal(t, []rune("test"), tmp.str)
	require.Equal(t, 4, tmp.length)
	require.Equal(t, 0, tmp.curChar)
	tmp.ReadChar()
	require.Equal(t, 1, tmp.curChar)
}

func TestStream(t *testing.T) {
	s := MakeStringStream("abcde")
	for i := 0; i < 5; i++ {
		got, err := s.ReadChar()
		require.NoError(t, err)
		require.Equal(t, int32('a'+i), got)
	}
	require.True(t, s.Empty())
	_, err := s.PeekChar()
	require.Error(t, err)
	_, err = s.ReadChar()
	require.Error(t, err)
}
