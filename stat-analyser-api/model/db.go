package model

import "fmt"

func InitDB() error {
	// Implement / Integrate CS:GO Demo DB
	return nil
}

// Incoming verdicts can be used here
func InsertVerdict(verdict Verdict) error {
	fmt.Print(verdict)
	return nil
}

func OutputVerdicts(SteamdID uint64) error {
	return nil
}
