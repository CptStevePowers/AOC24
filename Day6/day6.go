package main

import (
	"bufio"
	"fmt"
	"os"
)

type Guard struct {
	direction rune
	X, Y      uint
	Grid      *Grid
}

type Obstruction struct {
	X, Y uint
}

type Grid struct {
	width, height uint
	obstructions  map[[2]uint]Obstruction
	guard         Guard
}

func (g *Grid) IsObstructed(coord [2]uint) bool {
	x := coord[0]
	y := coord[1]
	if x > g.height-1 || y > g.width-1 {
		err := fmt.Errorf("out of bounds h: %v, w: %v, (%v)", g.height-1, g.width-1, coord)
		panic(err)
	}

	if g.obstructions[coord] == (Obstruction{}) {
		return false
	} else {
		return true
	}
}

func directionToVector(direction rune) (vector [2]int) {
	switch direction {
	case '^':
		vector = [2]int{0, -1}
	case 'v':
		vector = [2]int{0, 1}
	case '<':
		vector = [2]int{-1, 0}
	case '>':
		vector = [2]int{1, 0}
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

func (grid *Grid) NextLocationByVector(position [2]uint, vector [2]int) ([2]uint, error) {
	xNew, yNew := vector[0]+int(position[0]), vector[1]+int(position[1])
	if xNew < 0 || yNew < 0 || xNew > int(grid.width)-1 || yNew > int(grid.height)-1 {
		err := fmt.Errorf("out of bounds %v", [2]int{xNew, yNew})
		return position, err
	}
	return [2]uint{uint(xNew), uint(yNew)}, nil
}

func (g *Guard) Patrol() (newCoord [2]uint, err error) {
	vector := directionToVector(g.direction)

	newCoord, err = g.Grid.NextLocationByVector([2]uint{g.X, g.Y}, vector)
	if err != nil {
		return newCoord, err
	}

	if g.Grid.IsObstructed(newCoord) {
		//turn
		d := nextDirection(g.direction)
		g.direction = d
	} else {
		//step
		g.X = newCoord[0]
		g.Y = newCoord[1]
	}

	return [2]uint{g.X, g.Y}, err
}

func parseInput(p string) *Grid {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	o := make(map[[2]uint]Obstruction)
	grid := Grid{obstructions: o}

	isGuardSet := false
	for scanner.Scan() {
		line := scanner.Text()
		if grid.width == 0 {
			grid.width = uint(len(line))
		} else if grid.width != uint(len(line)) {
			err := fmt.Errorf("different line lengths %v is not equal to previous (%v)", len(line), grid.width)
			panic(err)
		}
		y := grid.height
		for x, v := range line {
			switch {
			case v == '#':
				obstruction := Obstruction{X: uint(x), Y: y}
				grid.obstructions[[2]uint{obstruction.X, obstruction.Y}] = obstruction
			case v == '^' || v == '<' || v == 'v' || v == '>':
				if isGuardSet {
					err := fmt.Errorf("Guard already set at (%v,%v)", grid.guard.X, grid.guard.Y)
					panic(err)
				}
				guard := Guard{X: uint(x), Y: y, Grid: &grid, direction: v}
				grid.guard = guard
				isGuardSet = true
			}
		}
		grid.height++
	}
	return &grid
}

func predictPatrol(grid *Grid) int {
	hist := make(map[[2]uint]rune)
	guardStartingCoord := [2]uint{grid.guard.X, grid.guard.Y}
	hist[guardStartingCoord] = grid.guard.direction

	for {
		coord, err := grid.guard.Patrol()
		if err != nil {
			fmt.Printf("Stopping patrol - %s\n", err)
			break
		}
		hist[coord] = grid.guard.direction
	}
	histKeys := make([][2]uint, 0, len(hist))
	for k := range hist {
		histKeys = append(histKeys, k)
	}

	return len(histKeys)
}

func findLoops(grid *Grid) int {
	hist := make(map[[2]uint]rune)
	guardStartingCoord := [2]uint{grid.guard.X, grid.guard.Y}
	hist[guardStartingCoord] = grid.guard.direction
	loops := make(map[[2]uint]rune)
	for {
		startingPosition, startingDirection := [2]uint{grid.guard.X, grid.guard.Y}, grid.guard.direction
		obstaclePosition, err := grid.NextLocationByVector(startingPosition, directionToVector(startingDirection))
		if err == nil && grid.checkForLoop(obstaclePosition, hist) {
			loops[obstaclePosition] = grid.guard.direction
		}

		coord, err := grid.guard.Patrol()
		if err != nil {
			break
		}
		hist[coord] = grid.guard.direction
	}
	return len(loops)
}

func (grid *Grid) checkForLoop(obstaclePosition [2]uint, hist map[[2]uint]rune) bool {
	startingPosition, startingDirection := [2]uint{grid.guard.X, grid.guard.Y}, grid.guard.direction
	if (grid.obstructions[obstaclePosition] != Obstruction{}) {
		// position already occupied
		return false
	}

	grid.obstructions[obstaclePosition] = Obstruction{X: obstaclePosition[0], Y: obstaclePosition[1]}
	defer delete(grid.obstructions, obstaclePosition)

	virtualGuard := Guard{
		Grid:      grid,
		X:         startingPosition[0],
		Y:         startingPosition[1],
		direction: startingDirection,
	}

	virtualGuardHistory := hist
	for {
		newPos, err := virtualGuard.Patrol()
		if err != nil {
			return false // expected behavior for out of bounds
		}

		if newPos == startingPosition && startingDirection == virtualGuard.direction {
			// fmt.Printf("Found loop for obstacle placement %v\n", obstaclePosition)
			return true
		}
		virtualGuardHistory[newPos] = virtualGuard.direction
	}
}

type LoopSet struct {
	Existing []*Obstruction
}

func main() {
	fmt.Printf("Hello Day6!\n")
	grid := parseInput("./input.txt")

	cpGrid := *grid
	fmt.Printf("Visited %v unique positions\n", predictPatrol(&cpGrid))

	cpGrid = *grid
	fmt.Printf("Found %v potential positions for obstacles to cause loop\n", findLoops(&cpGrid))
}
