package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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

func dirTree(output io.Writer, path string, printFiles bool) error {
	response := strings.Trim(buildTree("", "", path, printFiles), "\n")
	_, _ = fmt.Fprintln(output, response)
	return nil
}

func buildTree (result, prefix, path string, printFiles bool) string {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	var symbol string = `├───`
	//var description string = ""

	if !printFiles {
		files = filterOnlyDir(files)
	}

	lastIndex := len(files)

	for index, file := range files  {
		if index == lastIndex-1 {
			symbol = `└───`
		}

		//description = file.Name() + file.Size()
		description := fileDescription(file)
		result += prefix + symbol + description + "\n"

		if file.IsDir() {
			add := "│\t"
			if symbol == `└───` {
				add = "\t"
			}
			result = buildTree(result, prefix + add, path+"/"+file.Name(), printFiles)
		}
	}

	return result
}

func filterOnlyDir(files [] os.FileInfo) [] os.FileInfo {
	n := 0

	for _, file := range files {
		if file.IsDir() {
			files[n] = file
			n++
		}
	}

	return files[:n]
}

func fileDescription(file os.FileInfo) string {

	description := fmt.Sprintf("%v", file.Name())

	if !file.IsDir() {
		size := "(empty)"

		if file.Size() != 0 {
			size = fmt.Sprintf("(%vb)", file.Size())
		}

		description = fmt.Sprintf("%v %v", file.Name(), size)
	}

	return description
}
