package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sanderhahn/goexpr/eval"
)

const prompt = "> "

func main() {
	env := eval.NewEnvironment()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for scanner.Scan() {
		value, err := env.Eval(scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%g\n", value)
		}
		fmt.Print(prompt)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
