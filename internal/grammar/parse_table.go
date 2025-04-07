package grammar

type ParseTable map[string]map[string]map[string][]string

func (grammar *Grammar) GenerateParseTable() (ParseTable, err) {
	first := grammar.ComputeFirst()
	follow := grammar.ComputeFollow()

	parseTable := make(ParseTable)
	for nonTerminal := range grammar.NonTerminals {
		parseTable[nonTerminal] = make(map[string]map[string][]string)
	}

	for nonTerminal, rightHandSides := range grammar.Productions {
		for _, rightHandSide := range rightHandSides {

		}
	}
}
