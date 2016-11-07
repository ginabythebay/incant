package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func incantFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	result := filepath.Join(usr.HomeDir, "incantations.txt")
	_, err = os.Stat(result)
	return result, err
}

func main() {
	app := cli.NewApp()
	app.Name = "incant"
	app.Usage = "Searches for entries in ~/incantations.txt."
	app.Action = run

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	inPath, err := incantFile()
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	f, err := os.Open(inPath)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
	defer f.Close()

	if len(ctx.Args()) != 1 {
		panic(fmt.Sprintf("Error.  We expect exactly one argument, but we got %s", ctx.Args()))
	}

	search := ctx.Args().First()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		line := scanner.Text()
		if strings.Contains(line, search) {
			fmt.Println(line)
		}
	}

	return nil
}
