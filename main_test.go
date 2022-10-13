package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddSuccess(t *testing.T) {
	c := require.New(t)

	result := Add(1, 2)

	expect := 3

	c.Equal(expect, result)
}
