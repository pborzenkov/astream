package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "AlfaStream (https://potok.digital) report processing tool.")
	app.HelpFlag.Short('h')

	byOpCmd := addByOpCommand(app)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case byOpCmd.FullCommand():
		exit(byOp())
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
