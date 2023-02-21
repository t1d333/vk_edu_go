package uniq

import (
  "testing"
	"github.com/stretchr/testify/require"
)

var simpleTests = map[string]struct {
  input []string
  output []string
} {
  "Empty slice": {
    input: []string{},
    output: []string{},
  },

  "One line": {
    input: []string{"Test line"},
    output: []string{"Test line"},
  },
  "Some unique lines": {
    input: []string{"Test1","Test2", "Test3"}, 
    output: []string{"Test1", "Test2", "Test3"},
  },

  "Only repeated lines": {
    input: []string{"Test", "Test", "Test1", "Test1", "Test2", "Test2"}, 
    output: []string{"Test", "Test1", "Test2"},
  },

  "Repeated and unique lines": {
    input: []string{"Test", "Test", "unique", "", "123", "123", "Test"},
    output: []string{"Test", "unique", "", "123", "Test"},
  },

}

var testsWithCFlag = map[string]struct {
  input []string
  output []string
} {
  "One line": {
    input: []string{"Test line"},
    output: []string{"1 Test line"},
  },
  "Some unique lines": {
    input: []string{"Test1","Test2", "Test3"}, 
    output: []string{"1 Test1", "1 Test2", "1 Test3"},
  },

  "Only repeated lines": {
    input: []string{"Test", "Test", "Test1", "Test1", "Test2", "Test2"}, 
    output: []string{"2 Test", "2 Test1", "2 Test2"},
  },

  "Repeated and unique lines": {
    input: []string{"Test", "Test", "unique", "", "123", "123", "Test"},
    output: []string{"2 Test", "1 unique", "1 ", "2 123", "1 Test"},
  },
}


var testsWithDFlag = map[string]struct {
  input []string
  output []string
} {
  "One line": {
    input: []string{"Test line"},
    output: []string{},
  },
  "Some unique lines": {
    input: []string{"Test1","Test2", "Test3"}, 
    output: []string{},
  },

  "Only repeated lines": {
    input: []string{"Test", "Test", "Test1", "Test1", "Test2", "Test2"}, 
    output: []string{"Test", "Test1", "Test2"},
  },

  "Repeated and unique lines": {
    input: []string{"Test", "Test", "unique", "", "123", "123", "Test"},
    output: []string{"Test", "123"},
  },
}


var testsWithUFlag = map[string]struct {
  input []string
  output []string
} {
  "One line": {
    input: []string{"Test line"},
    output: []string{"Test line"},
  },
  "Some unique lines": {
    input: []string{"Test1","Test2", "Test3"}, 
    output: []string{"Test1","Test2", "Test3"},
  },

  "Only repeated lines": {
    input: []string{"Test", "Test", "Test1", "Test1", "Test2", "Test2"}, 
    output: []string{},
  },

  "Repeated and unique lines": {
    input: []string{"Test", "Test", "unique", "", "123", "123", "Test"},
    output: []string{"unique", "", "Test"},
  },
}

var testsWithFFlag = map[string]struct {
  FieldsNum int
  input []string
  output []string
} {
  "With one field": {
    FieldsNum: 1,
    input: []string{"a Test", "b Test", "c Test"},
    output: []string{"a Test"},
  },

  "With some fields": {
    FieldsNum: 2,
    input: []string{"Test123 abc", "Test123", "Unique", "Unique Test Test2"},
    output: []string{"Test123 abc", "Unique Test Test2"},
  },
}

var testsWithSFlag = map[string]struct {
  CharsNum int
  input []string
  output []string
} {
  "With one char": {
    CharsNum: 1,
    input: []string{"abc Test", "bbc Test", "cbc TestTest"},
    output: []string{"abc Test", "cbc TestTest"},
  },

  "With some chars": {
    CharsNum: 3,
    input: []string{"321Test", "123Test", "Unique", "TesTestTest", "UniTestTest"},
    output: []string{"321Test", "Unique", "TesTestTest"},
  },
}


func TestSimple(t *testing.T) {
  var opt Options
  for name, test := range simpleTests {
    got, _ := Uniq(test.input, opt)
    require.Equal(t, test.output, got, name)  
  }
}

func TestCFlag(t *testing.T) {
  var opt Options = Options{ShowCount: true}

  for name, test := range testsWithCFlag {
    got, _ := Uniq(test.input, opt)
    require.Equal(t, test.output, got, name)  
  }
}

func TestDFlag(t *testing.T) {
  var opt Options = Options{RepeatedOnly: true}
  for name, test := range testsWithDFlag {
    got, _ := Uniq(test.input, opt)
    require.Equal(t, test.output, got, name)  
  }
}

func TestUFlag(t *testing.T) {
  var opt Options = Options{UniqOnly: true}
  for name, test := range testsWithUFlag {
    got, _ := Uniq(test.input, opt)
    require.Equal(t, test.output, got, name)  
  }
}

func TestFFlag(t *testing.T) {
  for name, test := range testsWithFFlag {
    var opt Options = Options{FieldsNum: test.FieldsNum}
    got, _ := Uniq(test.input, opt)
    require.Equal(t, test.output, got, name)  
  }
}

func TestSFlag(t *testing.T) {
  for name, test := range testsWithSFlag {
    var opt Options = Options{CharsNum: test.CharsNum}
    got, _ := Uniq(test.input, opt)
    require.Equal(t, test.output, got, name)  
  }
}

func TestWithIncorrectFieldsNum(t *testing.T) {
  input := []string{"test"}
  opt := Options{FieldsNum: -10}
  _, err := Uniq(input, opt)
  require.Error(t, err)
}

func TestWithIncorrectCharsNum(t *testing.T) {
  input := []string{"test"}
  opt := Options{CharsNum: -10}
  _, err := Uniq(input, opt)
  require.Error(t, err)
}

func TestWithFlagsCAndD(t *testing.T) {
 input := []string{"test"}
  opt := Options{
    RepeatedOnly: true,
    ShowCount: true,
  }
  _, err := Uniq(input, opt)
  require.Error(t, err)
}

func TestWithFlagsCAndU(t *testing.T) {
 input := []string{"test"}
  opt := Options{
    UniqOnly: true,
    ShowCount: true,
  }
  _, err := Uniq(input, opt)
  require.Error(t, err)
}

func TestWithFlagsDAndU(t *testing.T) {
 input := []string{"test"}
  opt := Options{
    UniqOnly: true,
    RepeatedOnly: true,
  }
  _, err := Uniq(input, opt)
  require.Error(t, err)
}
