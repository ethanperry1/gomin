package main

import (
	"fmt"

	"github.com/ethanperry1/gomin/pkg/api"
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
			api.Package("pkg/profiles").File("profiles.go").Method("ProfilesByName", "New"),
		),
	)
	if err != nil {
		return err
	}

	fmt.Println(results)

	return nil
}
