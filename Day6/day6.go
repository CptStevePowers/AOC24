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
	Grid *Grid
}

type Grid struct {
	width, height uint
	obstructions  map[[2]uint]Obstruction
	guard         Guard
}

func (g *Grid) IsObstructed(coord [2]uint) (bool, error) {
	x := coord[0]
	y := coord[1]
	if x > g.height-1 || y > g.width-1 {
		err := fmt.Errorf("Out of bounds h: %v, w: %v, (%v)", g.height-1, g.width-1, coord)
		return false, err
	}

	if g.obstructions[coord] == (Obstruction{}) {
		return false, nil
	} else {
		return true, nil
	}
}

func (g *Guard) Walk() (newCoord [2]uint, err error) {
	xNew, yNew := g.X, g.Y
	switch g.direction {
	case '^':
		if yNew == 0 {
			err = fmt.Errorf("out of bounds (%v,-1)", g.X)
		} else {
			yNew--
		}
	case 'v':
		if yNew >= g.Grid.height-1 {
			err = fmt.Errorf("out of bounds (%v,%v)", g.X, g.Grid.height)
		} else {
			yNew++
		}
	case '<':
		if xNew == 0 {
			err = fmt.Errorf("out of bounds (-1,%v)", g.Y)
		} else {
			xNew--
		}
	case '>':
		if xNew >= g.Grid.width-1 {
			err = fmt.Errorf("out of bounds (%v, %v)", g.Grid.width, g.Y)
		} else {
			xNew++
		}
	}

	isBlocked, obstructionErr := g.Grid.IsObstructed([2]uint{xNew, yNew})
	if obstructionErr != nil {
		panic(obstructionErr)
	}

	if isBlocked {
		g.handleCollision()
	} else {
		g.X = xNew
		g.Y = yNew
	}

	return [2]uint{g.X, g.Y}, err
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
		return '?'
	}
}

func (g *Guard) handleCollision() {
	g.direction = nextDirection(g.direction)
}

func parseInput(p string) Grid {
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
				obstruction := Obstruction{X: uint(x), Y: y, Grid: &grid}
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
	return grid
}

func predictPatrol(grid Grid) int {
	hist := make(map[[2]uint]rune)
	guardStartingCoord := [2]uint{grid.guard.X, grid.guard.Y}
	hist[guardStartingCoord] = grid.guard.direction

	for {
		coord, err := grid.guard.Walk()
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

type LoopSet struct {
	Existing []*Obstruction
}

func findLoopSet(o *Obstruction, start *Obstruction, direction rune) (LoopSet, error) {
	switch direction {
	case '^':
		if o.X >= o.Grid.width-1 {
			err := fmt.Errorf("loop set impossible, out of bounds")
			return LoopSet{}, err
		}
		for x := o.X; x < o.Grid.width-2; x++ {
			coord := [2]uint{x, o.Y}
			hasObstacle, err := o.Grid.IsObstructed(coord)
			if err != nil {
				panic(err)
			}
			// find next obstacle
			if hasObstacle {
				if next := o.Grid.obstructions[coord]; next == *start {
					return LoopSet{Existing: []*Obstruction{o}}, nil
				} else {
					loop, err := findLoopSet(&next, start, nextDirection(direction))
					if err != nil {
						return LoopSet{}, err
					}
					loop.Existing = append(loop.Existing, o)
					return loop, nil
				}
			}
		}
	case '>':
		err := fmt.Errorf("not implemented")
		return LoopSet{}, err
	case '<':
		err := fmt.Errorf("not implemented")
		return LoopSet{}, err
	case 'v':
		err := fmt.Errorf("not implemented")
		return LoopSet{}, err
	default:
		err := fmt.Errorf("this should not happen")
		return LoopSet{}, err
	}
	err := fmt.Errorf("not implemented")
	return LoopSet{}, err
}

func main() {
	fmt.Printf("Hello Day6!\n")
	grid := parseInput("./example.txt")
	// fmt.Printf("Grid h: %v w: %v\n", grid.height, grid.width)
	// fmt.Printf("Guard at (%v,%v), facing: %c\n", grid.guard.X, grid.guard.Y, grid.guard.direction)
	// fmt.Printf("Obstructions at:")

	// for k := range grid.obstructions {
	// 	fmt.Printf(" %v", k)
	// }
	// fmt.Print("\n")
	fmt.Printf("Visited %v unique positions\n", predictPatrol(grid))
}
