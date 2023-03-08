package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/t1d333/vk_edu_golang/uniq/internal/reader"
	"github.com/t1d333/vk_edu_golang/uniq/internal/uniq"
	"github.com/t1d333/vk_edu_golang/uniq/internal/writer"
	"os"
)

func ParseFlags() (uniq.Options, []string) {
	var opt uniq.Options
	flag.BoolVar(&opt.ShowCount, "c", false, "Count the number of occurrences of a string in the input")
	flag.BoolVar(&opt.RepeatedOnly, "d", false, "Output only those lines that are repeated in the input.")
	flag.BoolVar(&opt.UniqOnly, "u", false, "Output only those lines that are not repeated in the input.")
	flag.BoolVar(&opt.IgnoreLetterCase, "i", false, "Ignore letter case.")
	flag.IntVar(&opt.FieldsNum, "f", 0, "Ignore the first num_fields of fields in the string.")
	flag.IntVar(&opt.CharsNum, "s", 0, "Ignore the first num_chars characters in the string. When used with the -f option, the first characters after num_fields fields are counted (ignoring the delimiter space after the last field).")
	flag.Parse()
	return opt, flag.CommandLine.Args()
}

func OpenInputOutput(files []string) (*os.File, *os.File, error) {
	var (
		err           error = nil
		filesLen      int   = len(files)
		input, output       = os.Stdin, os.Stdout
	)
	switch {
	case filesLen == 1:
		input, err = os.Open(files[0])
	case filesLen == 2:
		input, err = os.Open(files[0])
		if err != nil {
			return input, output, err
		}
		output, err = os.Create(files[1])
	case filesLen > 2:
		err = errors.New(fmt.Sprintf("extra operand ‘%s’", files[2]))
	}
	return input, output, err
}

func main() {

	options, files := ParseFlags()

	if err := options.IsValid(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	input, output, err := OpenInputOutput(files)
	defer input.Close()
	defer output.Close()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	lines, err := reader.ReadLines(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	lines = uniq.Uniq(lines, options)

	if err := writer.WriteLines(output, lines); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
