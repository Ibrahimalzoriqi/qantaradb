package planner

import (
	"testing"

	"github.com/OWNER/qantaradb/inspector"
	"github.com/OWNER/qantaradb/mapper"
)

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
