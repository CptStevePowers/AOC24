package main

import (
	"bufio"
	"fmt"
	"os"
)

type Field struct {
	value     byte
	X, Y      uint
	gridIndex uint
	Grid      *Grid
}

type Grid struct {
	width, height uint
	fields        []Field
}

func (*Grid) New(width, height uint) Grid {
	fields := make([]Field, 0, width*height)
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		field.gridIndex = uint(i)
		fields[i] = field
	}
	return Grid{width: width, height: height, fields: fields}
}

func (g *Grid) AddRow(row []Field) *Grid {
	if g.width == 0 {
		g.width = uint(len(row))
	} else if g.width != uint(len(row)) {
		err := fmt.Errorf("new row is not same size as grid width (%v != %v)", len(row), g.width)
		panic(err)
	}
	var newIndex uint = 0
	var y uint = 0
	if len(g.fields) > 0 {
		lastIndex := g.fields[len(g.fields)-1].gridIndex
		y = g.fields[lastIndex-1].Y + 1
		newIndex = lastIndex + 1

	}
	for i := 0; i < len(row); i++ {
		row[i].gridIndex = newIndex
		row[i].Y = y
		row[i].X = uint(i)
		row[i].Grid = g
		newIndex++
	}
	newFields := append(g.fields, row...)
	g.fields = newFields
	g.height++
	return g
}

func (g *Grid) GetFieldByCoord(x, y uint) *Field {
	i := y*g.width + x
	return &g.fields[i]
}

func (f *Field) Walk(x, y int) (*Field, error) {
	var err error = nil
	xNew := f.X
	if int(f.X)+x < 0 {
		xNew = 0
		err = fmt.Errorf("Out of bounds")
	} else if int(f.X)+x >= int(f.Grid.width) {
		xNew = f.Grid.width - 1
		err = fmt.Errorf("Out of bounds")
	} else {
		xNew = f.X + uint(x)
	}

	yNew := f.Y
	if int(f.Y)+y < 0 {
		yNew = 0
		err = fmt.Errorf("Out of bounds")
	} else if int(f.Y)+y >= int(f.Grid.height) {
		yNew = f.Grid.height - 1
		err = fmt.Errorf("Out of bounds")
	} else {
		yNew = f.Y + uint(y)
	}

	return f.Grid.GetFieldByCoord(xNew, yNew), err
}

func parseInput(p string) Grid {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()
	width := 0
	scanner := bufio.NewScanner(f)
	grid := Grid{}
	for scanner.Scan() {
		line := scanner.Text()
		if width == 0 {
			width = len(line)
		} else if width != len(line) {
			err := fmt.Errorf("different line lengths %v is not equal to previous (%v)", len(line), width)
			panic(err)
		}
		row := make([]Field, 0, width)
		for i := 0; i < len(line); i++ {
			field := Field{value: line[i]}
			row = append(row, field)
		}
		grid.AddRow(row)
	}
	return grid
}

func (g *Grid) findAllWords(word string) [][]Field {
	results := make([][]Field, 0, 8)
	for i := 0; i < len(g.fields); i++ {
		f := g.fields[i]
		if f.value != word[0] {
			continue
		} else {
			directions := [][]int{
				{-1, 1},  //down left
				{0, 1},   //down
				{1, 1},   //down right
				{-1, 0},  // left
				{1, 0},   // right
				{-1, -1}, //up left
				{0, -1},  //up
				{1, -1},  //up right
			}
			for j := 0; j < len(directions); j++ {
				pos := &f
				for k := 1; k < len(word); k++ {
					newField, err := pos.Walk(directions[j][0], directions[j][1])
					// fmt.Printf("Checking field %v %v: %v for %v\n", newField.X, newField.Y, string(newField.value), string(word[k]))
					if err != nil {
						break
					}
					if word[k] != newField.value {
						break
					}
					if k == len(word)-1 {
						// fmt.Printf("%s found at %v %v\n", word, f, newField)
						results = append(results, []Field{f, *newField})
					}
					pos = newField
				}
			}
		}
	}
	return results
}

func (g *Grid) findXMAS() []Field {
	results := make([]Field, 0, 8)
	for i := 0; i < len(g.fields); i++ {
		f := g.fields[i]
		if f.value != 'A' {
			continue
		} else {
			directions := [][]int{
				{-1, 1},  //down left
				{1, 1},   //down right
				{-1, -1}, //up left
				{1, -1},  //up right
			}
			downLeft, err := f.Walk(directions[0][0], directions[0][1])
			if err != nil {
				continue
			}
			upRight, err := f.Walk(directions[3][0], directions[3][1])
			if err != nil {
				continue
			}
			if !(downLeft.value == 'M' && upRight.value == 'S') && !(downLeft.value == 'S' && upRight.value == 'M') {
				continue
			}

			downRight, err := f.Walk(directions[1][0], directions[1][1])
			if err != nil {
				continue
			}
			upLeft, err := f.Walk(directions[2][0], directions[2][1])
			if err != nil {
				continue
			}
			if !(downRight.value == 'M' && upLeft.value == 'S') && !(downRight.value == 'S' && upLeft.value == 'M') {
				continue
			}
			results = append(results, f)
		}
	}
	return results
}

func main() {
	fmt.Printf("Hello Day4\n")
	grid := parseInput("./input.txt")
	// fmt.Printf("%v", grid)
	word := "XMAS"
	fmt.Printf("\"%s\" found %v times\n", word, len(grid.findAllWords("XMAS")))
	fmt.Printf("\"%s\" found \"X-MAS\" %v times\n", word, len(grid.findXMAS()))

}
