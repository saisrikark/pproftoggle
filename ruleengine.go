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
