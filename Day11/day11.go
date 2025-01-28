package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

// part2 credit to https://github.com/AllanTaylor314/AdventOfCode/blob/main/2024/11.py#L15
func main() {
	fmt.Print("Hello Day11\n")
	input := parseInput("./input.txt")
	stones := make(map[int64]int64)

	for _, n := range input {
		stones[n]++
	}

	start := time.Now()
	for i := 0; i < 75; i++ {
		newStones := make(map[int64]int64)
		for number := range stones {
			if number == 0 {
				newStones[1] += stones[number]
			} else if s := fmt.Sprintf("%v", number); len(s)%2 == 0 {
				s1, s2 := s[:int(len(s)/2)], s[int(len(s)/2):]
				n1, err := strconv.Atoi(s1)
				if err != nil {
					panic(err)
				}
				n2, err := strconv.Atoi(s2)
				if err != nil {
					panic(err)
				}
				newStones[int64(n1)] += stones[number]
				newStones[int64(n2)] += stones[number]
			} else {
				newStones[number*2024] += stones[number]
			}
		}
		stones = newStones
	}
	var count int64 = 0
	for key := range stones {
		count += stones[key]
	}
	fmt.Printf("Took %s seconds\n", time.Since(start))
	fmt.Printf("Number of stones: %v\n", count)
}
