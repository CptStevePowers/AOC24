package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseInput(p string) ([]int, []int) {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error opening the file %s", p)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	l1 := make([]int, 0, 32)
	l2 := make([]int, 0, 32)

	for s.Scan() {
		l := s.Text()
		splits := strings.Split(l, "   ")
		v1, err := strconv.Atoi(splits[0])
		if err != nil {
			fmt.Printf("Error parsing value of splits %s: %s", splits, splits[0])
		}
		l1 = append(l1, v1)

		v2, err := strconv.Atoi(splits[1])
		if err != nil {
			fmt.Printf("Error parsing value of splits %s: %s", splits, splits[1])
		}
		l2 = append(l2, v2)
	}
	return l1, l2
}

func calculateTotalDifference(l1, l2 []int) int {
	sort.Slice(l1, func(i, j int) bool {
		return l1[i] < l1[j]
	})
	sort.Slice(l2, func(i, j int) bool {
		return l2[i] < l2[j]
	})
	totalDiff := 0
	for i := 0; i < len(l1); i++ {
		diff := max(l1[i], l2[i]) - min(l1[i], l2[i])
		totalDiff += diff
	}
	return totalDiff
}

func calculateSimilarity(l1, l2 []int) int {
	maxValue := 0
	for i := 0; i < len(l2); i++ {
		if l2[i] > maxValue {
			maxValue = l2[i]
		}
	}

	similarities := make([]int, maxValue+1)
	for i := 0; i < len(l2); i++ {
		similarities[l2[i]] += 1
	}

	score := 0
	for i := 0; i < len(l1); i++ {
		score += similarities[l1[i]] * l1[i]
	}
	return score
}

func main() {
	l1, l2 := parseInput("./input.txt")
	totalDiff := calculateTotalDifference(l1, l2)
	fmt.Printf("totalDifference: %v\n", totalDiff)
	score := calculateSimilarity(l1, l2)
	fmt.Printf("Similarity score: %v\n", score)
}
