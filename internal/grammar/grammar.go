package grammar

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Grammar struct {
	Terminals    []string              `yaml:"terminals"`
	NonTerminals []string              `yaml:"nonterminals"`
	Start        string                `yaml:"start"`
	Productions  map[string][][]string `yaml:"productions"`

	// Cache for First and Follow sets
	firstCache  map[string]map[string]struct{}
	followCache map[string]map[string]struct{}
}

func LoadGrammar(filename string) (*Grammar, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var grammar Grammar
	err = yaml.Unmarshal(data, &grammar)
	if err != nil {
		return nil, err
	}

	grammar.ComputeFirst()
	grammar.ComputeFollow()

	return &grammar, nil
}

func (g *Grammar) ComputeFirst() map[string]map[string]struct{} {
	if g.firstCache != nil {
		return g.firstCache
	}

	first := map[string]map[string]struct{}{}
	for _, terminal := range g.Terminals {
		first[terminal] = map[string]struct{}{terminal: {}}
	}

	for _, nonTerminal := range g.NonTerminals {
		first[nonTerminal] = map[string]struct{}{}
	}

	changed := true
	for changed {
		changed = false
		for nonTerminal, productions := range g.Productions {
			for _, production := range productions {
				for _, symbol := range production {
					if _, isTerminal := first[symbol]; isTerminal {
						if _, exists := first[nonTerminal][symbol]; !exists {
							first[nonTerminal][symbol] = struct{}{}
							changed = true
						}
						break
					} else {
						for terminal := range first[symbol] {
							if _, exists := first[nonTerminal][terminal]; !exists {
								first[nonTerminal][terminal] = struct{}{}
								changed = true
							}
						}
						if _, exists := first[symbol][""]; !exists {
							break
						}
					}
				}
			}
		}
	}

	g.firstCache = first
	return first
}

func (g *Grammar) ComputeFollow() map[string]map[string]struct{} {
	if g.followCache != nil {
		return g.followCache
	}

	follow := map[string]map[string]struct{}{}
	for _, nonTerminal := range g.NonTerminals {
		follow[nonTerminal] = map[string]struct{}{}
	}

	follow[g.Start]["$"] = struct{}{}

	changed := true
	for changed {
		changed = false
		for nonTerminal, productions := range g.Productions {
			for _, production := range productions {
				for i, symbol := range production {
					if _, isNonTerminal := follow[symbol]; isNonTerminal {
						if i+1 < len(production) {
							nextSymbol := production[i+1]
							if _, isTerminal := follow[nextSymbol]; isTerminal {
								if _, exists := follow[symbol][nextSymbol]; !exists {
									follow[symbol][nextSymbol] = struct{}{}
									changed = true
								}
							} else {
								for terminal := range g.firstCache[nextSymbol] {
									if terminal != "" && terminal != "$" {
										if _, exists := follow[symbol][terminal]; !exists {
											follow[symbol][terminal] = struct{}{}
											changed = true
										}
									}
								}
								if _, exists := g.firstCache[nextSymbol][""]; exists {
									if _, exists := follow[symbol]["$"]; !exists {
										follow[symbol]["$"] = struct{}{}
										changed = true
									}
								}
							}
						} else if _, exists := follow[nonTerminal]["$"]; exists {
							if _, exists := follow[symbol]["$"]; !exists {
								follow[symbol]["$"] = struct{}{}
								changed = true
							}
						}
					}
				}
			}
		}
	}

	g.followCache = follow
	return follow
}
