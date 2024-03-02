package main

import (
	"github.com/ethanperry1/gomin/v0"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	evaluator, err := v0.CreateEvaluator(".", "profile", "github.com/ethanperry1/gomin")
	if err != nil {
		return err
	}

	results, err := evaluator.Evaluate(
		0.4,
		v0.Min(
			0.9,
			v0.AllFiles(),
		),
		v0.Min(
			0.0,
			v0.AllPackages().Filter("pkg").Functions().Filter("New"),
		),
		v0.Fallback(
			0.99,
			v0.AllFunctions(),
		),
		v0.Exclude(
			v0.AllPackages().Filter("v0"),
			v0.AllPackages().Filter("visitor"),
			v0.AllPackages().Filter("declarations"),
		),
	)
	if err != nil {
		return err
	}

	out := v0.NewMarkdownOutputter()
	return out.WriteOut(results)
}
