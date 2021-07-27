package opamanager_test

import (
	"context"
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
	require.NoError(t, err)
}
