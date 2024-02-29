
```go
evalCoverage(
    0.9,
    Fallbacks(0.9, 0.8, 0.7),
    Packages()
        .Min(0.9),
    Package("tokens")
        .Functions()
        .Min(0.25)
    Files()
        .Regex("_gen.go")
        .Exclude(),
    Package("pkg/evaluate")
        .File("errors.go")
        .Function("InvalidPackageDirectiveError", "Error")
        .Exclude(),
    Package("pkg/evaluate")
        .File("errors.go")
        .Function("swap")
        .Min(0.1),
    Package("pkg/processor")
        .Regex("mocks.go")
        .Exclude(),
    Files()
        .Regex("mocks.go")
        .Min(0.7)
)
```