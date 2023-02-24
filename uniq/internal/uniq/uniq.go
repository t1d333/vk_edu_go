package uniq

import (
  "errors"
  "fmt"
  "strings"
)

type Options struct {
  ShowCount bool
  RepeatedOnly bool
  UniqOnly bool
  IgnoreLetterCase bool
  FieldsNum int
  CharsNum int
}

func boolToInt(b bool) int {
  if b {
    return 1
  }
  return 0
}

func (o *Options) IsValid() error {
  c, u, d := boolToInt(o.ShowCount), boolToInt(o.UniqOnly), boolToInt(o.RepeatedOnly)
  if c + u + d > 1 {
    return errors.New("Usage: go run main.go [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
  }

  if o.FieldsNum < 0 {
    return errors.New(fmt.Sprintf("%d: %s", o.FieldsNum, "invalid number of fields to skip"))
  }

  if o.CharsNum < 0 {
    return errors.New(fmt.Sprintf("%d: %s", o.CharsNum, "invalid number of bytes to skip"))
  }

  return nil
}

func trimFields(str string, numFields int) string { 
  fields := strings.Split(str, " ") 
  if numFields >= len(fields) {
    return ""
  } else {
    return strings.Join(fields[numFields:], " ") 
  }
}

func trimChars(str string, numChars int) string {
  if numChars >= len(str) {
    return ""
  } else {
    return str[numChars:]
  }
}


func equalStrings(str1, str2 string, opt Options) bool {
  str1 = trimChars(trimFields(str1, opt.FieldsNum), opt.CharsNum)
  str2 = trimChars(trimFields(str2, opt.FieldsNum), opt.CharsNum)
  if opt.IgnoreLetterCase { 
    return strings.EqualFold(str1, str2)
  }
  return strings.Compare(str1, str2) == 0
}

func Uniq(lines []string, opt Options) []string {
  var (
    i, j = 0, 0
    linesCount = len(lines)
    result = make([]string, 0)
  )

  for i = 0; i < linesCount; i = j {
    line := lines[i]
    for j = i; j < linesCount && equalStrings(lines[i], lines[j], opt); j++ {}

    if (opt.RepeatedOnly && ((j - i) == 1)) || (opt.UniqOnly && ((j - i) > 1)) {
      continue
    }

    if opt.ShowCount {
      line = fmt.Sprintf("%d %s", j - i, line)
    }

    result = append(result, line)
  }
  return result
}
