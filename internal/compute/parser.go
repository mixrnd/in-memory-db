package compute

import (
	"errors"
	"strings"
)

var (
	ErrUnknownCommand      = errors.New("unknown command")
	ErrWrongArgumentNumber = errors.New("wrong argument number")
)

type Command string

const (
	GetCommand Command = "GET"
	SetCommand Command = "SET"
	DelCommand Command = "DEL"
)

const (
	maxCommandParts          = 3
	commandIndex             = 0
	firstArgIndex            = 1
	secondArgIndex           = 2
	oneCommandArgPartNumber  = 2
	twoCommandArgsPartNumber = 3
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(cmd string) (Query, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 || len(parts) > maxCommandParts {
		return Query{}, ErrUnknownCommand
	}

	command := Command(strings.ToUpper(parts[commandIndex]))
	switch command {
	case GetCommand, DelCommand:
		if len(parts) != oneCommandArgPartNumber {
			return Query{}, ErrWrongArgumentNumber
		}
		return NewQuery(command, []string{parts[firstArgIndex]}), nil
	case SetCommand:
		if len(parts) != twoCommandArgsPartNumber {
			return Query{}, ErrWrongArgumentNumber
		}
		return NewQuery(command, parts[firstArgIndex:secondArgIndex+1]), nil
	}

	return Query{}, ErrUnknownCommand
}
