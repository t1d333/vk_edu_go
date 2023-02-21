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

func checkOptions(opt Options) error {
  if (opt.ShowCount && opt.RepeatedOnly) || (opt.ShowCount && opt.UniqOnly) || (opt.RepeatedOnly && opt.UniqOnly) {
    return errors.New("Usage: go run main.go [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
  }

  if opt.FieldsNum < 0 {
    return errors.New(fmt.Sprintf("%d: %s", opt.FieldsNum, "invalid number of fields to skip"))
  }

  if opt.CharsNum < 0 {
    return errors.New(fmt.Sprintf("%d: %s", opt.CharsNum, "invalid number of bytes to skip"))
  }

  return nil
}

func trimFields(str string, numFields int) string { 
  fields := strings.Split(str, " ") 
  for index := range fields {
    if index == numFields {
      break
    }
    fields = fields[1:]
  }

  return strings.Join(fields, " ") 
}

func trimChars(str string, numChars int) string {
  for index := range str {
    if index == numChars {
      break
    }
    str = str[1:]
  }
  return str
}

func equalStrings(str1, str2 string, IgnoreLetterCase bool) bool {
  if IgnoreLetterCase { 
    return strings.EqualFold(str1, str2)
  }
  return strings.Compare(str1, str2) == 0
}

func Uniq(lines []string, opt Options) ([]string, error) {
  if err := checkOptions(opt); err != nil {
    return lines, err
  }

  var (
    i, j = 0, 0
    linesCount = len(lines)
    result = make([]string, 0)
  )

  for i = 0; i < linesCount; i = j {
    for j = i; (j < linesCount); j++ {
      linePrev, lineNext := trimFields(lines[i], opt.FieldsNum), trimFields(lines[j], opt.FieldsNum)
      linePrev, lineNext = trimChars(linePrev, opt.CharsNum), trimChars(lineNext, opt.CharsNum)
      if equalStrings(linePrev, lineNext, opt.IgnoreLetterCase) {
        continue
      }
      break
    }
    if (opt.RepeatedOnly && ((j - i) == 1)) || (opt.UniqOnly && ((j - i) > 1)) {
      continue
    }
    line := lines[i]
    if opt.ShowCount {
      line = fmt.Sprintf("%d %s", j - i, line)
    }
    result = append(result, line)
  }
  return result, nil
}
