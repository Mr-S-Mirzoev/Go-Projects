package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func printer(out io.Writer, name string, k int, isDir bool, isLastIter []bool, size int64, printFiles bool) {
	//fmt.Println(k, len(isLastIter), isLastIter[k])
	for i := 0; i < k; i++ {
		if isLastIter[i] {
			fmt.Fprint(out, "\t")
		} else {
			fmt.Fprint(out, "│\t")
		}
	}
	if isLastIter[k] {
		fmt.Fprint(out, "└───")
	} else {
		fmt.Fprint(out, "├───")
	}
	if !isDir && printFiles {
		fmt.Fprint(out, name)
		if size == 0 {
			fmt.Fprint(out, " (empty)")
		} else {
			fmt.Fprint(out, " (", size, "b)")
		}

	} else if isDir {
		fmt.Fprint(out, name)
	}
	fmt.Fprintln(out)
}

func dirTreeRec(out io.Writer, path string, printFiles bool, k int, isLastIter []bool) error {
	filesUnfiltered, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	newIsLastIter := make([]bool, len(isLastIter))
	files := make([]os.FileInfo, 0)
	if !printFiles {
		for _, file := range filesUnfiltered {
			newPath := path + "/" + file.Name()
			fileInfo, err := os.Stat(newPath)
			if err != nil {
				log.Fatal(err)
				return err
			}
			if fileInfo.IsDir() {
				files = append(files, file)
			}
		}
	} else {
		files = filesUnfiltered
	}
	for x, file := range files {
		newPath := path + "/" + file.Name()
		fileInfo, err := os.Stat(newPath)
		if err != nil {
			log.Fatal(err)
			return err
		}
		isLast := x == len(files)-1
		newIsLastIter = make([]bool, len(isLastIter))
		copy(newIsLastIter, isLastIter)
		newIsLastIter = append(newIsLastIter, isLast)
		printer(out, file.Name(), k, fileInfo.IsDir(), newIsLastIter, fileInfo.Size(), printFiles)
		if fileInfo.IsDir() {
			newErr := dirTreeRec(out, newPath, printFiles, k+1, newIsLastIter)
			if newErr != nil {
				return newErr
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeRec(out, path, printFiles, 0, make([]bool, 0))
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	/*fileInfo, err = os.Stat("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File name:", fileInfo.Name())
	fmt.Println("Size in bytes:", fileInfo.Size())
	fmt.Println("Permissions:", fileInfo.Mode())
	fmt.Println("Last modified:", fileInfo.ModTime())
	fmt.Println("Is Directory: ", fileInfo.IsDir())
	fmt.Printf("System interface type: %T\n", fileInfo.Sys())
	fmt.Printf("System info: %+v\n\n", fileInfo.Sys()) */
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
