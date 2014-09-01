package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	. "github.com/sanderhahn/goexpr/eval"
)

const prompt = "> "

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for scanner.Scan() {
		val, err := Eval(scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("%g\n", val)
		}
		fmt.Print(prompt)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
