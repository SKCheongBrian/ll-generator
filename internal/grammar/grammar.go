package grammar

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Production struct {
	Lhs string
	Rhs []string
}

type Grammar struct {
	Terminals    map[string]bool       `yaml:"terminals"`
	NonTerminals map[string]bool       `yaml:"nonterminals"`
	Start        string                `yaml:"start"`
	Productions  map[string][][]string `yaml:"productions"`

	// Cache for First and Follow sets and epsilons
	firstCache   map[string]map[string]bool
	followCache  map[string]map[string]bool
	epsilonCache map[string]bool
	augStart     string
}

func LoadGrammar(filename string) (*Grammar, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tempgrammar struct {
		Terminals    []string              `yaml:"terminals"`
		NonTerminals []string              `yaml:"nonterminals"`
		Start        string                `yaml:"start"`
		Productions  map[string][][]string `yaml:"productions"`
	}
	err = yaml.Unmarshal(data, &tempgrammar)
	if err != nil {
		return nil, err
	}

	var grammar Grammar
	grammar.Terminals = make(map[string]bool)
	grammar.NonTerminals = make(map[string]bool)
	for _, terminal := range tempgrammar.Terminals {
		grammar.Terminals[terminal] = true
	}
	for _, nonTerminal := range tempgrammar.NonTerminals {
		grammar.NonTerminals[nonTerminal] = true
	}
	grammar.Start = tempgrammar.Start
	grammar.Productions = tempgrammar.Productions

	newStart := grammar.Start + "'"
	grammar.NonTerminals[newStart] = true
	grammar.Productions[newStart] = [][]string{{grammar.Start, "$"}}

	grammar.Terminals["$"] = true

	grammar.augStart = newStart

	grammar.ComputeFirst()
	grammar.ComputeFollow()

	return &grammar, nil
}

func isEpsilonProduction(rhs []string) bool {
	return len(rhs) == 1 && rhs[0] == ""
}

func (grammar *Grammar) ComputeEpsilon() map[string]bool {
	if grammar.epsilonCache != nil {
		return grammar.epsilonCache
	}

	epsilon := make(map[string]bool)
	var changed bool = true
	for changed {
		changed = false
		for nonTerminal, productions := range grammar.Productions {
			for _, production := range productions {
				if isEpsilonProduction(production) && !epsilon[nonTerminal] {
					epsilon[nonTerminal] = true
					changed = true
				}

				nullable := true
				for _, symbol := range production {
					if !epsilon[symbol] {
						nullable = false
						break
					}
				}
				if nullable && !epsilon[nonTerminal] {
					epsilon[nonTerminal] = true
					changed = true
				}
			}
		}
	}

	grammar.epsilonCache = epsilon
	return epsilon
}

func (grammar *Grammar) ComputeFirst() map[string]map[string]bool {
	if grammar.firstCache != nil {
		return grammar.firstCache
	}

	epsilon := grammar.ComputeEpsilon()

	// First set for terminals is the terminal itself
	first := make(map[string]map[string]bool)
	for terminal := range grammar.Terminals {
		first[terminal] = map[string]bool{terminal: true}
	}

	// Initially empty for non-terminals
	for nonTerminal := range grammar.NonTerminals {
		first[nonTerminal] = make(map[string]bool)
	}

	// Fixpoint iteration to compute First sets
	changed := true
	for changed {
		changed = false
		for nonTerminal, rightHandSides := range grammar.Productions {
			for _, rightHandSide := range rightHandSides {
				for _, symbol := range rightHandSide {
					// union the first of the symbol with the first of the non-terminal
					updated_first_nt := make(map[string]bool)
					for k := range first[nonTerminal] {
						updated_first_nt[k] = true
					}
					for k := range first[symbol] {
						if !first[nonTerminal][k] {
							updated_first_nt[k] = true
							changed = true
						}
					}
					first[nonTerminal] = updated_first_nt

					// if the current symbol cannot create an epsilon
					// then we are done
					if !epsilon[symbol] {
						break
					}
				}
			}
		}
	}

	grammar.firstCache = first
	return first
}

func (grammar *Grammar) ComputeFollow() map[string]map[string]bool {
	if grammar.followCache != nil {
		return grammar.followCache
	}

	epsilon := grammar.ComputeEpsilon()
	first := grammar.ComputeFirst()

	follow := make(map[string]map[string]bool)

	for nonTerminal := range grammar.NonTerminals {
		follow[nonTerminal] = make(map[string]bool)
	}

	changed := true
	for changed {
		changed = false
		for nonTerminal, rightHandSides := range grammar.Productions {
			for _, rightHandSide := range rightHandSides {
				aux := follow[nonTerminal]
				for _, symbol := range Reversed(rightHandSide) {
					if follow[symbol] != nil {
						updated_follow_symbol := make(map[string]bool)
						for k := range follow[symbol] {
							updated_follow_symbol[k] = true
						}
						for k := range aux {
							if !follow[symbol][k] {
								updated_follow_symbol[k] = true
								changed = true
							}
						}
						follow[symbol] = updated_follow_symbol
					}

					if epsilon[symbol] {
						for k := range first[symbol] {
							aux[k] = true
						}
					} else {
						updated_aux := make(map[string]bool)
						for k := range first[symbol] {
							updated_aux[k] = true
						}
						aux = updated_aux
					}
				}
			}
		}
	}

	grammar.followCache = follow
	return follow
}

// Generates the first set for a given sequence of symbols.
func (grammar *Grammar) FirstOfSequence(symbols []string) map[string]bool {
	result := make(map[string]bool)

	if isEpsilonProduction(symbols) {
		result[""] = true
		return result
	}

	epsilon := grammar.ComputeEpsilon()
	first := grammar.ComputeFirst()
	for _, symbol := range symbols {
		for t := range first[symbol] {
			result[t] = true
		}
		if !epsilon[symbol] {
			delete(result, "")
			return result
		}
	}

	result[""] = true
	return result
}

func (grammar *Grammar) CanDeriveEpsilon(symbols []string) bool {
	epsilon := grammar.ComputeEpsilon()
	for _, symbol := range symbols {
		if !epsilon[symbol] {
			return false
		}
	}
	return true
}

func Reversed(arr []string) []string {
	result := make([]string, len(arr))

	for i, j := len(arr)-1, 0; i >= 0; i, j = i-1, j+1 {
		result[j] = arr[i]
	}

	return result
}
