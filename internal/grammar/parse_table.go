package grammar

// maps non-terminal -> symbol -> non-terminal -> right hand side
type ParseTable map[string]map[string]Production

func (grammar *Grammar) GenerateParseTable() ParseTable {
	grammar.ComputeFirst()
	follow := grammar.ComputeFollow()

	parseTable := make(ParseTable)
	for nonTerminal := range grammar.NonTerminals {
		parseTable[nonTerminal] = make(map[string]Production)
	}

	for nonTerminal, productions := range grammar.Productions {
		for _, production := range productions {
			firstSet := grammar.FirstOfSequence(production.Rhs)

			// 1. table[A][t] = A -> rightHandSide for all t in first(righthandside)
			for t := range firstSet {
				if t == "" {
					continue
				}
				parseTable[nonTerminal][t] = production
			}

			// 2. If "" in first(A), add A -> righthandside into table[A][t] for all t in follow(A)
			if firstSet[""] {
				followSet := follow[nonTerminal]
				for t := range followSet {
					parseTable[nonTerminal][t] = production
				}
			}
		}
	}
	return parseTable
}
