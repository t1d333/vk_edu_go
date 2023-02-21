package calc

import (
	"errors"
	"strconv"
	"unicode"
	"github.com/t1d333/vk_edu_golang/calc/internal/stream"
)


//https://github.com/bmstu-iu9/scheme-labs/blob/master/lect13.md
// Грамматика
// E → T | E + T | E - T
// T → F | T * F | T / F
// F → N | -N | ( E )

func makeOperation(lhs, rhs int, op byte) (int, error) {
  switch op {
    case '*':
      return lhs * rhs, nil
    case '/':
      if rhs == 0 {
        return 0, errors.New("Division by zero")
      }
      return lhs / rhs, nil
    case '+':
      return lhs + rhs, nil
    case '-':
      return lhs - rhs, nil
    default:
      return 0, errors.New("Unknown operation")
  }
}

func EvalExpr(s *stream.Stream) (int, error) {
  expr, err := EvalTerm(s)
  for !s.Empty() && (err == nil) {
    char, _ := s.ReadChar()
    if (char == '+') || (char == '-') {
     term, err := EvalTerm(s)
     if err != nil {
       return 0, err
     }
     expr, _ = makeOperation(expr, term, char)
    } else {
      return 0, errors.New("Bad Expression")
    }
  }
  return expr, err
}

func EvalTerm(s *stream.Stream) (int, error) {
  term, err := EvalFactor(s)
  for !s.Empty() && (err == nil) {
    char, _ := s.PeekChar()
    if char == '*' || char == '/' {
      s.ReadChar()
      factor, err := EvalFactor(s)
      if err != nil {
        return 0, err
      }
      term, err = makeOperation(term, factor, char)
      if err != nil {
        return 0, err
      }
    } else {
      break
    }
  }
  return term, err
}


func EvalFactor(s *stream.Stream) (int, error) {
  for !s.Empty() {
    char, _ := s.PeekChar()
    switch {
    case unicode.IsNumber(rune(char)):
      return EvalNumber(s)
    case char == '-': 
      return EvalNumber(s)
    case char == '(':
      expr, err := ExtractExpr(s)
      if err != nil {
        return 0, err
      }
      a, b := EvalExpr(expr)
      return a, b
    default:
      return 0, errors.New("Bad Expression")
    }
  } 
  return 0, errors.New("Bad Expression")
}


func EvalNumber(s *stream.Stream) (int, error) {
  tmp := make([]byte, 0)
  sign := 1
  if char, _ := s.PeekChar(); char == '-' {
    sign = -1
    s.ReadChar()
  }
  
  for char, _ := s.PeekChar(); unicode.IsDigit(rune(char)) && !s.Empty(); {
    char, _ = s.ReadChar()
    tmp = append(tmp, char) 
    char, _ = s.PeekChar()
  }
  num, _ := strconv.Atoi(string(tmp))
  return sign * num, nil
}

func ExtractExpr(s *stream.Stream) (*stream.Stream, error) {
  expr := make([]byte, 0)
  openedBracketsCounter := 1 
  s.ReadChar()
  for !s.Empty() && (openedBracketsCounter != 0) {
    char, _ := s.PeekChar()
    switch char {
    case '(' :
        openedBracketsCounter++    
    case ')':
        openedBracketsCounter--
    }
    s.ReadChar()
    if openedBracketsCounter == 0 {
      break
    }
    expr = append(expr, char)
  }
  if openedBracketsCounter != 0 {
    return nil, errors.New("Bad expression")
  }
  return stream.MakeStream(string(expr)), nil
}

func Calc(expression string) (int, error) {
  s := stream.MakeStream(expression)
  return EvalExpr(s)
}

