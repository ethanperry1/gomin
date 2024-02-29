package main

import (
	"fmt"
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
	evaluator, err := v0.CreateEvaluator(".", "profile", "github.com/ethanperry1/gomin")
	if err != nil {
		return err
	}

	results, err := evaluator.Evaluate(
		0.4,
		v0.Min(
			0.7,
			v0.Package("pkg/profiles").File("profiles.go").Method("ProfilesByName", "Get"),
			v0.AllPackages(),
		),
	)
	if err != nil {
		return err
	}

	fmt.Println(v0.Validate(results))

	file, err := os.Create("coverage_table.txt")
	if err != nil {
		return err
	}

	writer := v0.NewWriter()
	return writer.Write(file, v0.Format(results, v0.FileDepth))
}
