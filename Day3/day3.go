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

func findAndApplyMul(s string) int {
	exp, err := regexp.Compile(`mul\([0-9]{1,3},[0-9]{1,3}\)`)
	if err != nil {
		fmt.Printf("Faulty regular expression\n")
		panic(err)
	}
	matches := exp.FindAllString(s, -1)
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
	return sum
}

func main() {
	fmt.Printf("Hi Day3\n")
	input := parseInput("./input.txt")
	sum := findAndApplyMul(input)
	fmt.Printf("Part 1 result: %v\n", sum)
	doStrings := strings.Split(input, "do()")
	sum = 0
	for i := 0; i < len(doStrings); i++ {
		splits := strings.Split(doStrings[i], "don't()")
		sum += findAndApplyMul(splits[0])
	}
	fmt.Printf("Part 2 result: %v\n", sum)
}
