package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	// "debug/gosym"
	"strconv"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	const endGraphic = "└───";
	const preGraphic = "|";
	const intermediateGraphic = "├───";

	var paths []string;
	var prevNodeStair string;
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, errOpen := os.Open(path);
			if errOpen != nil {
				panic(err.Error());
			}
			defer file.Close();
			stat, _ := file.Stat();

			fileSize := stat.Size();
			if fileSize == 0 {
				path = path + " (empty)"
			} else {
				path = path + " (" + strconv.Itoa(int(stat.Size())) + "b)"
			}
		}

		numberOfSeparators := strings.Count(path, string(os.PathSeparator))
		lastSeparatorIndex := strings.LastIndex(path, string(os.PathSeparator))
		if (numberOfSeparators == 1) {
			nodeStair := path[:lastSeparatorIndex]

			if (nodeStair == prevNodeStair) {
				path = strings.Repeat("\t", numberOfSeparators) + intermediateGraphic + path[lastSeparatorIndex+1:]
			} else {
				path = strings.Repeat("\t", numberOfSeparators) + endGraphic + path[lastSeparatorIndex+1:]
			}
			paths = append(paths, path);
			prevNodeStair = nodeStair
			return nil;
		}

		if (numberOfSeparators == 0) {
			var nodeStair string
			if info.IsDir() {
				nodeStair = path
			}

			if (nodeStair == prevNodeStair) {
				path = endGraphic + path
			} else {
				path = intermediateGraphic + path
			}
			paths = append(paths, path);
			prevNodeStair = nodeStair
			return nil;
		}

		preLastSeparatorIndex := strings.LastIndex(path[:lastSeparatorIndex], string(os.PathSeparator))
		nodeStair := path[preLastSeparatorIndex + 1 : lastSeparatorIndex]
		if (nodeStair == prevNodeStair) {
			path = strings.Repeat("\t", numberOfSeparators) +  intermediateGraphic + path[lastSeparatorIndex+1:]
		} else {
			path = strings.Repeat("\t", numberOfSeparators - 1) + preGraphic + "\t" + endGraphic + path[lastSeparatorIndex+1:]
		}

		paths = append(paths, path);
		prevNodeStair = nodeStair
		return nil;
	})
	if err != nil {
		panic(err.Error());
	}

	if printFiles {
		for i := 0; i < len(paths); i++ {
			fmt.Fprintln(out, paths[i]);
		}
	}

	return nil;
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
