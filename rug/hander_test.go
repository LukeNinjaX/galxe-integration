package rug

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRug(t *testing.T) {
	rug, err := NewRug(nil)
	require.Equal(t, nil, err)

	_ = rug
}
