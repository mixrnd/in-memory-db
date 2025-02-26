package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_ToSting(t *testing.T) {
	tests := []struct {
		name           string
		query          Query
		expectedString string
	}{
		{
			name: "get query",
			query: Query{
				command: GetCommand,
				args:    []string{"config"},
			},
			expectedString: "GET config",
		},
		{
			name: "set query",
			query: Query{
				command: SetCommand,
				args:    []string{"config", "123"},
			},
			expectedString: "SET config 123",
		},
		{
			name: "del query",
			query: Query{
				command: DelCommand,
				args:    []string{"config"},
			},
			expectedString: "DEL config",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedString, test.query.ToSting())
		})
	}
}
