package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Coordinate [2]int
type Frequency rune

type Antenna struct {
	Coordinate Coordinate
	Frequency  Frequency
}

type Grid struct {
	Width, Height int
	Antennas      map[Coordinate]Antenna
}

func parseInput(p string) *Grid {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	antennas := make(map[Coordinate]Antenna)
	grid := Grid{Antennas: antennas}

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
			case v == '.':
				continue
			case v != '.':
				frequency := Frequency(v)
				coord := Coordinate{x, y}
				location := [2]int{x, y}
				grid.Antennas[coord] = Antenna{Coordinate: location, Frequency: frequency}
			}
		}
		grid.Height++
	}
	return &grid
}

func FindAntinodes(a, b Antenna) []Coordinate {
	antinodes := make([]Coordinate, 0, 2)
	if a.Frequency != b.Frequency {
		return antinodes
	}
	vectorAB := [2]int{b.Coordinate[0] - a.Coordinate[0], b.Coordinate[1] - a.Coordinate[1]}
	antinodes = append(antinodes, Coordinate{a.Coordinate[0] - vectorAB[0], a.Coordinate[1] - vectorAB[1]})
	antinodes = append(antinodes, Coordinate{b.Coordinate[0] + vectorAB[0], b.Coordinate[1] + vectorAB[1]})
	return antinodes
}

func main() {
	fmt.Printf("Hello Day8!\n")
	grid := parseInput("./input.txt")
	antinodes := make(map[Coordinate][]Frequency)

	keys := make([]Coordinate, 0, len(grid.Antennas))
	for c := range grid.Antennas {
		keys = append(keys, c)
	}

	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			antennaA := grid.Antennas[keys[i]]
			antennaB := grid.Antennas[keys[j]]
			if antennaA.Frequency != antennaB.Frequency {
				continue
			}
			frequency := antennaA.Frequency
			a := FindAntinodes(antennaA, antennaB)
			for _, coord := range a {
				if coord[0] >= grid.Width || coord[0] < 0 || coord[1] >= grid.Height || coord[1] < 0 {
					continue
				}
				if len(antinodes[coord]) == 0 {
					antinodes[coord] = make([]Frequency, 0, 1)
				}
				if slices.Index(antinodes[coord], frequency) < 0 {
					antinodes[coord] = append(antinodes[coord], frequency)
				}
			}
		}
	}

	// print grid
	fmt.Print("Grid:\n")
	s := ""
	for y := 0; y < grid.Width; y++ {
		for x := 0; x < grid.Height; x++ {
			r := '.'
			if len(antinodes[Coordinate{x, y}]) > 0 {
				r = '#'
			}
			s += string(r)
		}
		s += "\n"
	}
	fmt.Print(s)

	fmt.Printf("Result: %v\n", len(antinodes))
}
