package main

import (
	"fmt"
	"github.com/t1d333/vk_edu_golang/calc/internal/calc"
	"os"
)

func main() {
	expression := ""
	fmt.Scanln(&expression)
	result, err := calc.Calc(expression)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(result)
}
