package stream


import (
  "errors"
)

type Stream struct {
  str string
  curChar int
  length int
}

func MakeStream(str string) *Stream {
  return &Stream {
    str: str,
    curChar: 0,
    length: len(str),
  }
}

func (s *Stream) PeekChar() (byte, error) {
  if s.curChar == s.length {
    return ' ', errors.New("Stream is empty") 
  }
  return  s.str[s.curChar], nil 
}

func (s *Stream) ReadChar() (byte, error) {
  tmp, err := s.PeekChar()
  if err == nil {
    s.curChar++
  }
  return tmp, err 
}

func (s *Stream) Empty() bool {
  return s.curChar == s.length
}
