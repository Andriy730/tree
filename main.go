package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	verticalLine             = "│"
	verticalLineWithContinue = "├"
	horizontalLine           = "───"
	verticalLineEnd          = "└"
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

func dirTree(out io.Writer, path string, printFiles bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path to " + path + " does not exist")
	}

	fmt.Println(path)
	printTree(out, path, printFiles, "")

	return nil
}

func printTree(out io.Writer, path string, printFiles bool, indent string) {
	files, _ := ioutil.ReadDir(path)
	var newIndent string
	newIndent = indent + verticalLine + "\t"
	if !printFiles {
		files = filterFiles(files)
	}
	for i, file := range files {
		start := indent
		if i == len(files)-1 {
			start += verticalLineEnd
			newIndent = indent + "\t"
		} else {
			start += verticalLineWithContinue
		}
		var size string
		if file.IsDir() {
			size = ""
		} else {
			if file.Size() > 0 {
				size = " (" + strconv.FormatInt(file.Size(), 10) + "b" + ")"
			} else {
				size = " (empty)"
			}
		}
		fmt.Fprint(out, start+horizontalLine+file.Name()+size+"\n")
		if file.IsDir() {
			printTree(out, path+"/"+file.Name(), printFiles, newIndent)
		}
	}
}

func filterFiles(files []os.FileInfo) []os.FileInfo {
	filteredFiles := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}
