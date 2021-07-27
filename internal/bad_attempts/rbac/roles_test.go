package roles

import (
	"context"
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/topdown"

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

	tr := topdown.NewBufferTracer()
	res, err := query.Eval(ctx,
		rego.EvalInput(map[string]interface{}{
			"user":   "steven",
			"op":     "read",
			"object": "workspace",
			"owner":  "steven",
		}),
		rego.EvalQueryTracer(tr),
	)
	require.NoError(t, err)
	fmt.Println(res)
	//fmt.Println(tr.Enabled())
	//for _, e := range *tr {
	//	fmt.Println("\t", e)
	//}
	//var buf bytes.Buffer
	//topdown.PrettyTrace(&buf, *tr)
	//fmt.Println(buf.String())
}
