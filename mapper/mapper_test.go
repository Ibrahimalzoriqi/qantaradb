package mapper

import (
	"testing"

	"github.com/OWNER/qantaradb/inspector"
)

func TestMapTableTypeKeepsQualifiedTinyintStatusNumeric(t *testing.T) {
	column := inspector.Column{Name: "status", DataType: "tinyint", ColumnType: "tinyint(1)"}
	config := Config{
		Tinyint1AsBool:          true,
		Tinyint1BooleanDenylist: []string{"offline_payments.status"},
	}

	targetType, _, err := MapTableType("offline_payments", column, config)
	if err != nil {
		t.Fatalf("MapTableType() error = %v", err)
	}
	if targetType != "smallint" {
		t.Fatalf("target type = %q, want smallint", targetType)
	}
}

func TestMapTableTypePreservesBooleanDefaultOutsideDenylist(t *testing.T) {
	column := inspector.Column{Name: "is_active", DataType: "tinyint", ColumnType: "tinyint(1)"}
	config := Config{
		Tinyint1AsBool:          true,
		Tinyint1BooleanDenylist: []string{"offline_payments.status"},
	}

	targetType, _, err := MapTableType("offline_payments", column, config)
	if err != nil {
		t.Fatalf("MapTableType() error = %v", err)
	}
	if targetType != "boolean" {
		t.Fatalf("target type = %q, want boolean", targetType)
	}
}
