package reader

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReaderWithEmptyInput(t *testing.T) {
	reader := bytes.NewBufferString("")
	got, err := ReadLines(reader)
	require.Equal(t, nil, err, "Error occurred")
	require.Equal(t, 0, len(got))
}

func TestReaderWithOneLine(t *testing.T) {
	str := "test tring"
	reader := bytes.NewBufferString(str)
	expected := []string{str}
	got, err := ReadLines(reader)

	if err != nil {
		t.Error("Error occurred", err)
	}

	if got[0] != str {
		t.Errorf("\nExpected: %v\nGot: %v\n", expected, got)
	}
}

func TestReaderWithSomeLines(t *testing.T) {
	str := "Test1\nTest2\nTest3"
	reader := bytes.NewBufferString(str)
	expected := []string{"Test1", "Test2", "Test3"}
	got, err := ReadLines(reader)
	require.Equal(t, nil, err, "Error occurred")
	require.Equal(t, expected, got)
}
