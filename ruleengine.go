package pproftoggle

type Rule interface {
	// Name returns the name of the rule
	Name() string
	// Matches determines whether the
	Matches() (bool, error)
}

type status struct {
	hasMatched   bool
	rulesMatched []Rule
}

func getStatus(rules []Rule) (status, error) {
	var st status
	st.rulesMatched = make([]Rule, 0)

	for _, rule := range rules {
		matches, err := rule.Matches()
		if err != nil {
			return st, err
		} else if matches {
			st.hasMatched = true
			st.rulesMatched = append(st.rulesMatched, rule)
		}
	}

	return st, nil
}
