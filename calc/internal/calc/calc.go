package calc

import (
	"errors"
	"github.com/t1d333/vk_edu_golang/calc/internal/stream"
	"github.com/t1d333/vk_edu_golang/calc/internal/stringstream"
	"strconv"
	"unicode"
)

var (
	BadEexpressionError   = errors.New("bad expression")
	DivByZeroError        = errors.New("division by zero")
	UnknownOperationError = errors.New("unknown operation")
)

// https://github.com/bmstu-iu9/scheme-labs/blob/master/lect13.md
// Грамматика
// E → T | E + T | E - T
// T → F | T * F | T / F
// F → N | -N | ( E )

func makeOperation(lhs, rhs int, op rune) (int, error) {
	switch op {
	case '*':
		return lhs * rhs, nil
	case '/':
		if rhs == 0 {
			return 0, DivByZeroError
		}
		return lhs / rhs, nil
	case '+':
		return lhs + rhs, nil
	case '-':
		return lhs - rhs, nil
	default:
		return 0, UnknownOperationError
	}
}

func evalExpr(s stream.Stream) (int, error) {
	expr, err := evalTerm(s)

	for !s.Empty() && err == nil {
		char, _ := s.ReadChar()

		if char == '+' || char == '-' {
			term, err := evalTerm(s)
			if err != nil {
				return 0, err
			}
			expr, _ = makeOperation(expr, term, char)

		} else {
			return 0, BadEexpressionError
		}
	}

	return expr, err
}

func evalTerm(s stream.Stream) (int, error) {
	term, err := evalFactor(s)

	for !s.Empty() && err == nil {
		char, _ := s.PeekChar()

		if char == '*' || char == '/' {
			s.ReadChar()
			factor, err := evalFactor(s)

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

func evalFactor(s stream.Stream) (int, error) {
	char, _ := s.PeekChar()

	switch {
	case unicode.IsNumber(char):
		return evalNumber(s)
	case char == '-' || char == '+':
		return evalNumber(s)
	case char == '(':
		expr, err := extractExpr(s)
		if err != nil {
			return 0, err
		}
		return evalExpr(expr)
	default:
		return 0, BadEexpressionError
	}
}

func evalNumber(s stream.Stream) (int, error) {
	tmp := make([]rune, 0)
	sign := 1

	if char, _ := s.PeekChar(); char == '-' {
		sign = -1
		s.ReadChar()
	} else if char == '+' {
		s.ReadChar()
	}

	for char, _ := s.PeekChar(); unicode.IsDigit(char) && !s.Empty(); {
		char, _ = s.ReadChar()
		tmp = append(tmp, char)
		char, _ = s.PeekChar()
	}

	num, _ := strconv.Atoi(string(tmp))

	return sign * num, nil
}

func extractExpr(s stream.Stream) (stream.Stream, error) {
	expr := make([]rune, 0)
	openedBracketsCounter := 1
	s.ReadChar()

	for !s.Empty() && openedBracketsCounter != 0 {
		char, _ := s.ReadChar()

		switch char {
		case '(':
			openedBracketsCounter++
		case ')':
			openedBracketsCounter--
		}

		if openedBracketsCounter == 0 {
			break
		}

		expr = append(expr, char)
	}

	if openedBracketsCounter != 0 {
		return nil, BadEexpressionError
	}

	return stringstream.MakeStringStream(string(expr)), nil
}

func Calc(expression string) (int, error) {
	var s stream.Stream = stringstream.MakeStringStream(expression)
	return evalExpr(s)
}
