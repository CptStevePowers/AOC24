package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

func parseInput(p string) FileSystem {
	f, err := os.Open(p)
	if err != nil {
		fmt.Printf("Error reading file contents of %s\n", p)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	if err != nil {
		panic(err)
	}

	fs := FileSystem{Files: make([]File, 0)}
	var fileId int64 = 0
	isFile := true
	var fsPointer int64 = 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return fs
			}
			panic(err)
		}
		l := string(b)
		v, err := strconv.Atoi(string(l))
		if err != nil {
			panic(err)
		}
		if isFile {
			fs.AddFile(File{Id: fileId, Size: int64(v), Start: fsPointer})
			fileId++
		} else {
			fs.AddFile(File{Id: -1, Size: int64(v), Start: fsPointer})
		}
		isFile = !isFile
		fsPointer += int64(v)
	}
	return fs
}

func (fs FileSystem) String() string {
	s := ""
	for _, v := range fs.Files {
		for i := int64(0); i < v.Size; i++ {
			if v.Id < 0 {
				s += "."
			} else {
				s += fmt.Sprintf("%v", v.Id)
			}
		}
	}
	return s
}

func (fs *FileSystem) CheckSum() int64 {
	var sum int64 = 0
	var pos int64 = 0
	for _, v := range fs.Files {
		for i := int64(0); i < v.Size; i++ {
			if v.Id > 0 {
				r := pos * v.Id
				sum += r
			}
			pos++
		}
	}
	return sum
}

type File struct {
	Id, Start, Size int64
}

type FileSystem struct {
	Files []File
}

func (fs *FileSystem) AddFile(f File) *File {
	if fs.Files == nil {
		fs.Files = make([]File, 1)
	}
	fs.Files = append(fs.Files, f)
	return &fs.Files[len(fs.Files)-1]
}

// alrighty this is gonna go bad if we dont establish the correct datamodel :D
// something like Block with start, id, size
func (fs *FileSystem) Part2() *FileSystem {
	fileId := slices.MaxFunc(fs.Files, func(a, b File) int { return int(a.Id) - int(b.Id) }).Id
	for fileId > -1 {
		filePos := slices.IndexFunc(fs.Files, func(f File) bool { return f.Id == fileId })
		if filePos < 0 {
			panic(fmt.Errorf("i should not be smaller 0"))
		}
		fileSize := fs.Files[filePos].Size
		gapPos := slices.IndexFunc(fs.Files[:filePos], func(f File) bool {
			return f.Id < 0 && f.Size >= fileSize
		})

		if gapPos > -1 && gapPos < filePos {
			newGap := File{Start: fs.Files[filePos].Start, Size: fs.Files[filePos].Size, Id: -1}
			fs.AddFile(newGap)
			fs.Files[filePos].Start = fs.Files[gapPos].Start
			if leftOverBlocks := fs.Files[gapPos].Size - fs.Files[filePos].Size; leftOverBlocks > 0 {
				fs.Files[gapPos].Start = fs.Files[gapPos].Start + fs.Files[filePos].Size
				fs.Files[gapPos].Size = leftOverBlocks
			} else {
				fs.Files = append(fs.Files[:gapPos], fs.Files[gapPos+1:]...)
			}
		}
		fileId--
	}
	slices.SortFunc(fs.Files, func(a, b File) int { return int(a.Start - b.Start) })
	return fs
}

func main() {
	fmt.Print("Hello Day9\n")
	unformattedDisk := parseInput("./input.txt")
	// fmt.Printf("%v\n", unformattedDisk)
	formattedDisk := unformattedDisk.Part2()
	// fmt.Printf("\n%v\n", formattedDisk)
	fmt.Printf("Part2: %v\n", formattedDisk.CheckSum())
}
