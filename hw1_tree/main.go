package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	// "strings"
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
	file, err := os.Open(path)
	if err != nil {
		return err 
	}
	st, err := file.Stat();
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return nil
	}
	entries, err := file.Readdir(-1);
	if err != nil {
		return err
	}
	for i, entry := range entries {
		err = _dirTree(out, path + "/" + entry.Name(), printFiles, "", i == len(entries) - 1)
		if err != nil {
			return err
		}
	}
	return nil
}

func _dirTree(out io.Writer, path string, printFiles bool, prefix string, last bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err 
	}
	st, err := file.Stat();
	if err != nil {
		return err
	}
	if !st.IsDir() {
		if printFiles {
			printEntry(out, path, prefix, last)
		}
		return nil
	}
	printEntry(out, st.Name(), prefix, last)
	entries, err := file.Readdir(-1);
	if err != nil {
		return err
	}
	for i, entry := range entries {
		if last {
			err = _dirTree(out, path + "/" + entry.Name(), printFiles, prefix + "	", i == len(entries) - 1)
		} else {
			err = _dirTree(out, path + "/" + entry.Name(), printFiles, prefix + "│\t", i == len(entries) - 1)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func printEntry(out io.Writer, filename string, prefix string, last bool) {
	prefix = fprefix(prefix, last)
	out.Write([]byte(prefix + filepath.Base(filename) + "\n"))
}

func printFile(out io.Writer, filename string, size int, prefix string, last bool) {
	prefix = fprefix(prefix, last)
	s := prefix + filepath.Base(filename)
	if size > 0 {
		s += fmt.Sprintf(" (%db)\n", size)
	} else {
		s += fmt.Sprint(" (empty)\n")
	}
	out.Write([]byte(s))
}

func fprefix(prefix string, last bool) string {
	if last {
		return prefix + "└───"
	} else {
		return prefix + "├───"
	}
}
