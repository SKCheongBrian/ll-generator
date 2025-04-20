package grammar

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadGrammar(t *testing.T) {
	tempFile, err := os.CreateTemp("", "grammar_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	yamlContent := `
terminals:
  - a
  - b
nonterminals:
  - S
  - A
start: S
productions:
  S:
    - [A, a]
    - [b]
  A:
    - [a]
`
	if _, err := tempFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	grammar, err := LoadGrammar(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load grammar: %v", err)
	}

	// Check start symbol
	if grammar.Start != "S" {
		t.Errorf("Expected start symbol 'S', got '%s'", grammar.Start)
	}

	// Check terminals
	expectedTerminals := map[string]bool{"a": true, "b": true, "$": true}
	if !reflect.DeepEqual(grammar.Terminals, expectedTerminals) {
		t.Errorf("Expected terminals %v, got %v", expectedTerminals, grammar.Terminals)
	}

	// Check non-terminals
	expectedNonTerminals := map[string]bool{"S": true, "A": true, "S'": true}
	if !reflect.DeepEqual(grammar.NonTerminals, expectedNonTerminals) {
		t.Errorf("Expected non-terminals %v, got %v", expectedNonTerminals, grammar.NonTerminals)
	}

	expectedProductions := map[string][]Production{
		"S": {
			{
				Lhs: "S",
				Rhs: []string{
					"A", "a",
				},
			},
			{
				Lhs: "S",
				Rhs: []string{
					"b",
				},
			},
		},

		"A": {
			{
				Lhs: "A",
				Rhs: []string{
					"a",
				},
			},
		},

		"S'": {
			{
				Lhs: "S'",
				Rhs: []string{
					"S", "$",
				},
			},
		},
	}

	if !reflect.DeepEqual(grammar.Productions, expectedProductions) {
		t.Errorf("Expected productions %v, got %v", expectedProductions, grammar.Productions)
	}
}

func TestFollow1(t *testing.T) {
	grammar := &Grammar{
		Terminals: map[string]bool{
			"a": true, "b": true, "+": true, "*": true, "$": true,
		},
		NonTerminals: map[string]bool{
			"E": true, "E'": true, "T": true, "T'": true, "F": true,
		},
		Start: "E",
		Productions: map[string][]Production{
			"E":  {{Lhs: "E", Rhs: []string{"T", "E'"}}},
			"E'": {{Lhs: "E'", Rhs: []string{"+", "T", "E'"}}, {Lhs: "E'", Rhs: []string{""}}},
			"T":  {{Lhs: "T", Rhs: []string{"F", "T'"}}},
			"T'": {{Lhs: "T'", Rhs: []string{"*", "F", "T'"}}, {Lhs: "T'", Rhs: []string{""}}},
			"F":  {{Lhs: "F", Rhs: []string{"a"}}, {Lhs: "F", Rhs: []string{"b"}}},
		},
		augStart: "S",
	}

	// Add augmented start production
	grammar.NonTerminals["S"] = true
	grammar.Productions["S"] = append(grammar.Productions["S"], Production{
		Lhs: "S", Rhs: []string{"E", "$"},
	})
	grammar.Terminals["$"] = true

	follow := grammar.ComputeFollow()

	expected := map[string]map[string]bool{
		"E":  {"$": true},
		"E'": {"$": true},
		"T":  {"+": true, "$": true},
		"T'": {"+": true, "$": true},
		"F":  {"*": true, "+": true, "$": true},
	}

	for nt, expectedSet := range expected {
		actualSet := follow[nt]
		if !reflect.DeepEqual(actualSet, expectedSet) {
			t.Errorf("FOLLOW(%s): expected %v, got %v", nt, expectedSet, actualSet)
		}
	}
}

func TestComputeEpsilon(t *testing.T) {
	grammar := &Grammar{
		Productions: map[string][]Production{
			"S": {
				{
					Lhs: "S",
					Rhs: []string{"A", "a"},
				},
				{
					Lhs: "S",
					Rhs: []string{"B"},
				},
			},
			"A": {
				{
					Lhs: "A",
					Rhs: []string{"a"},
				},
				{
					Lhs: "A",
					Rhs: []string{""},
				},
			},
			"B": {
				{
					Lhs: "B",
					Rhs: []string{"b"},
				},
				{
					Lhs: "B",
					Rhs: []string{""},
				},
			},
		},
	}

	epsilon := grammar.ComputeEpsilon()

	expectedEpsilon := map[string]bool{"A": true, "B": true, "S": true}
	if !reflect.DeepEqual(epsilon, expectedEpsilon) {
		t.Errorf("Expected epsilon %v, got %v", expectedEpsilon, epsilon)
	}
}
