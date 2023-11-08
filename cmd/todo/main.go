package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
	"errors"

	gotodoapp "github.com/ddl-killer/go-todo-app"
)

const (
	todoFile = ".todos.json"
)

func init() {
	add = flag.Bool("add", false, "add a new todo")
	complete = flag.Int("complete", -1, "mark a todo as completd")
	delete = flag.Int("delete", -1, "delete a todo")
	list = flag.Bool("list", false, "list all todos")

	flag.Parse()
}

var (
	add      *bool
	complete *int
	delete   *int
	list     *bool
)

func main() {
	todos := gotodoapp.Todos{}
	err := todos.Load(todoFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task ,err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		err = todos.Save(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Save(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *delete > 0:
		err := todos.Delete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Save(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
    case *list:
		todos.Print()
	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	text := scanner.Text()
	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}
	return text,nil
}
