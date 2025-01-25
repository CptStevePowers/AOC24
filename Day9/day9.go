package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Disk []int64

func parseInput(p string) Disk {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	// line, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}

	disk := make(Disk, 0)
	var fileId int64 = 0
	isFile := true
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return disk
			}
			panic(err)
		}
		l := string(b)
		v, err := strconv.Atoi(string(l))
		if err != nil {
			panic(err)
		}
		if isFile {
			for j := 0; j < v; j++ {
				disk = append(disk, fileId)
			}
			fileId++
		} else {
			for j := 0; j < v; j++ {
				disk = append(disk, -1)
			}
		}
		isFile = !isFile
	}
	return disk
}

func (disk Disk) FormatDisk() Disk {
	leftPointer := 0
	rightPointer := len(disk) - 1

	for leftPointer < rightPointer {
		if disk[leftPointer] < 0 && disk[rightPointer] > -1 {
			disk[leftPointer], disk[rightPointer] = disk[rightPointer], disk[leftPointer]
		}

		for disk[leftPointer] > -1 {
			leftPointer++
		}
		for disk[rightPointer] < 0 {
			rightPointer--
		}
	}
	return disk
}

func (disk Disk) String() string {
	s := ""
	for _, v := range disk {
		if v < 0 {
			s += "."
		} else {
			s += fmt.Sprintf("%v", v)
		}
		s += " "
	}
	return s
}

func (disk Disk) CheckSum() int64 {
	var sum int64 = 0
	var pos int64 = 0
	for _, v := range disk {
		if v > 0 {
			r := pos * v
			sum += r
		}
		pos++
	}
	return sum
}

func main() {
	fmt.Print("Hello Day9\n")
	unformattedDisk := parseInput("./input.txt")
	fmt.Printf("%v\n", unformattedDisk)
	formattedDisk := unformattedDisk.FormatDisk()
	fmt.Printf("\n%v\n", formattedDisk)
	fmt.Printf("Checksum: %v\n", formattedDisk.CheckSum())
}
