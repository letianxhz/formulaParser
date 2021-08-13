package main

import (
	"bufio"
	"fmt"
	"formulaParser/formula"
	"os"
	"strings"
	"time"
)

func main() {
	for {
		fmt.Print("input /> ")
		f := bufio.NewReader(os.Stdin)
		s, err := f.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if s == "exit" || s == "quit" || s == "q" {
			fmt.Println("bye")
			break
		}
		start := time.Now()
		formula.Exec(s)
		cost := time.Since(start)
		fmt.Println("time: " + cost.String())
	}
}



