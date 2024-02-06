package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	ext  string
	size int64
	list bool
}

func main() {
	root := flag.String("root", "", "Root dir to start")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Min file size")
	flag.Parse()

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	return nil
}
