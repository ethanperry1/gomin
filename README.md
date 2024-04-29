# 🏋️‍♀️ GOmin 🏋️

Raise the bar for your code quality and unit test coverage.

## Example Usage

```sh
go get github.com/ethanperry1/gomin
```

```go
results, err := evaluator.Evaluate(
	0.4,
	v0.Min( // Rule #1
		0.9,
		v0.AllFiles(),
	),
	v0.Min( // Rule #2
		0.1,
		v0.AllPackages().Filter("pkg/").Functions().Filter("New"),
	),
	v0.Fallback( // Rule #3
		0.95,
		v0.AllFunctions(),
	),
	v0.Exclude(
		v0.AllPackages().Filter("v0"),
		v0.AllPackages().Filter("visitor"),
		v0.AllPackages().Filter("declarations"),
	),
)
```

What is this doing?
1. Validating that all files have at least 90% unit test coverage.
2. Validating that functions which matches "New" in any package which matches "pkg/" have at least 10% coverage.
3. For all functions which have not yet had rules applied, validate that these have at least 95% coverage.
4. Exclude any packages matching v0, visitor, or declarations from the global unit test coverage calculation.

In this example, the global coverage threshold was 40%, which is set with the first parameter.

[View an example generated report.](./coverage_report.md). Reports will begin with a green circle when successful and red in failure conditions.
