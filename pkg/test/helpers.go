package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func NewStruct(t *testing.T, d map[string]interface{}) *structpb.Struct {
	t.Helper()

	strct, err := structpb.NewStruct(d)
	require.NoError(t, err)
	return strct
}
