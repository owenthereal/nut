package main

import (
	"fmt"
	"os"

	"github.com/gophergala/nut/vendor/_nuts/github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "nut"
	app.Usage = "Vendor Go dependencies"
	app.Version = "0.0.1"
	app.Author = ""
	app.Email = ""

	app.Commands = []cli.Command{
		newCmd,
		installCmd,
	}

	app.Run(os.Args)
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
