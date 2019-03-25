package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func filterFiles(filesInfo []os.FileInfo) (filteredFilesInfo []os.FileInfo) {
	for _, fi := range filesInfo {
		if fi.IsDir() {
			filteredFilesInfo = append(filteredFilesInfo, fi)
		}
	}
	return
}

func dirTreePrefix(output io.Writer, directory string, printFiles bool, prefix string) error {
	file, err := os.Open(directory)
	check(err)
	filesInfo, err := file.Readdir(0)
	check(err)

	if !printFiles {
		filesInfo = filterFiles(filesInfo)
	}

	sort.Slice(filesInfo, func(i, j int) bool {
		return filesInfo[i].Name() < filesInfo[j].Name()
	})

	level := strings.Count(directory, string(os.PathSeparator))

	for i, fi := range filesInfo {
		var sign string
		var sizeStr string
		isLast := i == len(filesInfo)-1
		if isLast {
			sign += "└───"
		} else {
			sign += "├───"
		}
		if !fi.IsDir() {
			sizeStr += " ("
			if size := fi.Size(); size > 0 {
				sizeStr += strconv.Itoa(int(size))
				sizeStr += "b"
			} else {
				sizeStr += "empty"
			}
			sizeStr += ")"
		}
		fmt.Fprintln(output, prefix+sign+fi.Name()+sizeStr)
		if fi.IsDir() {
			if isLast {
				prefix += "	"
			} else {
				prefix += "│	"
			}
			dirTreePrefix(output, directory+string(os.PathSeparator)+fi.Name(), printFiles, prefix)
			prefix = strings.Repeat("│	", level)
		}
	}

	return nil
}

func dirTree(output io.Writer, directory string, printFiles bool) error {
	return dirTreePrefix(output, directory, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
