package stringstream


import (
  "errors"
)

type StringStream struct {
  str []rune
  curChar int
  length int
}

func MakeStringStream(str string) *StringStream {
  return &StringStream {
    str: []rune(str),
    curChar: 0,
    length: len(str),
  }
}

func (s *StringStream) PeekChar() (rune, error) {
  if s.curChar == s.length {
    return ' ', errors.New("Stream is empty") 
  }
  return  s.str[s.curChar], nil 
}

func (s *StringStream) ReadChar() (rune, error) {
  tmp, err := s.PeekChar()
  if err == nil {
    s.curChar++
  }
  return tmp, err 
}

func (s *StringStream) Empty() bool {
  return s.curChar == s.length
}
