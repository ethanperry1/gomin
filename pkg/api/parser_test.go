package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	parser := NewProfileParser("../..", "profile", "github.com/ethanperry1/gomin")
	_, err := parser.CreateNodeTree()
	require.NoError(t, err)
}