package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput(p string) []int64 {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	numbers := make([]int64, 0)
	for scanner.Scan() {
		line := scanner.Text()
		blocks := strings.Split(line, " ")
		for i := 0; i < len(blocks); i++ {
			num, err := strconv.Atoi(blocks[i])
			if err != nil {
				panic(err)
			}
			numbers = append(numbers, int64(num))
		}
	}
	return numbers
}

//stone rules...
//... stones with 0 are replaced with 1
//... if even number of digits replace with two stones: left stone left half, right stone right half of digits (remove leading zeroes)
//... otherwise multiply stone by 2024

func Blink(stones []int64) []int64 {
	out := make([]int64, 0, len(stones))
	for _, stone := range stones {
		if stone == 0 {
			out = append(out, 1)
		} else if s := fmt.Sprintf("%v", stone); len(s)%2 == 0 {
			s1, err := strconv.Atoi(s[:len(s)/2])
			if err != nil {
				panic(err)
			}
			out = append(out, int64(s1))
			s2, err := strconv.Atoi(s[len(s)/2:])
			if err != nil {
				panic(err)
			}
			out = append(out, int64(s2))
		} else {
			out = append(out, stone*2024)
		}
	}
	return out
}

func main() {
	fmt.Print("Hello Day11\n")
	stones := parseInput("./input.txt")
	// fmt.Printf("%v\n", stones)
	for i := 0; i < 75; i++ {
		stones = Blink(stones)
		fmt.Printf("%v, %v\n", i, len(stones))
		// fmt.Printf("After %v blinks:\n%v\n", i+1, stones)
	}
	fmt.Printf("Number of stones: %v", len(stones))
}
