package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Coordinate [2]int
type History map[Coordinate][]rune

type Guard struct {
	Direction rune
	Position  Coordinate
	Grid      *Grid
	History   History
}

type Grid struct {
	Width, Height int
	Obstructions  map[Coordinate]bool
	Guard         Guard
}

func (g *Grid) IsObstructed(coord Coordinate) bool {
	if !g.IsInBounds(coord) {
		err := fmt.Errorf("out of bounds h: %v, w: %v, (%v)", g.Height-1, g.Width-1, coord)
		panic(err)
	}

	if g.Obstructions[coord] {
		return true
	} else {
		return false
	}
}

func directionToVector(direction rune) (vector Coordinate) {
	switch direction {
	case '^':
		vector = Coordinate{0, -1}
	case 'v':
		vector = Coordinate{0, 1}
	case '<':
		vector = Coordinate{-1, 0}
	case '>':
		vector = Coordinate{1, 0}
	default:
		err := fmt.Errorf("direction must not be '%c', valid directions: '^'|'>'|'v'|'<'", direction)
		panic(err)
	}
	return vector
}

func nextDirection(current rune) rune {
	switch current {
	case '^':
		return '>'
	case '>':
		return 'v'
	case 'v':
		return '<'
	case '<':
		return '^'
	default:
		err := fmt.Errorf("%c is invalid input", current)
		panic(err)
	}
}

func (grid *Grid) IsInBounds(position Coordinate) bool {
	x, y := position[0], position[1]
	if x >= 0 && y >= 0 && x < int(grid.Width) && y < int(grid.Height) {
		return true
	}
	return false
}

func (guard *Guard) Patrol() (newPosition Coordinate, err error) {
	vector := directionToVector(guard.Direction)

	xNew, yNew := vector[0]+guard.Position[0], vector[1]+guard.Position[1]
	newPosition = Coordinate{xNew, yNew}
	if !guard.Grid.IsInBounds(newPosition) {
		err = fmt.Errorf("out of bounds")
		return newPosition, err
	}

	if guard.Grid.IsObstructed(newPosition) {
		d := nextDirection(guard.Direction)
		guard.Direction = d
	} else {
		guard.Position = newPosition
	}

	return guard.Position, err
}

func parseInput(p string) *Grid {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	o := make(map[Coordinate]bool)
	grid := Grid{Obstructions: o}

	isGuardSet := false
	for scanner.Scan() {
		line := scanner.Text()
		if grid.Width == 0 {
			grid.Width = int(len(line))
		} else if grid.Width != int(len(line)) {
			err := fmt.Errorf("different line lengths %v is not equal to previous (%v)", len(line), grid.Width)
			panic(err)
		}
		y := grid.Height
		for x, v := range line {
			switch {
			case v == '#':
				grid.Obstructions[Coordinate{x, y}] = true
			case v == '^' || v == '<' || v == 'v' || v == '>':
				if isGuardSet {
					err := fmt.Errorf("Guard already set at (%v,%v)", grid.Guard.Position[0], grid.Guard.Position[1])
					panic(err)
				}
				guard := Guard{Position: Coordinate{x, y}, Grid: &grid, Direction: v}
				grid.Guard = guard
				isGuardSet = true
			}
		}
		grid.Height++
	}
	return &grid
}

func predictPatrol(guard Guard) History {
	guard.History = make(History)
	guardStartingPosition := guard.Position
	guard.History[guardStartingPosition] = append(guard.History[guardStartingPosition], guard.Direction)

	for {
		coord, err := guard.Patrol()
		if err != nil {
			fmt.Printf("Stopping patrol - %s\n", err)
			break
		}
		guard.History[coord] = append(guard.History[coord], guard.Direction)
	}

	return guard.History
}

func findLocationsForLoops(grid *Grid, history History) int {
	guardStartingPosition := grid.Guard.Position
	grid.Guard.History = make(History)
	grid.Guard.History[guardStartingPosition] = append(grid.Guard.History[guardStartingPosition], grid.Guard.Direction)

	loops := make(map[Coordinate]rune)
	for obstaclePosition := range history {
		if grid.checkForLoop(obstaclePosition) {
			loops[obstaclePosition] = grid.Guard.Direction
		}
	}

	return len(loops)
}

func (grid *Grid) checkForLoop(obstaclePosition Coordinate) bool {
	startingPosition, startingDirection := grid.Guard.Position, grid.Guard.Direction
	if grid.Obstructions[obstaclePosition] {
		return false
	}

	grid.Obstructions[obstaclePosition] = true
	defer delete(grid.Obstructions, obstaclePosition)

	virtualGuard := Guard{
		Grid:      grid,
		Position:  startingPosition,
		Direction: startingDirection,
		History:   make(History),
	}

	for {
		newPos, err := virtualGuard.Patrol()
		if err != nil {
			return false // expected behavior for out of bounds
		}

		if slices.Index(virtualGuard.History[newPos], virtualGuard.Direction) > -1 {
			return true
		}
		virtualGuard.History[newPos] = append(virtualGuard.History[newPos], virtualGuard.Direction)
	}
}

func main() {
	fmt.Printf("Hello Day6!\n")
	grid := parseInput("./input.txt")

	cpGrid := *grid
	uniquePositions := predictPatrol(cpGrid.Guard)
	fmt.Printf("Visited %v unique positions\n", len(uniquePositions))

	cpGrid = *grid
	fmt.Printf("Found %v potential positions for obstacles to cause loop\n", findLocationsForLoops(&cpGrid, uniquePositions))
}
