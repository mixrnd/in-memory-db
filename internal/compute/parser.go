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

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(cmd string) (Query, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 || len(parts) > 3 {
		return Query{}, ErrUnknownCommand
	}

	command := Command(strings.ToUpper(parts[0]))
	switch command {
	case GetCommand, DelCommand:
		if len(parts) != 2 {
			return Query{}, ErrWrongArgumentNumber
		}
		return NewQuery(command, []string{parts[1]}), nil
	case SetCommand:
		if len(parts) != 3 {
			return Query{}, ErrWrongArgumentNumber
		}
		return NewQuery(command, parts[1:3]), nil
	}

	return Query{}, ErrUnknownCommand
}
