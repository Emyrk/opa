package roles

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/open-policy-agent/opa/metrics"

	"github.com/open-policy-agent/opa/bundle"

	"github.com/open-policy-agent/opa/storage"

	"github.com/open-policy-agent/opa/storage/inmem"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
)

//go:embed *.rego
var policies embed.FS

func Load() (*rego.Rego, error) {
	var (
		ctx = context.Background()
	)

	b, err := bundle.NewFSLoader(policies)
	if err != nil {
		return nil, err
	}
	bunR := bundle.NewCustomReader(b)
	bun, err := bunR.Read()
	if err != nil {
		return nil, err
	}

	store := inmem.New()
	txn, err := store.NewTransaction(ctx, storage.TransactionParams{
		Write: true,
	})
	if err != nil {
		return nil, err
	}

	compiler := ast.NewCompiler()

	err = bundle.Activate(&bundle.ActivateOpts{
		Ctx:      ctx,
		Store:    store,
		Txn:      txn,
		TxnCtx:   nil,
		Compiler: compiler,
		Bundles: map[string]*bundle.Bundle{
			"bundle": &bun,
		},
		ExtraModules: nil,
		Metrics:      metrics.New(),
	})
	if err != nil {
		return nil, err
	}

	q := rego.New(
		rego.Store(store),
		rego.Imports([]string{"data.workspace"}),
	)

	return q, nil
}

func OPA() (*rego.Rego, error) {
	var (
		ctx = context.Background()
	)

	store := inmem.New()
	dir, err := policies.ReadDir(".")
	if err != nil {
		return nil, err
	}

	txn, err := store.NewTransaction(ctx, storage.TransactionParams{
		Write: true,
	})
	if err != nil {
		return nil, err
	}

	for _, file := range dir {
		data, err := policies.ReadFile(file.Name())
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(file.Name())
		fName := strings.TrimSuffix(file.Name(), ext)

		err = store.Write(ctx, txn, storage.AddOp, storage.MustParsePath("/"+fName), data)
		if err != nil {
			return nil, err
		}

		err = store.UpsertPolicy(ctx, txn, fName, data)
		if err != nil {
			return nil, err
		}
	}

	err = store.Write(ctx, txn, storage.AddOp, storage.MustParsePath("/test"), "a")
	if err != nil {
		return nil, err
	}

	pols, _ := store.ListPolicies(ctx, txn)
	fmt.Println(pols)

	err = store.Commit(ctx, txn)
	if err != nil {
		return nil, err
	}

	return rego.New(
		rego.Store(store),
		//rego.Imports([]string{"data.workspace.rego"}),
		rego.Query("data.workspace.rego"),
	), nil
}

//
//
//func Store() {
//	store := inmem.New()
//	dir, err := policies.ReadDir(".")
//	if err != nil {
//		panic(err)
//	}
//
//	var opts []func(r *rego.Rego)
//
//	for _, file := range dir {
//		module, err := policies.ReadFile(file.Name())
//		if err != nil {
//			panic(err)
//		}
//
//		store.Write()
//
//		opts = append(opts, rego.Module(file.Name(), string(module)))
//	}
//
//	opts = append(opts, rego.Query("rbac.workspace.allow"))
//	opts = append(opts, rego.Imports([]string{"data.rbac.workspace"}))
//	r := rego.New(opts...)
//	return r
//
//	store := inmem.New()
//	store.
//		txn, _ := store.NewTransaction(context.Background())
//
//}
//

func compiler() *ast.Compiler {
	dir, err := policies.ReadDir(".")
	if err != nil {
		panic(err)
	}

	modules := make(map[string]string)

	for _, file := range dir {
		module, err := policies.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}

		modules[file.Name()] = string(module)
	}

	compiler, err := ast.CompileModules(modules)
	if err != nil {
		panic(err)
	}
	return compiler
}

func Workspace() *rego.Rego {
	return rego.New(
		rego.Compiler(compiler()),
		rego.Imports([]string{"data.rbac.workspace"}),
		rego.Query("rbac.workspace.allow"),
	)
}

func Policies() *rego.Rego {
	dir, err := policies.ReadDir(".")
	if err != nil {
		panic(err)
	}

	var opts []func(r *rego.Rego)

	for _, file := range dir {
		module, err := policies.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}

		opts = append(opts, rego.Module(file.Name(), string(module)))
	}

	opts = append(opts, rego.Query("data.workspace.allow"))
	opts = append(opts, rego.Imports([]string{"data.workspace"}))
	r := rego.New(opts...)
	return r
}
