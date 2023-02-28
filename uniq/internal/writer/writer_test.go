package writer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriterWithEmptySlice(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	lines := make([]string, 0)
	err := WriteLines(buf, lines)

	require.Equal(t, nil, err, "Error occurred")
	require.Equal(t, "", buf.String())
}

func TestWriterWithOneLine(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	lines := []string{"Test"}
	err := WriteLines(buf, lines)
	require.Equal(t, nil, err, "Error occurred")
	require.Equal(t, "Test\n", buf.String())
}

func TestWriterWithSomeLines(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	expected := "Test1\nTest2\nTest3\n"
	lines := []string{"Test1", "Test2", "Test3"}
	err := WriteLines(buf, lines)
	require.Equal(t, nil, err, "Error occurred")
	require.Equal(t, expected, buf.String())
}
