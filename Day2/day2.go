package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput(p string) [][]int {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error opening the file %s", p)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	reports := make([][]int, 0, 32)

	for s.Scan() {
		l := s.Text()
		splits := strings.Split(l, " ")
		values := make([]int, 0, len(splits))
		for i := 0; i < len(splits); i++ {
			v, err := strconv.Atoi(splits[i])
			if err != nil {
				fmt.Printf("Fatal while trying to convert %v\n", splits[i])
				break
			}
			values = append(values, v)
		}
		reports = append(reports, values)
	}
	return reports
}

func filterSafeReports(reports [][]int) [][]int {
	output := make([][]int, 0, len(reports))
	for i := 0; i < len(reports); i++ {
		r := reports[i]
		if isReportSafe(r) {
			output = append(output, r)
		}
	}
	return output
}

func isReportSafe(r []int) bool {
	if len(r) < 2 {
		return true
	}

	isIncreasing := false
	if r[1]-r[0] > 0 {
		isIncreasing = true
	}

	for i := 0; i < len(r)-1; i++ {
		diff := r[i+1] - r[i]
		if !(diff > 0 && diff < 4 && isIncreasing) && !(diff < 0 && diff > -4 && !isIncreasing) {
			return false
		}
	}
	return true
}

func filterSafeReportsDampened(reports [][]int) [][]int {
	output := make([][]int, 0, len(reports))
	for i := 0; i < len(reports); i++ {
		r := reports[i]
		if isReportSafe(r) {
			output = append(output, r)
			continue
		}
		for j := 0; j < len(r); j++ {
			newReport := make([]int, len(r))
			copy(newReport, r)
			newReport = append(newReport[:j], r[j+1:]...)
			if isReportSafe(newReport) {
				output = append(output, r)
				break
			}
		}
	}
	return output
}

func main() {
	fmt.Printf("Hi Day2\n")
	reports := parseInput("./input.txt")
	safeReports := filterSafeReports(reports)
	fmt.Printf("There are %v safe reports\n", len(safeReports))
	okReports := filterSafeReportsDampened(reports)
	fmt.Printf("There are %v dampened safe reports\n", len(okReports))
}
