package roles

import (
	"context"
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/require"
)

func TestRoles(t *testing.T) {
	var (
		ctx = context.Background()
	)

	//r, err := Load()
	//require.NoError(t, err)

	r := Policies()

	query, err := r.PrepareForEval(ctx)
	require.NoError(t, err)

	res, err := query.Eval(ctx, rego.EvalInput(map[string]interface{}{
		"user":   "steven",
		"op":     "read",
		"object": "workspace",
		"owner":  "steven",
	}))
	require.NoError(t, err)
	fmt.Println(res)
}
