package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

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
	app.Usage = "Searches for entries in ~/incantations.txt, using simple substring matching."
	app.UsageText = "incant [options] <searchterm>"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "r, run",
			Usage: "If run is set, and there is exactly a one line match, we will exec $SHELL and run that line in it.",
		},
	}

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

	var line string
	var matchCount int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		line = scanner.Text()
		if strings.Contains(line, search) {
			fmt.Println(line)
			matchCount++
		}
	}

	if ctx.Bool("run") {
		if matchCount == 0 {
			fmt.Println("Ignoring run option as there was no matching line.")
			return nil
		}
		if matchCount != 1 {
			fmt.Printf("Ignoring run option as it only works with 1 match line and there were %d matching lines.\n", matchCount)
			os.Exit(1)
		}

		doRun(line)
	}

	return nil
}

func doRun(line string) {
	shell, found := os.LookupEnv("SHELL")
	if !found {
		fmt.Println("Unable to run.  Set SHELL environment first.")
		os.Exit(1)
	}
	args := []string{shell, "-c", line}
	env := os.Environ()
	err := syscall.Exec(shell, args, env)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
