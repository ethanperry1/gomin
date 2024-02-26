package main

import (
	"fmt"
	"os"
	"strconv"

	"gomin/pkg/evaluate"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	root := os.Getenv("ROOT")
	profile := os.Getenv("PROFILE")
	name := os.Getenv("NAME")
	overall := os.Getenv("OVERALL_MIN_COV")
	defPackage := os.Getenv("DEFAULT_MIN_PACKAGE_COV")
	defFile := os.Getenv("DEFAULT_MIN_PACKAGE_COV")
	defblock := os.Getenv("DEFAULT_MIN_PACKAGE_COV")

	minOverallCov, err := strconv.ParseFloat(overall, 64)
	if err != nil {
		fmt.Println("no valid overall coverage bar supplied, defaulting to 0.0")
	}

	minPackageCov, err := strconv.ParseFloat(defPackage, 64)
	if err != nil {
		fmt.Println("no valid default minimum package coverage supplied, defaulting to 0.0")
	}

	minFileCov, err := strconv.ParseFloat(defFile, 64)
	if err != nil {
		fmt.Println("no valid default minimum file coverage supplied, defaulting to 0.0")
	}

	minBlockCov, err := strconv.ParseFloat(defblock, 64)
	if err != nil {
		fmt.Println("no valid default minimum function block coverage supplied, defaulting to 0.0")
	}

	evaluator := evaluate.New(name, root, profile,
		evaluate.InitDefaultPackageMinimum(minPackageCov),
		evaluate.InitDefaultFileMinimum(minFileCov),
		evaluate.InitDefaultFunctionMinimum(minBlockCov),
	)
	cov, err := evaluator.EvalCoverage()
	if err != nil {
		return err
	}

	overallCoverage := float64(cov.Covered()) / float64(cov.Statements())
	fmt.Printf("Evaluated overall coverage of %0.2f -- %d statement(s) covered of %d total statement(s).\n", overallCoverage, cov.Covered(), cov.Statements())

	if overallCoverage < minOverallCov {
		return fmt.Errorf("expected coverage of at least %0.2f -- coverage bar was not met", minOverallCov)
	}

	return nil
}
