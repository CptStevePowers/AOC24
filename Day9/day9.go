package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
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

func (disk Disk) FormatDiskP1() Disk {
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

type File struct {
	Id, Start, Size int64
}

type FileSystem struct {
	Files []File
}

func (fs FileSystem) String() string {
	s := ""
	for _, file := range fs.Files {
		s += fmt.Sprintf("%v(%v, %v) ", file.Id, file.Start, file.Size)
	}
	return s
}

func (fs *FileSystem) GetFile(id int64) *File {
	if i := slices.IndexFunc(fs.Files, func(f File) bool { return f.Id == id }); i > -1 {
		return &fs.Files[i]
	} else {
		return nil
	}
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
func NewFileSystem(disk Disk) FileSystem {
	var pos int64 = 0
	j := pos + 1
	fs := FileSystem{Files: make([]File, 0)}
	for j < int64(len(disk)) && pos < int64(len(disk)) {
		j = pos + 1
		fileId := disk[pos]
		f := File{Id: fileId, Start: pos, Size: 1}
		for j < int64(len(disk)) && disk[j] == fileId {
			f.Size++
			j++
		}

		if len(fs.Files) > 0 && fs.Files[len(fs.Files)-1].Id == f.Id {
			fs.Files[len(fs.Files)-1].Size += f.Size
		} else {
			fs.AddFile(f)
		}
		pos = pos + f.Size
	}
	return fs
}

func (disk Disk) Part2() Disk {
	fs := NewFileSystem(disk)
	fileId := slices.MaxFunc(fs.Files, func(a, b File) int { return int(a.Id) - int(b.Id) }).Id
	for fileId > -1 {
		filePos := slices.IndexFunc(fs.Files, func(f File) bool { return f.Id == fileId })
		if filePos < 0 {
			panic(fmt.Errorf("i should not be smaller 0"))
		}
		fileSize := fs.Files[filePos].Size
		gapPos := slices.IndexFunc(fs.Files, func(f File) bool {
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
	return fs.ToDisk()
}

func (fs FileSystem) ToDisk() Disk {
	disk := make(Disk, 0, len(fs.Files))
	for i := 0; i < len(fs.Files); i++ {
		for j := 0; j < int(fs.Files[i].Size); j++ {
			disk = append(disk, fs.Files[i].Id)
		}
	}
	return disk
}

func main() {
	fmt.Print("Hello Day9\n")
	unformattedDisk := parseInput("./input-benchmark.txt")
	fmt.Printf("%v\n", unformattedDisk)
	formattedDisk := unformattedDisk.Part2()
	fmt.Printf("\n%v\n", formattedDisk)
	fmt.Printf("Part2: %v\n", formattedDisk.CheckSum())
}
