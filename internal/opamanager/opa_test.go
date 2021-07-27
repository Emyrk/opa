package opamanager_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cdr/opa/internal/opamanager"
	"github.com/stretchr/testify/require"
)

func TestOPAWorkspace(t *testing.T) {

	var (
		ctx = context.Background()
	)

	m, err := opamanager.LoadOPAManager(ctx, opamanager.Policies)
	require.NoError(t, err)

	err = opamanager.Q.CanAccessWorkspace(ctx, m, opamanager.AccessResourceInput{
		Actor: opamanager.ActorInput{
			Op:   "read",
			User: "steven",
		},
		Object: opamanager.ObjectInput{
			ID:     "1234",
			Owner:  "dean",
			Shared: []string{"steven"},
			Type:   "workspace",
		},
	})
	fmt.Println(m.ListPolicies(ctx))
	fmt.Println(m.Read(ctx, "/rbac/roles"))
	require.NoError(t, err)
}
