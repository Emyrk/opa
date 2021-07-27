package opamanager

import (
	"context"

	"github.com/open-policy-agent/opa/rego"
	"golang.org/x/xerrors"
)

type opaQuery string

const (
	QueryResourceWorkspace opaQuery = "workspace"
)

func (q opaQuery) Query(ctx context.Context, m *OPAManager, input interface{}) (rego.ResultSet, error) {
	p := m.getPrep(q)
	return p.Eval(ctx, rego.EvalInput(input))
}

var Q q

type q struct{}

type AccessResourceInput struct {
	Actor  ActorInput  `json:"actor"`
	Object ObjectInput `json:"object"`
}

type ActorInput struct {
	Op   string `json:"op"`
	User string `json:"user"`
	// Do not populate roles when making the input...
	Roles []string `json:"roles"`
}

type ObjectInput struct {
	ID     string   `json:"id"`
	Owner  string   `json:"owner"`
	Shared []string `json:"shared"`
	Type   string   `json:"type"`
}

func (q) CanAccessWorkspace(ctx context.Context, m *OPAManager, input AccessResourceInput) error {
	input.Actor.Roles = UserRoles[input.Actor.User] // Fix user roles from the lookup table
	res, err := QueryResourceWorkspace.Query(ctx, m, input)
	if err != nil {
		return err
	}

	return checkResultSet(res)
}

func checkResultSet(set rego.ResultSet) error {
	if len(set) == 0 {
		return xerrors.Errorf("undefined decision")
	}

	// Not sure what this is...
	if len(set) > 1 {
		return xerrors.Errorf("got 2 results back for the query")
	}

	result := set[0]
	allow, ok := result.Bindings["allow"]
	if !ok {
		return xerrors.Errorf("no allow variable set")
	}

	aB, ok := allow.(bool)
	if !ok {
		return xerrors.Errorf("allow var is not a boolean: %v", allow)
	}
	if !aB {
		return xerrors.Errorf("rejected, allow is false")
	}

	return nil
}
