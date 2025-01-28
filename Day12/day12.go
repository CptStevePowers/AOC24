package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"time"
)

type Coordinate [2]int

type Plot struct {
	Coordinate Coordinate
	Plant      string
	Region     *Region
}

type Garden map[Coordinate]*Plot

func (g Garden) Width() int {
	w := 0
	for plot := range g {
		if plot[1] > w {
			w = plot[1]
		}
	}
	return w + 1
}

func (g Garden) Height() int {
	h := 0
	for plot := range g {
		if plot[0] > h {
			h = plot[0]
		}
	}
	return h + 1
}

func parseInput(p string) Garden {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	garden := make(Garden)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x := 0; x < len(line); x++ {
			coord := Coordinate{x, y}
			garden[coord] = &Plot{Coordinate: coord, Plant: string(line[x])}
		}
		y++
	}
	return garden
}

// regions are places with the same kind of plant
// know area and perimeter of region
// price = area * perimeter
type Region struct {
	Plots []*Plot
	Plant string
}

func MergeRegions(a, b Region) Region {
	if a.Plant != b.Plant {
		panic(fmt.Errorf("plants not equal! %v != %v", a.Plant, b.Plant))
	}
	a.Plots = append(a.Plots, b.Plots...)
	for _, plot := range a.Plots {
		plot.Region = &a
	}
	return a
}

type Vector [2]int

var validDirections []Vector = []Vector{{-1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func (r *Region) Area() int {
	return len(r.Plots)
}

func (r *Region) Perimeter() int {
	perimeter := 0
	regionMap := make(map[Coordinate]bool)
	for _, plot := range r.Plots {
		regionMap[plot.Coordinate] = true
	}
	for coord := range regionMap {
		for _, direction := range validDirections {
			nextCoord := Coordinate{coord[0] + direction[0], coord[1] + direction[1]}
			if !regionMap[nextCoord] {
				perimeter++
			}
		}
	}
	return perimeter
}

func Move(pos Coordinate, dir Vector) Coordinate {
	return Coordinate{pos[0] + dir[0], pos[1] + dir[1]}
}

type Fence [2]Coordinate

func NewFence(a, b Coordinate) Fence {

	if a[0] < b[0] {
		return Fence{a, b}
	} else if a[0] > b[0] {
		return Fence{b, a}
	} else if a[1] < b[1] {
		return Fence{a, b}
	} else if a[1] > b[1] {
		return Fence{b, a}
	}
	return Fence{a, b}
}

func (r *Region) Sides() int {
	regionMap := make(map[Coordinate]bool)
	for _, plot := range r.Plots {
		regionMap[plot.Coordinate] = true
	}

	allFences := make([]Fence, 0)
	for currentCoordinate := range regionMap {
		for _, direction := range validDirections {
			nextCoordinate := Move(currentCoordinate, direction)
			if !regionMap[nextCoordinate] {
				allFences = append(allFences, NewFence(currentCoordinate, nextCoordinate))
			}
		}
	}

	horizontalFences := make(map[int][]int)
	verticalFences := make(map[int][]int)
	for _, fence := range allFences {
		if fence[0][0] == fence[1][0] {
			verticalFences[fence[0][0]] = append(verticalFences[fence[0][0]], fence[0][1])
		} else if fence[0][1] == fence[1][1] {
			horizontalFences[fence[0][1]] = append(horizontalFences[fence[0][1]], fence[0][0])
		} else {
			panic(fmt.Errorf("this should really not happen help!"))
		}
	}

	//TODO: HELP ME
	// sides := 0
	// for row := range horizontalFences {
	// 	cols := horizontalFences[row]
	// 	cols = slices.SortedFunc(cols, func(a, b int) int { return a - b })

	// }

	return len(allFences)
}

type Queue []Coordinate

func Push(q Queue, coord Coordinate) Queue {
	return append(Queue{coord}, q...)
}

func GetRegions(garden Garden) map[*Region]bool {
	queue := make(Queue, len(garden))
	for key := range garden {
		queue = append(queue, key)
	}

	for len(queue) > 0 {
		var currentPlot *Plot = nil
		if currentPlot == nil {
			coord := queue[0]
			currentPlot = garden[coord]
		}

		if r := slices.Index(queue, currentPlot.Coordinate); r > -1 {
			queue = append(queue[:r], queue[r+1:]...)
		}

		if currentPlot.Region == nil {
			currentPlot.Region = &Region{Plant: currentPlot.Plant, Plots: []*Plot{currentPlot}}
		}

		//exists valid region in neighbors?
		for _, direction := range validDirections {
			neighborCoord := Move(currentPlot.Coordinate, direction)
			neighborPlot := garden[neighborCoord]
			if neighborPlot == nil {
				continue
			}
			if neighborPlot.Plant != currentPlot.Plant {
				continue
			}
			if neighborPlot.Region == nil {
				neighborPlot.Region = currentPlot.Region
				currentPlot.Region.Plots = append(currentPlot.Region.Plots, neighborPlot)
				queue = Push(queue, neighborPlot.Coordinate)
				continue
			}
			if neighborPlot.Region != currentPlot.Region {
				MergeRegions(*neighborPlot.Region, *currentPlot.Region)
			}
		}
	}

	regions := make(map[*Region]bool)
	for coord := range garden {
		regions[garden[coord].Region] = true
	}
	return regions
}

func Part1(regions map[*Region]bool) int {
	total := 0
	for r := range regions {
		area := r.Area()
		perimeter := r.Perimeter()
		price := area * perimeter
		// fmt.Printf("A region of %s plants with price %v * %v = %v\n", r.Plant, area, perimeter, price)
		total += price
	}
	return total
}

func Part2(regions map[*Region]bool) int {
	total := 0
	for r := range regions {
		area := r.Area()
		perimeter := r.Sides()
		price := area * perimeter
		// fmt.Printf("A region of %s plants with price %v * %v = %v\n", r.Plant, area, perimeter, price)
		total += price
	}
	return total

}

func main() {
	fmt.Print("Hello Day12\n")
	garden := parseInput("./input.txt")
	start := time.Now()
	regions := GetRegions(garden)
	fmt.Printf("Get Regions done after %v\n", time.Since(start))
	start = time.Now()
	total := Part1(regions)
	fmt.Printf("Part1 done after %v\n", time.Since(start))
	fmt.Printf("Total price: %v\n", total)
	start = time.Now()
	total = Part2(regions)
	fmt.Printf("Part2 done after %v\n", time.Since(start))
	fmt.Printf("Total price: %v\n", total)
}
