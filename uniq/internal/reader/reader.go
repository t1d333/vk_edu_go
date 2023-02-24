package reader

import (
	"bufio"
	"io"
)



func ReadLines(rd io.Reader) ([]string, error) {
  var (
    sc = bufio.NewScanner(rd)
    result  = make([]string, 0)
    err error = nil
  )

  for sc.Scan() {
    result = append(result, sc.Text())
  }
  return result, err
}
