package grammar

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadGrammar(t *testing.T) {
	tempFile, err := os.CreateTemp("", "grammar_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

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
	tempFile.Close()

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
	if len(grammar.Terminals) != len(expectedTerminals) {
		t.Errorf("Expected %d terminals, got %d", len(expectedTerminals), len(grammar.Terminals))
	}
	for terminal := range expectedTerminals {
		if !grammar.Terminals[terminal] {
			t.Errorf("Expected terminal '%s' but it was not found", terminal)
		}
	}

	// Check non-terminals
	expectedNonTerminals := map[string]bool{"S": true, "A": true, "S'": true}
	if len(grammar.NonTerminals) != len(expectedNonTerminals) {
		t.Errorf("Expected %d non-terminals, got %d", len(expectedNonTerminals), len(grammar.NonTerminals))
	}
	for nonTerminal := range expectedNonTerminals {
		if !grammar.NonTerminals[nonTerminal] {
			t.Errorf("Expected non-terminal '%s', but it was not found", nonTerminal)
		}
	}

	// Check productions
	expectedProductions := map[string][][]string{
		"S":  {{"A", "a"}, {"b"}},
		"A":  {{"a"}},
		"S'": {{"S", "$"}},
	}
	for nonTerminal, rules := range expectedProductions {
		if len(grammar.Productions[nonTerminal]) != len(rules) {
			t.Errorf("Expected %d productions for '%s', got %d", len(rules), nonTerminal, len(grammar.Productions[nonTerminal]))
		}
		for i, rule := range rules {
			for j, symbol := range rule {
				if grammar.Productions[nonTerminal][i][j] != symbol {
					t.Errorf("Expected symbol '%s' in production, got '%s'", symbol, grammar.Productions[nonTerminal][i][j])
				}
			}
		}
	}

	fmt.Println(grammar.firstCache)
	fmt.Println(grammar.followCache)
}

func TestComputeEpsilon(t *testing.T) {
	grammar := &Grammar{
		Productions: map[string][][]string{
			"S": {{"A", "a"}, {"B"}},
			"A": {{"a"}, {}},
			"B": {{"b"}, {}},
		},
	}

	epsilon := grammar.ComputeEpsilon()

	if !epsilon["A"] {
		t.Errorf("Expected 'A' to be epsilon, but it was not")
	}
	if !epsilon["S"] {
		t.Errorf("Expected 'S' to be epsilon, but it was not")
	}
	if !epsilon["B"] {
		t.Errorf("Expected 'B' to  be epsilon, but it was not")
	}

	fmt.Println(epsilon)
}
