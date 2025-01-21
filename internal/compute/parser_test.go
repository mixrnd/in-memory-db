package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name          string
		cmd           string
		expectedQuery Query
		expectedError error
	}{
		{
			name: "correct get query",
			cmd:  "GET config",
			expectedQuery: Query{
				command: GetCommand,
				args:    []string{"config"},
			},
		},
		{
			name: "correct get query, args separated by tab",
			cmd:  "GET	config",
			expectedQuery: Query{
				command: GetCommand,
				args:    []string{"config"},
			},
		},
		{
			name:          "incorrect get query, no args",
			cmd:           "GET ",
			expectedError: ErrWrongArgumentNumber,
		},
		{
			name: "correct set query",
			cmd:  "SET config 123",
			expectedQuery: Query{
				command: SetCommand,
				args:    []string{"config", "123"},
			},
		},
		{
			name: "correct set query, args separated by tabs and spaces",
			cmd:  "SET	 config   	123",
			expectedQuery: Query{
				command: SetCommand,
				args:    []string{"config", "123"},
			},
		},
		{
			name: "correct del query",
			cmd:  "DEL config",
			expectedQuery: Query{
				command: DelCommand,
				args:    []string{"config"},
			},
		},
		{
			name:          "incorrect del query, no args",
			cmd:           "DEL ",
			expectedError: ErrWrongArgumentNumber,
		},
		{
			name: "correct del query, lower",
			cmd:  "del config",
			expectedQuery: Query{
				command: DelCommand,
				args:    []string{"config"},
			},
		},
		{
			name:          "unknown command",
			cmd:           "PUT a b",
			expectedError: ErrUnknownCommand,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewParser()
			q, err := p.Parse(test.cmd)

			if test.expectedError != nil {
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedQuery, q)
			}
		})
	}
}
