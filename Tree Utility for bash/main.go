package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func dirTreeFmt(out io.Writer, path string, printFiles bool, pos int, prefix, prefixForLast string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	num := len(files)

	if printFiles {
		for pos, file := range files {
			if file.IsDir() && pos != (num-1) {
				fmt.Printf("%s├───%s\n", prefix, file.Name())

				prevPref := prefix
				prevPrefForLast := prefixForLast
				prefix += "|\t"
				if pos == num-2 {
					prefixForLast += "|\t"
				} else {
					prefixForLast += " \t"
				}
				fileName := path + "/" + file.Name() + "/"

				dirTreeFmt(out, fileName, true, pos, prefix, prefix)

				prefix = prevPref
				prefixForLast = prevPrefForLast
			} else if file.IsDir() && pos == (num-1) {
				fmt.Printf("%s└───%s\n", prefixForLast, file.Name())

				prevPref := prefix
				prevPrefForLast := prefixForLast
				prefix += "|\t"
				prefixForLast += " \t"
				fileName := path + "/" + file.Name() + "/"

				dirTreeFmt(out, fileName, true, pos, prefixForLast, prefixForLast)

				prefix = prevPref
				prefixForLast = prevPrefForLast
			} else {
				if pos != (num - 1) {
					if file.Size() == 0 {
						fmt.Printf("%s├───%s (empty)\n", prefix, file.Name())
					} else {
						fmt.Printf("%s├───%s (%db)\n", prefix, file.Name(), file.Size())
					}
				} else if file.Size() == 0 {
					fmt.Printf("%s└───%s (empty)\n", prefixForLast, file.Name())
				} else {
					fmt.Printf("%s└───%s (%db)\n", prefix, file.Name(), file.Size())
				}
			}
		}

		return err
	} else {
		num = 0
		pos := 0

		for _, file := range files {
			if file.IsDir() {
				num += 1
			}
		}

		for _, file := range files {
			if file.IsDir() && pos != (num-1) {
				pos += 1

				fmt.Printf("%s├───%s\n", prefix, file.Name())

				prevPref := prefix
				prevPrefForLast := prefixForLast
				prefix += "|\t"
				if pos == num-2 {
					prefixForLast += "|\t"
				} else {
					prefixForLast += " \t"
				}
				fileName := path + "/" + file.Name() + "/"

				dirTreeFmt(out, fileName, false, pos, prefix, prefix)

				prefix = prevPref
				prefixForLast = prevPrefForLast
			} else if file.IsDir() && pos == (num-1) {
				fmt.Printf("%s└───%s\n", prefixForLast, file.Name())

				prevPref := prefix
				prevPrefForLast := prefixForLast
				prefix += "|\t"
				prefixForLast += " \t"
				fileName := path + "/" + file.Name() + "/"

				dirTreeFmt(out, fileName, false, pos, prefixForLast, prefixForLast)

				prefix = prevPref
				prefixForLast = prevPrefForLast
			}
		}

		return err
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	_, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fileName := "./" + path + "/"

	return dirTreeFmt(out, fileName, printFiles, 0, "", "")
}

func main() {
	var out io.Writer
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]

	fmt.Println(path)
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
