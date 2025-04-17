package grammar

// maps non-terminal -> symbol -> non-terminal -> right hand side
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
			firstSet := grammar.FirstOfSequence(rightHandSide)

			// 1. table[A][t] = A -> rightHandSide for all t in first(righthandside)
			for t := range firstSet {
				if t == "" {
					continue
				}
				if parseTable[nonTerminal][t] == nil {
					parseTable[nonTerminal][t] = make(map[string][]string)
				}
				parseTable[nonTerminal][t][nonTerminal] = rightHandSide
			}

			// 2. If "" in first(A), add A -> righthandside into table[A][t] for all t in follow(A)
			if firstSet[""] {
				followSet := follow[nonTerminal]
				for t := range followSet {
					if parseTable[nonTerminal][t] == nil {
						parseTable[nonTerminal][t] = make(map[string][]string)
					}
					parseTable[nonTerminal][t][nonTerminal] = rightHandSide
				}
			}
		}
	}
	return parseTable, nil
}
