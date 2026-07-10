package planner

import (
	"testing"

	"github.com/OWNER/qantaradb/inspector"
	"github.com/OWNER/qantaradb/mapper"
)

func TestSortTablesTopologically(t *testing.T) {
	tables := []inspector.Table{
		{
			Name:        "order_items",
			ForeignKeys: []inspector.ForeignKey{{ReferencedTable: "orders"}},
		},
		{
			Name:        "orders",
			ForeignKeys: []inspector.ForeignKey{{ReferencedTable: "users"}, {ReferencedTable: "restaurants"}},
		},
		{Name: "users", ForeignKeys: []inspector.ForeignKey{}},
		{Name: "restaurants", ForeignKeys: []inspector.ForeignKey{}},
	}

	order, err := sortTablesTopologically(tables, map[string]bool{})
	if err != nil {
		t.Fatalf("unexpected topological sorting error: %v", err)
	}

	positions := make(map[string]int)
	for index, name := range order {
		positions[name] = index
	}

	if len(order) != 4 {
		t.Errorf("expected 4 tables in sorted order, got %d", len(order))
	}
	if positions["users"] >= positions["orders"] {
		t.Errorf("expected users to come before orders")
	}
	if positions["restaurants"] >= positions["orders"] {
		t.Errorf("expected restaurants to come before orders")
	}
	if positions["orders"] >= positions["order_items"] {
		t.Errorf("expected orders to come before order_items")
	}
}

func TestCircularDependencyFallback(t *testing.T) {
	tables := []inspector.Table{
		{Name: "tableA", ForeignKeys: []inspector.ForeignKey{{ReferencedTable: "tableB"}}},
		{Name: "tableB", ForeignKeys: []inspector.ForeignKey{{ReferencedTable: "tableA"}}},
	}

	if _, err := sortTablesTopologically(tables, map[string]bool{}); err == nil {
		t.Error("expected circular dependency error, got nil")
	}
}

func TestCreatePlanAppliesQualifiedTinyintOverride(t *testing.T) {
	schema := &inspector.SchemaInfo{Tables: []inspector.Table{
		{
			Name: "offline_payments",
			Columns: []inspector.Column{
				{Name: "id", DataType: "bigint", ColumnType: "bigint", Extra: "auto_increment"},
				{Name: "status", DataType: "tinyint", ColumnType: "tinyint(1)"},
			},
		},
	}}

	plan, err := CreatePlan(schema, mapper.Config{
		Tinyint1AsBool:          true,
		Tinyint1BooleanDenylist: []string{"offline_payments.status"},
	}, nil)
	if err != nil {
		t.Fatalf("CreatePlan() error = %v", err)
	}

	if len(plan.Tables) != 1 || len(plan.Tables[0].Columns) != 2 {
		t.Fatalf("unexpected plan shape: %#v", plan.Tables)
	}
	if got := plan.Tables[0].Columns[1].TargetType; got != "smallint" {
		t.Fatalf("status target type = %q, want smallint", got)
	}
}
