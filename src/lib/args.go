package llpm

import (
	"errors"
	"fmt"
	"os"
)

type ArgParser struct {
	command string
	args    []string
	current int
	err     error
}

func NewArgParser(commands []string) (*ArgParser, error) {
	var result ArgParser
	if len(commands) <= 1 {
		return &result, errors.New(
			"No command provided. Use \"llpm help\" for a list of available commands")
	}
	result.command = commands[1]
	result.current = 0
	args := make([]string, 0)
	if len(commands) >= 3 {
		args = commands[2:]
	}
	result.args = args
	return &result, nil
}

func (argParser *ArgParser) GetCommand() string {
	return argParser.command
}

func (argParser *ArgParser) GetLen() int {
	return len(argParser.args) - argParser.current
}

func (argParser *ArgParser) getNext(argName string) (string, error) {
	if argParser.GetLen() <= 0 {
		return "", errors.New(fmt.Sprintf(
			"Missing argument \"%s\" for command \"%s\"",
			argName, argParser.GetCommand()))
	}
	result := argParser.args[argParser.current]
	argParser.current += 1
	return result, nil
}

func (argParser *ArgParser) GetNext(argName string) string {
	result, err := argParser.getNext(argName)
	if err != nil {
		argParser.err = err
	}
	return result
}

func (argParser *ArgParser) GetNextOr(argName string, defaultValue string) string {
	result, err := argParser.getNext(argName)
	if err != nil {
		return defaultValue
	}
	return result
}

func (argParser *ArgParser) GetRemaining() []string {
	if argParser.GetLen() == 0 {
		return make([]string, 0)
	}
	result := argParser.args[argParser.current:]
	argParser.current += len(result)
	return result
}

func (argParser *ArgParser) ExitIfErr() {
	if argParser.err != nil {
		fmt.Println(argParser.err)
		os.Exit(1)
	}
}
