package postgres

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	corev3 "github.com/sensu/sensu-go/api/core/v3"
	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/backend/store/patch"
	storev2 "github.com/sensu/sensu-go/backend/store/v2"
)

func testWithPostgresStoreV2(t *testing.T, fn func(storev2.Interface)) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping postgres test")
		return
	}
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		t.Skip("skipping postgres test")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := pgxpool.Connect(ctx, pgURL)
	if err != nil {
		t.Fatal(err)
	}
	dbName := "sensu" + strings.ReplaceAll(uuid.New().String(), "-", "")
	if _, err := db.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName)); err != nil {
		t.Fatal(err)
	}
	defer dropAll(context.Background(), dbName, pgURL)
	db.Close()
	db, err = pgxpool.Connect(ctx, fmt.Sprintf("dbname=%s ", dbName)+pgURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := upgrade(ctx, db); err != nil {
		t.Fatal(err)
	}
	fn(NewStoreV2(db, nil))
}

func testCreateNamespace(t *testing.T, s storev2.Interface, name string) {
	t.Helper()
	ctx := context.Background()
	namespace := corev3.FixtureNamespace(name)
	req := storev2.NewResourceRequestFromResource(ctx, namespace)
	req.UsePostgres = true
	wrapper := WrapNamespace(namespace)
	if err := s.CreateOrUpdate(req, wrapper); err != nil {
		t.Error(err)
	}
}

func testCreateEntityConfig(t *testing.T, s storev2.Interface, name string) {
	t.Helper()
	ctx := context.Background()
	cfg := corev3.FixtureEntityConfig(name)
	req := storev2.NewResourceRequestFromResource(ctx, cfg)
	req.UsePostgres = true
	wrapper := WrapEntityConfig(cfg)
	if err := s.CreateOrUpdate(req, wrapper); err != nil {
		t.Error(err)
	}
}

func testCreateEntityState(t *testing.T, s storev2.Interface, name string) {
	t.Helper()
	ctx := context.Background()
	state := corev3.FixtureEntityState(name)
	req := storev2.NewResourceRequestFromResource(ctx, state)
	req.UsePostgres = true
	wrapper := WrapEntityState(state)
	if err := s.CreateOrUpdate(req, wrapper); err != nil {
		t.Error(err)
	}
}

func TestStoreCreateOrUpdate(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name        string
		args        args
		verifyQuery string
		beforeHook  func(*testing.T, storev2.Interface)
		want        int
	}{
		{
			name: "entity configs can be created and updated",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			verifyQuery: fmt.Sprintf("SELECT * FROM %s", entityConfigStoreName),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
			want: 1,
		},
		{
			name: "entity states can be created and updated",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			verifyQuery: fmt.Sprintf("SELECT * FROM %s", entityStateStoreName),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
				testCreateEntityConfig(t, s, "foo")
			},
			want: 1,
		},
		{
			name: "namespaces can be created and updated",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			verifyQuery: fmt.Sprintf("SELECT * FROM %s", namespaceStoreName),
			want:        1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				if err := s.CreateOrUpdate(tt.args.req, tt.args.wrapper); err != nil {
					t.Error(err)
				}

				// Repeating the call to the store should succeed
				if err := s.CreateOrUpdate(tt.args.req, tt.args.wrapper); err != nil {
					t.Error(err)
				}
				rows, err := s.(*StoreV2).db.Query(context.Background(), tt.verifyQuery)
				if err != nil {
					t.Fatal(err)
				}
				defer rows.Close()
				got := 0
				for rows.Next() {
					got++
				}
				if got != tt.want {
					t.Errorf("bad row count: got %d, want %d", got, tt.want)
				}
			})
		})
	}
}

func TestStoreUpdateIfExists(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "entity configs can be updated if one exists",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "entity states can be updated if one exists",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "namespaces can be updated if one exists",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				cfg := corev3.FixtureEntityConfig("foo")
				ctx := context.Background()
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				wrapper := WrapEntityConfig(cfg)

				// UpdateIfExists should fail
				if err := s.UpdateIfExists(req, wrapper); err == nil {
					t.Error("expected non-nil error")
				} else {
					if _, ok := err.(*store.ErrNotFound); !ok {
						t.Errorf("wrong error: %s", err)
					}
				}
				if err := s.CreateOrUpdate(req, wrapper); err != nil {
					t.Fatal(err)
				}

				// UpdateIfExists should succeed
				if err := s.UpdateIfExists(req, wrapper); err != nil {
					t.Error(err)
				}
			})
		})
	}
}

func TestStoreCreateIfNotExists(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "entity configs can be created if one does not exist",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "entity states can be created if one does not exist",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "namespaces can be created if one does not exist",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				cfg := corev3.FixtureEntityConfig("foo")
				ctx := context.Background()
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				wrapper := WrapEntityConfig(cfg)

				// CreateIfNotExists should succeed
				if err := s.CreateIfNotExists(req, wrapper); err != nil {
					t.Fatal(err)
				}

				// CreateIfNotExists should fail
				if err := s.CreateIfNotExists(req, wrapper); err == nil {
					t.Error("expected non-nil error")
				} else if _, ok := err.(*store.ErrAlreadyExists); !ok {
					t.Errorf("wrong error: %s", err)
				}

				// UpdateIfExists should succeed
				if err := s.UpdateIfExists(req, wrapper); err != nil {
					t.Error(err)
				}
			})
		})
	}
}

func TestStoreGet(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "an entity config can be retrieved",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "an entity state can be retrieved",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "a namespace can be retrieved",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				cfg := corev3.FixtureEntityConfig("foo")
				ctx := context.Background()
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				wrapper := WrapEntityConfig(cfg)

				// CreateIfNotExists should succeed
				if err := s.CreateOrUpdate(req, wrapper); err != nil {
					t.Fatal(err)
				}
				got, err := s.Get(req)
				if err != nil {
					t.Fatal(err)
				}
				if want := wrapper; !reflect.DeepEqual(got, wrapper) {
					t.Errorf("bad resource; got %#v, want %#v", got, want)
				}
			})
		})
	}
}

func TestStoreDelete(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "an entity config can be deleted",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "an entity state can be deleted",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "a namespace can be deleted",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				cfg := corev3.FixtureEntityConfig("foo")
				ctx := context.Background()
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				wrapper := WrapEntityConfig(cfg)
				// CreateIfNotExists should succeed
				if err := s.CreateIfNotExists(req, wrapper); err != nil {
					t.Fatal(err)
				}
				if err := s.Delete(req); err != nil {
					t.Fatal(err)
				}
				if err := s.Delete(req); err == nil {
					t.Error("expected non-nil error")
				} else if _, ok := err.(*store.ErrNotFound); !ok {
					t.Errorf("expected ErrNotFound: got %s", err)
				}
				if _, err := s.Get(req); err == nil {
					t.Error("expected non-nil error")
				} else if _, ok := err.(*store.ErrNotFound); !ok {
					t.Errorf("expected ErrNotFound: got %s", err)
				}
			})
		})
	}
}

func TestStoreList(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "entity configs can be listed",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "entity states can be listed",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "namespaces can be listed",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				for i := 0; i < 10; i++ {
					// create 10 resources
					cfg := corev3.FixtureEntityConfig(fmt.Sprintf("foo-%d", i))
					ctx := context.Background()
					req := storev2.NewResourceRequestFromResource(ctx, cfg)
					req.UsePostgres = true
					wrapper := WrapEntityConfig(cfg)
					if err := s.CreateIfNotExists(req, wrapper); err != nil {
						t.Fatal(err)
					}
				}
				ctx := context.Background()
				req := storev2.NewResourceRequest(ctx, "default", "anything", new(corev3.EntityConfig).StoreName())
				req.UsePostgres = true
				pred := &store.SelectionPredicate{Limit: 5}
				// Test listing with limit of 5
				list, err := s.List(req, pred)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := list.Len(), 5; got != want {
					t.Errorf("wrong number of items: got %d, want %d", got, want)
				}
				if got, want := pred.Continue, `{"offset":5}`; got != want {
					t.Errorf("bad continue token: got %q, want %q", got, want)
				}
				// get the rest of the list
				pred.Limit = 6
				list, err = s.List(req, pred)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := list.Len(), 5; got != want {
					t.Errorf("wrong number of items: got %d, want %d", got, want)
				}
				if pred.Continue != "" {
					t.Error("expected empty continue token")
				}
				// Test listing from all namespaces
				req.Namespace = ""
				pred = &store.SelectionPredicate{Limit: 5}
				list, err = s.List(req, pred)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := list.Len(), 5; got != want {
					t.Errorf("wrong number of items: got %d, want %d", got, want)
				}
				if got, want := pred.Continue, `{"offset":5}`; got != want {
					t.Errorf("bad continue token: got %q, want %q", got, want)
				}
				pred.Limit = 6
				// get the rest of the list
				list, err = s.List(req, pred)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := list.Len(), 5; got != want {
					t.Errorf("wrong number of items: got %d, want %d", got, want)
				}
				if pred.Continue != "" {
					t.Error("expected empty continue token")
				}
				pred.Limit = 5
				// Test listing in descending order
				pred.Continue = ""
				req.SortOrder = storev2.SortDescend
				list, err = s.List(req, pred)
				if err != nil {
					t.Fatal(err)
				}
				if got := list.Len(); got == 0 {
					t.Fatalf("wrong number of items: got %d, want > %d", got, 0)
				}
				firstObj, err := list.(WrapList)[0].Unwrap()
				if err != nil {
					t.Fatal(err)
				}
				if got, want := firstObj.GetMetadata().Name, "foo-9"; got != want {
					t.Errorf("unexpected first item in list: got %s, want %s", got, want)
				}
				// Test listing in ascending order
				pred.Continue = ""
				req.SortOrder = storev2.SortAscend
				list, err = s.List(req, pred)
				if err != nil {
					t.Fatal(err)
				}
				if got := list.Len(); got == 0 {
					t.Fatalf("wrong number of items: got %d, want > %d", got, 0)
				}
				firstObj, err = list.(WrapList)[0].Unwrap()
				if err != nil {
					t.Fatal(err)
				}
				if got, want := firstObj.GetMetadata().Name, "foo-0"; got != want {
					t.Errorf("unexpected first item in list: got %s, want %s", got, want)
				}
			})
		})
	}
}

func TestStoreExists(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "can check if an entity config exists",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "can check if an entity state exists",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "can check if a namespace exists",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				cfg := corev3.FixtureEntityConfig("foo")
				ctx := context.Background()
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				// Exists should return false
				got, err := s.Exists(req)
				if err != nil {
					t.Fatal(err)
				}
				if want := false; got != want {
					t.Errorf("got true, want false")
				}

				// Create a resource under the default namespace
				wrapper := WrapEntityConfig(cfg)
				// CreateIfNotExists should succeed
				if err := s.CreateIfNotExists(req, wrapper); err != nil {
					t.Fatal(err)
				}
				got, err = s.Exists(req)
				if err != nil {
					t.Fatal(err)
				}
				if want := true; got != want {
					t.Errorf("got false, want true")
				}
			})
		})
	}
}

func TestStorePatch(t *testing.T) {
	type args struct {
		req        storev2.ResourceRequest
		wrapper    storev2.Wrapper
		patcher    patch.Patcher
		conditions *store.ETagCondition
	}
	tests := []struct {
		name       string
		args       args
		beforeHook func(*testing.T, storev2.Interface)
	}{
		{
			name: "an entity config can be patched",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
					patcher: &patch.Merge{
						MergePatch: []byte(`{"metadata":{"labels":{"food":"hummus"}}}`),
					},
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
		{
			name: "an entity state can be patched",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
					patcher: &patch.Merge{
						MergePatch: []byte(`{"metadata":{"labels":{"food":"hummus"}}}`),
					},
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
				testCreateEntityConfig(t, s, "foo")
			},
		},
		{
			name: "a namespace can be patched",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
					patcher: &patch.Merge{
						MergePatch: []byte(`{"metadata":{"labels":{"food":"hummus"}}}`),
					},
				}
			}(),
			beforeHook: func(t *testing.T, s storev2.Interface) {
				testCreateNamespace(t, s, "default")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				if tt.beforeHook != nil {
					tt.beforeHook(t, s)
				}
				if err := s.CreateOrUpdate(tt.args.req, tt.args.wrapper); err != nil {
					t.Error(err)
				}
				if err := s.Patch(tt.args.req, tt.args.wrapper, tt.args.patcher, nil); err != nil {
					t.Fatal(err)
				}

				updatedWrapper, err := s.Get(tt.args.req)
				if err != nil {
					t.Fatal(err)
				}

				updated, err := updatedWrapper.Unwrap()
				if err != nil {
					t.Fatal(err)
				}

				if got, want := updated.GetMetadata().Labels["food"], "hummus"; got != want {
					t.Errorf("bad patched labels: got %q, want %q", got, want)
				}
			})
		})
	}
}

func TestStoreGetMultiple(t *testing.T) {
	type args struct {
		req     storev2.ResourceRequest
		wrapper storev2.Wrapper
	}
	tests := []struct {
		name string
		args args
		reqs func(*testing.T, storev2.Interface) []storev2.ResourceRequest
		test func(*testing.T, storev2.Wrapper)
	}{
		{
			name: "multiple entity configs can be retrieved",
			args: func() args {
				ctx := context.Background()
				cfg := corev3.FixtureEntityConfig("foo")
				req := storev2.NewResourceRequestFromResource(ctx, cfg)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityConfig(cfg),
				}
			}(),
			reqs: func(t *testing.T, s storev2.Interface) []storev2.ResourceRequest {
				testCreateNamespace(t, s, "default")
				ctx := context.Background()
				reqs := make([]storev2.ResourceRequest, 0)
				for i := 0; i < 10; i++ {
					entityName := fmt.Sprintf("foo-%d", i)
					cfg := corev3.FixtureEntityConfig(entityName)
					req := storev2.NewResourceRequestFromResource(ctx, cfg)
					req.UsePostgres = true
					reqs = append(reqs, req)
					testCreateEntityConfig(t, s, entityName)
				}
				return reqs
			},
			test: func(t *testing.T, wrapper storev2.Wrapper) {
				var cfg corev3.EntityConfig
				if err := wrapper.UnwrapInto(&cfg); err != nil {
					t.Error(err)
				}

				if got, want := len(cfg.Subscriptions), 2; got != want {
					t.Errorf("wrong number of subscriptions, got = %v, want %v", got, want)
				}
				if got, want := len(cfg.KeepaliveHandlers), 1; got != want {
					t.Errorf("wrong number of keepalive handlers, got = %v, want %v", got, want)
				}
			},
		},
		{
			name: "multiple entity states can be retrieved",
			args: func() args {
				ctx := context.Background()
				state := corev3.FixtureEntityState("foo")
				req := storev2.NewResourceRequestFromResource(ctx, state)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapEntityState(state),
				}
			}(),
			reqs: func(t *testing.T, s storev2.Interface) []storev2.ResourceRequest {
				testCreateNamespace(t, s, "default")
				ctx := context.Background()
				reqs := make([]storev2.ResourceRequest, 0)
				for i := 0; i < 10; i++ {
					entityName := fmt.Sprintf("foo-%d", i)
					state := corev3.FixtureEntityState(entityName)
					req := storev2.NewResourceRequestFromResource(ctx, state)
					req.UsePostgres = true
					reqs = append(reqs, req)
					testCreateEntityConfig(t, s, entityName)
					testCreateEntityState(t, s, entityName)
				}
				return reqs
			},
			test: func(t *testing.T, wrapper storev2.Wrapper) {
				var state corev3.EntityState
				if err := wrapper.UnwrapInto(&state); err != nil {
					t.Error(err)
				}

				if got, want := state.LastSeen, int64(12345); got != want {
					t.Errorf("wrong last_seen value, got = %v, want = %v", got, want)
				}
				if got, want := state.System.Arch, "amd64"; got != want {
					t.Errorf("wrong system arch value, got = %v, want %v", got, want)
				}
			},
		},
		{
			name: "multiple namespaces can be retrieved",
			args: func() args {
				ctx := context.Background()
				ns := corev3.FixtureNamespace("bar")
				req := storev2.NewResourceRequestFromResource(ctx, ns)
				req.UsePostgres = true
				return args{
					req:     req,
					wrapper: WrapNamespace(ns),
				}
			}(),
			reqs: func(t *testing.T, s storev2.Interface) []storev2.ResourceRequest {
				ctx := context.Background()
				reqs := make([]storev2.ResourceRequest, 0)
				for i := 0; i < 10; i++ {
					namespaceName := fmt.Sprintf("foo-%d", i)
					namespace := corev3.FixtureNamespace(namespaceName)
					req := storev2.NewResourceRequestFromResource(ctx, namespace)
					req.UsePostgres = true
					reqs = append(reqs, req)
					testCreateNamespace(t, s, namespaceName)
				}
				return reqs
			},
			test: func(t *testing.T, wrapper storev2.Wrapper) {
				var namespace corev3.Namespace
				if err := wrapper.UnwrapInto(&namespace); err != nil {
					t.Error(err)
				}

				if got, want := len(namespace.Metadata.Labels), 0; got != want {
					t.Errorf("wrong number of labels, got = %v, want = %v", got, want)
				}
				if got, want := len(namespace.Metadata.Annotations), 0; got != want {
					t.Errorf("wrong number of annotations, got = %v, want = %v", got, want)
				}
				if got, want := namespace.Metadata.Name, "foo-"; !strings.Contains(got, want) {
					t.Errorf("wrong namespace name, got = %v, want name to contain = %v", got, want)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWithPostgresStoreV2(t, func(s storev2.Interface) {
				reqs := tt.reqs(t, s)
				result, err := s.(*StoreV2).GetMultiple(context.Background(), reqs[:5])
				if err != nil {
					t.Fatal(err)
				}
				if got, want := len(result), 5; got != want {
					t.Fatalf("bad number of results: got %d, want %d", got, want)
				}
				for i := 0; i < 5; i++ {
					wrapper, ok := result[reqs[i]]
					if !ok {
						t.Errorf("missing result %d", i)
						continue
					}
					tt.test(t, wrapper)
				}
				req := reqs[0]
				req.Name = "notexists"
				result, err = s.(*StoreV2).GetMultiple(context.Background(), []storev2.ResourceRequest{req})
				if err != nil {
					t.Fatal(err)
				}
				if got, want := len(result), 0; got != want {
					t.Fatalf("wrong result length: got %d, want %d", got, want)
				}
			})
		})
	}
}
