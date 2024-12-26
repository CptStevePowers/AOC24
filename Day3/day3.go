package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseInput(p string) string {
	bytes, err := os.ReadFile(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	text := string(bytes)
	return text
}

func main() {
	fmt.Printf("Hi Day3\n")
	input := parseInput("./input.txt")
	exp, err := regexp.Compile(`mul\([0-9]{1,3},[0-9]{1,3}\)`)
	if err != nil {
		fmt.Printf("Faulty regular expression\n")
		panic(err)
	}
	matches := exp.FindAllString(input, -1)
	fmt.Printf("Found %v matches\n", len(matches))
	sum := 0
	for i := 0; i < len(matches); i++ {
		s := matches[i]
		s, _ = strings.CutPrefix(s, "mul(")
		s, _ = strings.CutSuffix(s, ")")
		splits := strings.Split(s, ",")
		a, err := strconv.Atoi(splits[0])
		if err != nil {
			fmt.Printf("Something went wrong when parsing %s", splits[0])
			panic(err)
		}
		b, err := strconv.Atoi(splits[1])
		if err != nil {
			fmt.Printf("Something went wrong when parsing %s", splits[1])
			panic(err)
		}
		sum += a * b
	}
	fmt.Printf("Result: %v", sum)

}
