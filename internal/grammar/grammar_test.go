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

	if grammar.Start != "S" {
		t.Errorf("Expected start symbol 'S', got '%s'", grammar.Start)
	}

	expectedTerminals := []string{"a", "b"}
	if len(grammar.Terminals) != len(expectedTerminals) {
		t.Errorf("Expected %d terminals, got %d",
			len(expectedTerminals), len(grammar.Terminals))
	}
	for i, terminal := range expectedTerminals {
		if grammar.Terminals[i] != terminal {
			t.Errorf("Expected terminal '%s', got '%s'",
				terminal, grammar.Terminals[i])
		}
	}

	expectedNonTerminals := []string{"S", "A"}
	if len(grammar.NonTerminals) != len(expectedNonTerminals) {
		t.Errorf("Expected %d nonterminals, got %d", len(expectedNonTerminals), len(grammar.NonTerminals))
	}
	for i, nonTerminal := range expectedNonTerminals {
		if grammar.NonTerminals[i] != nonTerminal {
			t.Errorf("Expected nonterminal '%s', got '%s'", nonTerminal, grammar.NonTerminals[i])
		}
	}

	expectedProductions := map[string][][]string{
		"S": {{"A", "a"}, {"b"}},
		"A": {{"a"}},
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
