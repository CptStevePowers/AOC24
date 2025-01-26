package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Coordinate [2]int

const maxIncline int = 1

var validDirections []Coordinate = []Coordinate{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

// a trail...
// ...starts at 0 and ends at 9
// ...steadily increases by 1
// ...can not go diagonal
type TrailMap struct {
	Width, Height int
	Grid          map[Coordinate]int
}

func (tm *TrailMap) IsInBounds(coord Coordinate) bool {
	return coord[0] >= 0 && coord[1] >= 0 && coord[0] < tm.Width && coord[1] < tm.Height
}

type Trail struct {
	TrailHead Coordinate
	TrailEnd  Coordinate
	Length    int
}

type Stepper struct {
	History []Coordinate
}

func (stepper *Stepper) Position() Coordinate {
	if len(stepper.History) < 1 {
		panic(fmt.Errorf("stepper has no history"))
	}
	return stepper.History[len(stepper.History)-1]
}

func (stepper *Stepper) StepCount() int {
	return len(stepper.History)
}

func (stepper *Stepper) Walk(tm TrailMap, direction Coordinate) (newCoord Coordinate, err error) {
	if slices.Index(validDirections, direction) < 0 {
		return Coordinate{}, fmt.Errorf("direction %v must be in %v", direction, validDirections)
	}
	stepperPos := stepper.Position()
	newCoord = Coordinate{stepperPos[0] + direction[0], stepperPos[1] + direction[1]}
	if !tm.IsInBounds(newCoord) {
		return Coordinate{}, fmt.Errorf("out of bounds. %v must be within grid w: %v, h: %v", newCoord, tm.Width, tm.Height)
	}

	if tm.Grid[stepperPos]+maxIncline != tm.Grid[newCoord] {
		return Coordinate{}, fmt.Errorf("terrain difference exceeded max incline %v=%v+%v", tm.Grid[newCoord], tm.Grid[stepperPos], maxIncline)
	}
	stepper.History = append(stepper.History, newCoord)
	return newCoord, nil
}

func parseInput(p string) TrailMap {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	tm := TrailMap{Width: 0, Height: 0, Grid: make(map[Coordinate]int)}
	for scanner.Scan() {
		line := scanner.Text()
		if tm.Width == 0 {
			tm.Width = len(line)
		} else if tm.Width != len(line) {
			err := fmt.Errorf("line lengths don't match")
			panic(err)
		}
		for i := 0; i < len(line); i++ {
			n, err := strconv.Atoi(string(line[i]))
			if err != nil {
				panic(err)
			}
			coord := Coordinate{i, tm.Height}
			tm.Grid[coord] = n
		}
		tm.Height++
	}
	return tm
}

func ScoreTrail(trailMap TrailMap, trailHead Coordinate) int {
	if trailMap.Grid[trailHead] != 0 {
		return -1
	}
	steppers := []Stepper{{History: []Coordinate{trailHead}}}
	completedTrails := make([]Stepper, 0)
	for len(steppers) > 0 {
		current := steppers[0]
		steppers = steppers[1:]
		for _, d := range validDirections {
			newStepper := Stepper{History: make([]Coordinate, len(current.History))}
			copy(newStepper.History, current.History)
			if newPos, err := newStepper.Walk(trailMap, d); err != nil {
				continue
			} else if trailMap.Grid[newPos] == 9 {
				completedTrails = append(completedTrails, newStepper)
			} else {
				steppers = append(steppers, newStepper)
			}
		}
	}

	// finishedSteppersByDestination := make(map[Coordinate][]Stepper)
	// for _, stepper := range completedTrails {
	// finishedSteppersByDestination[stepper.Position()] = append(finishedSteppersByDestination[stepper.Position()], stepper)
	// }

	score := len(completedTrails)
	return score
}

func main() {
	fmt.Print("Hello Day10\n")
	tm := parseInput("./input.txt")
	trailHeads := make([]Coordinate, 0)
	for key := range tm.Grid {
		if tm.Grid[key] == 0 {
			trailHeads = append(trailHeads, key)
		}
	}
	total := 0
	for _, th := range trailHeads {
		total += ScoreTrail(tm, th)
	}
	fmt.Printf("%v", total)
}
