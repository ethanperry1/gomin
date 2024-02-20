package main

import (
	"fmt"
	"os"

	"gobar/pkg/evaluate"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

//gobar:min:1.0
func run() error {
	root := os.Getenv("ROOT")
	profile := os.Getenv("PROFILE")

	evaluator := evaluate.New("gobar", root, profile)
	cov, err := evaluator.EvalCoverage()
	if err != nil {
		return err
	}

	fmt.Printf("Evaluated overall coverage of %0.2f%% -- %d statements covered of %d total statements.", float64(cov.Covered()) / float64(cov.Statements()), cov.Covered(), cov.Statements())

	return nil
}

