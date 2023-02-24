package calc

import (
  "testing"
  "github.com/stretchr/testify/require"
)

var simpleTests = []struct {
  input string
  output int
} {
  {"1", 1},
  {"(((10)))", 10},
  {"-123", -123},
  {"+12", 12},
  {"1+2+3", 6},
  {"5-1", 4,},
  {"3*6", 18,},
  {"9/4", 2,},
  {"23+103", 126},
  {"23-103", -80},
  {"13*13", 169},
  {"1+2*3", 7},
  {"10/0", 0},
  {"-15/3", -5},
  {"(-10)-(+12)", -22},
  {"16/(-4)", -4},
  {"3*3+1", 10},
  {"33/10*2", 6},
  {"(3+10)+(2+12)", 27},
  {"((16+4)*2)/(1+1)", 20},
  {"(3+1)*(1+1)", 8},
  {"((((3+1)*2)+4)+10)*2", 44},
}

var badExpressions = []string {
   "1+2+",
   "*2",
   "12*",
   "()",
   "(1+2))",
   "(1+)",
   "((1*2)",
   "1+a",
   "10/0",
}

func TestSimple(t *testing.T) {
  for _, test := range simpleTests {
    got, _ := Calc(test.input)
    require.Equal(t, test.output, got)  
  }
}

func TestBadExpressions(t *testing.T) {
  for _, expr := range badExpressions {
    _, err := Calc(expr)
    require.Error(t, err)
  }
}

