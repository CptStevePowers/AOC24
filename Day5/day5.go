package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func parseInput(p string) (rules map[int][]int, updates [][]int) {
	rules = make(map[int][]int)
	updates = make([][]int, 0, 8)
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()

	ruleExp, err := regexp.Compile(`^[\d]+\|[\d]+$`)
	if err != nil {
		fmt.Printf("Error compiling regular expression")
		panic(err)
	}
	updateExp, err := regexp.Compile(`^([\d]+,)+[\d]+$`)
	if err != nil {
		fmt.Printf("Error compiling regular expression")
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if ruleExp.MatchString(line) {
			splits := strings.Split(line, "|")
			n1, err := strconv.Atoi(splits[0])
			if err != nil {
				fmt.Printf("Could not parse %v", splits[0])
				panic(err)
			}

			n2, err := strconv.Atoi(splits[1])
			if err != nil {
				fmt.Printf("Could not parse %v", splits[0])
				panic(err)
			}
			rules[n2] = append(rules[n2], n1) // potentially long arrays could be sped up with sorting or other structures...
		}
		if updateExp.MatchString(line) {
			splits := strings.Split(line, ",")
			update := make([]int, 0, len(splits))
			for s := range splits {
				n, err := strconv.Atoi(splits[s])
				if err != nil {
					fmt.Printf("Could not parse %v", splits[s])
					panic(err)
				}
				update = append(update, n)
			}
			updates = append(updates, update)
		}
	}
	return rules, updates
}

func main() {
	fmt.Printf("Hello Day5\n")
	rules, updates := parseInput("./input.txt")
	validUpdates := make([][]int, 0, len(updates))
updateLoop:
	for _, update := range updates {
		for i := 0; i < len(update); i++ {
			number := update[i]
			for j := i + 1; j < len(update); j++ {
				if slices.Index(rules[number], update[j]) > -1 {
					continue updateLoop
				}
			}
		}
		validUpdates = append(validUpdates, update)
	}
	sum := 0
	for _, update := range validUpdates {
		middle := int(math.Floor(float64(len(update) / 2)))
		sum += update[middle]
	}
	fmt.Printf("Sum of middles: %v\n", sum)
}
