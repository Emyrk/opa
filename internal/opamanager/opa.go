package opamanager

import (
	"context"
	"io/fs"

	"golang.org/x/xerrors"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/bundle"
	"github.com/open-policy-agent/opa/metrics"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

// OPAManager will store policies and data in a memory store
type OPAManager struct {
	store    storage.Store
	compiler *ast.Compiler
	metrics  metrics.Metrics

	prepared map[opaQuery]rego.PreparedEvalQuery
}

// LoadOPAManager will load all policies and data into memory
func LoadOPAManager(ctx context.Context, f fs.FS) (*OPAManager, error) {
	m := new(OPAManager)
	m.prepared = make(map[opaQuery]rego.PreparedEvalQuery)

	// Filesystem loader
	fsL, err := bundle.NewFSLoader(f)
	if err != nil {
		return nil, xerrors.Errorf("fs loader: %w", err)
	}

	r := bundle.NewCustomReader(fsL)
	b, err := r.Read()
	if err != nil {
		return nil, xerrors.Errorf("read bundle: %w", err)
	}

	m.store = inmem.New()
	txn, err := m.store.NewTransaction(ctx, storage.TransactionParams{
		Write: true,
	})
	if err != nil {
		return nil, xerrors.Errorf("start txn: %w", err)
	}

	m.compiler = ast.NewCompiler().WithPathConflictsCheck(storage.NonEmpty(ctx, m.store, txn))
	// Not sure what these metrics are for.
	m.metrics = metrics.New()

	m.metrics.Info()
	// Activate the bundle into the store
	err = bundle.Activate(&bundle.ActivateOpts{
		Ctx:      ctx,
		Store:    m.store,
		Txn:      txn,
		Compiler: m.compiler,
		Metrics:  m.metrics,
		Bundles:  map[string]*bundle.Bundle{".": &b},
	})
	if err != nil {
		return nil, xerrors.Errorf("activate bundle: %w", err)
	}

	err = m.store.Commit(ctx, txn)
	if err != nil {
		m.store.Abort(ctx, txn)
		return nil, xerrors.Errorf("commit txn: %w", err)
	}

	err = m.prepareQueries(ctx)
	if err != nil {
		return nil, xerrors.Errorf("prepare queries: %w", err)
	}

	return m, nil
}

func (m *OPAManager) getPrep(q opaQuery) rego.PreparedEvalQuery {
	return m.prepared[q]
}

func (m *OPAManager) prepareQueries(ctx context.Context) error {
	all := []struct {
		Name    opaQuery
		Imports []string
		Query   string
	}{
		{QueryResourceWorkspace, []string{"data.rbac.resources.workspace.allow"},
			"allow := data.rbac.resources.workspace.allow; x := 0",
		},
	}

	for _, one := range all {
		err := m.prepQ(ctx, one.Name, one.Query, one.Imports)
		if err != nil {
			return xerrors.Errorf("query %q prep: %w", one.Name, err)
		}
	}

	return nil
}

func (m *OPAManager) prepQ(ctx context.Context, queryName opaQuery, query string, imports []string, opts ...rego.PrepareOption) error {
	p, err := rego.New(append(m.regoOpts(),
		rego.Imports(imports),
		rego.Query(query),
	)...).PrepareForEval(ctx, opts...)
	if err != nil {
		return xerrors.Errorf("prepare query: %w", err)
	}

	m.prepared[queryName] = p
	return nil
}

// This is not thread safe
func (m *OPAManager) regoOpts() []func(r *rego.Rego) {
	return []func(r *rego.Rego){
		rego.Store(m.store),
		rego.Compiler(m.compiler),
	}
}
