# 🏋️‍♀️ GObar 🏋️

Raise the bar for your code quality and unit test coverage.

_Warning:_ This library is still under development and may contain unresolved bugs.

## Usage

Gobar provides a number of special directives which can be embedded directly in code through comments. These directives shape the way in which unit test coverage is evaluated. For example, we can expect the unit test coverage of a particular function to be at least 65% so long as the gobar directive is applied:

```go
//gobar:min:0.65
func MyFunction() {
    if true {
        fmt.Println("Hello world!")
    }
}
```

## Directives

__Types__

1. _min_ | verifies that the block, file, or package has a minimum number of covered statements.
2. _exclude_ | excludes a block, file, or package from overall unit test coverage calculations.
3. _pkg_ | denotes that the directive should affect the package level.
4. _default_ | can be used to set default behaviors for calculating coverage on files or blocks.

__Formatting__

```go
//gobar:<directive_name>:<directive_param(s)>
```

__Placement__

```go
// Package level directive.
//gobar:pkg:min:0.5
package mypack

//-- or --

// File level directive.
//gobar:min:0.5
package mypack

//-- or --

// Function level directive.
//gobar:min:0.5
func MyFunction() {...}

//-- or --

// Function variable level directive.
//gobar:min:0.5
var myFunc = func() {...}
```

__Running__

As executable:

```sh
# Profile must be present:
# go test ./... -coverprofile=profile

# Optional parameters:
# export OVERALL_MIN_COV=0.8
# export DEFAULT_MIN_PACKAGE_COV=0.7
# export DEFAULT_MIN_FILE_COV=0.6
# export DEFAULT_MIN_BLOCK_COV=0.5

export ROOT=. # Root directory of project.
export PROFILE=profile # Path to profile.
export NAME=gobar # Name of project in go.mod
go run ./cmd
```

With evaluation API:

```go
import "gobar/pkg/evaluate"

const (
    name="gobar"
    root="/home/user/work/gobar"
    profile="profile"
)

func main() {
    evaluator := evaluate.New(name, root, profile)
    cov, err := evaluator.EvalCoverage()
    if err != nil {
        panic(err)
    }
}

```

__Rules__

1. Function directives must be be function documentation, implying that there cannot be whitespace separating the function and comment like this:

```go
//gobar:...

func MyFunc() {...} // This doesn't work!
```

2. The same principle applies to file level directives, which must not have whitespace separating the directive and the package declaration.

3. Other comments __can__ sit in between the directive and the package or function declaration, like this:

```go
// Some documentation comment.
//gobar:min:0.25
// Some other documentation comment.
func MyFunc() {...} // This will work!
```

4. There __can__ be whitespace separating the directive and the slashes, like this:

```go
//      gobar:min:0.25 this is okay!
func MyFunc() {...}
```

5. There __can__ be other comments following the directive, but they cannot come before:

```go
//gobar:min:0.5 this is a acceptable.
// This is just an ordinary comment -- gobar:min:0.5
```

6. Multiple directives can be applied to the same file or block of code, but both will be evaluated, and they will be evaluated in the order in which they appear:

```go
//gobar:min:0.5 -- first we will check for at least 50% unit test coverage.
//gobar:exclude -- then we will exclude this from the overall unit test coverage bar.
func MyFunc() {...}
```

7. Excluding a function from coverage will make it's coverage look like 0 statements and 0 covered statements, implying coverage of 0%.

```go
//gobar:exclude -- the block now has a coverage of 0/0.
//gobar:min:0.5 -- as a result, this is impossible to satisfy.
func MyFunc() {...}
```

8. Directives will apply to both variable declarations if more than one is defined:

```go
//gobar:exclude
var funcOne, funcTwo = func() {}, func() {} // Both will have the same exclusion directive.
```

## Examples

### Multiple directives on one file.

```go
//gobar:pkg:min:0.55
//gobar:pkg:default:file:min:0.5
//gobar:min:0.65
//gobar:default:block:min:0.35
//gobar:pkg:default:block:min:0.30
//gobar:pkg:exclude
package mypack
```

In this case, the following conditions would be true:
1. The package is expected to have at least 40% unit test coverage.
2. All files in the package must have at least 50% unit test coverage.
3. This particular file must have at least 65% unit test coverage.
4. All function blocks in this file must have at least 35% unit test coverage.
5. All function blocks in the whole package must have at least 30% unit test coverage.
6. The package itself will be excluded from the overall unit test coverage calculation.

### Multiple directives on var declaration with special var declarations.

```go
//gobar:min:0.5
var (
    //gobar:min:0.4
    var f1, f2 = func(){...}, func(){...}
    //gobar:min:0.3
    var f3 = func() {...}
    var f4 = func() {...}
)
```

The following statements would be true in this circumstance:
1. Function f1 and function f2 are expected to have 40% unit test coverage.
2. Function f3 is expected to have 30% unit test coverage.
3. Function f4 is expected to have 50% unit test coverage.

### File default declaration on one function.

```go
// gobar:default:file:min:0.5
func MyFunc() {...}
```

This will not impact the coverage bar in any way, because the child function block directive cannot impact its parent file.