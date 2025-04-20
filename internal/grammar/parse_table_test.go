package grammar

import (
	"fmt"
	"testing"
)

var exampleGrammar Grammar = Grammar{
	Start: "S",
	NonTerminals: map[string]bool{
		"S":  true,
		"E":  true,
		"T":  true,
		"T'": true,
		"F":  true,
		"S'": true, // augmented start
	},
	Terminals: map[string]bool{
		"(":  true,
		")":  true,
		"id": true,
		"*":  true,
		"+":  true,
		"$":  true, // end of input marker
	},
	Productions: map[string][]Production{
		"S": {
			{
				Lhs: "S",
				Rhs: []string{
					"T", "E",
				},
			},
		},

		"E": {
			{
				Lhs: "E",
				Rhs: []string{
					"+", "T", "E",
				},
			},
			{
				Lhs: "E",
				Rhs: []string{
					"",
				},
			},
		},

		"T": {
			{
				Lhs: "T",
				Rhs: []string{
					"F", "T'",
				},
			},
		},

		"T'": {
			{
				Lhs: "T'",
				Rhs: []string{
					"*", "F", "T'",
				},
			},
			{
				Lhs: "T'",
				Rhs: []string{
					"",
				},
			},
		},

		"F": {
			{
				Lhs: "F",
				Rhs: []string{
					"id",
				},
			},
			{
				Lhs: "F",
				Rhs: []string{
					"(", "S", ")",
				},
			},
		},

		"S'": {
			{Lhs: "S'", Rhs: []string{"S", "$"}},
		},
	},
	augStart: "S'",
}

func TestParseTableGeneration(t *testing.T) {
	var parseTable ParseTable = exampleGrammar.GenerateParseTable()
	for nt, entry := range parseTable {
		fmt.Printf("Non-terminal: %s\n", nt)
		for t, prod := range entry {
			fmt.Printf("    Lookahead: %s |-> %s --> %v\n", t, prod.Lhs, prod.Rhs)
		}
	}
}
