package main

import (
	"os"

	"github.com/ethanperry1/gomin/v0"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	evaluator, err := api.CreateEvaluator(".", "profile", "github.com/ethanperry1/gomin")
	if err != nil {
		return err
	}

	results, err := evaluator.Evaluate(
		api.Min(
			0.7,
			api.Package("pkg/profiles").File("profiles.go").Method("ProfilesByName", "Get"),
			api.AllPackages(),
		),
	)
	if err != nil {
		return err
	}

	file, err := os.Create("coverage_table.txt")
	if err != nil {
		return err
	}

	writer := api.NewWriter()
	return writer.Write(file, api.Format(results, api.FileDepth))
}
