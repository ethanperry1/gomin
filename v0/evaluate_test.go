package api

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func TestEvaluator(t *testing.T) {
// 	evaluator, err := CreateEvaluator("", "", "")
// 	require.NoError(t, err)

// 	evaluator.Evaluate(
// 		Exclude(
// 			Package("my_pack").Files().Filter("my_file").Functions().Filter("my_func"),
// 			AllFiles().Filter("file"),
// 			AllFunctions().Filter("funky"),
// 			Package("p").File("a").Literal(0),
// 			Package("p").File("a").Method("", ""),
// 		),
// 	)
// }