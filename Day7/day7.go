package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operator interface {
	Calc(a, b int) int
}

type Add struct{}

func (add Add) Calc(a, b int) int {
	return a + b
}

type Multiply struct{}

func (mul Multiply) Calc(a, b int) int {
	return a * b
}

type Concat struct{}

func (con Concat) Calc(a, b int) int {
	num, err := strconv.Atoi(fmt.Sprintf("%v%v", a, b))
	if err != nil {
		panic(err)
	}
	return num
}

var possibleOperators []Operator = []Operator{Add{}, Multiply{}, Concat{}}

type Equation struct {
	Numbers   []int
	Operators []Operator
	Result    int
}

func (e *Equation) Solve() (int, error) {
	if len(e.Numbers) < 1 {
		err := fmt.Errorf("numbers array is empty")
		return 0, err
	}
	result := e.Numbers[0]
	for i := 0; i < len(e.Operators); i++ {
		result = e.Operators[i].Calc(result, e.Numbers[i+1])
	}
	return result, nil
}

func (e *Equation) FindCombination(index int) bool {
	if len(e.Numbers) < 1 {
		return false
	}

	for i := 0; i < len(possibleOperators); i++ {
		e.Operators[index] = possibleOperators[i]
		if index < len(e.Operators)-1 {
			if e.FindCombination(index + 1) {
				return true
			}
		} else {
			s, err := e.Solve()
			if err != nil {
				panic(err)
			}
			if s == e.Result {
				return true
			}
		}
	}
	return false
}

func parseInput(p string) []Equation {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	equations := []Equation{}

	for scanner.Scan() {
		line := scanner.Text()
		e := Equation{}
		s := strings.Split(line, ":")
		e.Result, err = strconv.Atoi(s[0])
		if err != nil {
			panic(err)
		}
		numStrings := strings.TrimLeft(s[1], " ")
		s = strings.Split(numStrings, " ")
		for i := range s {
			num, err := strconv.Atoi(s[i])
			if err != nil {
				panic(err)
			}
			e.Numbers = append(e.Numbers, num)
		}
		equations = append(equations, e)
	}
	return equations
}

func main() {
	fmt.Printf("Hello Day7!\n")
	equations := parseInput("./input.txt")

	correctEquations := make(map[string]Equation)
	for _, e := range equations {
		//init operators
		e.Operators = make([]Operator, len(e.Numbers)-1)
		if e.FindCombination(0) {
			s := ""
			for _, v := range e.Numbers {
				s = s + "," + fmt.Sprint(v)
				s = strings.TrimLeft(s, ",")
			}
			correctEquations[s] = e
		}
	}

	sum := 0
	for k := range correctEquations {
		sum += correctEquations[k].Result
	}

	fmt.Printf("Result %v\n", sum)
}
