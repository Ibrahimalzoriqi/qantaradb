package main

import "testing"

func TestValidateLocalExecutionGateAcceptsDedicatedStagingTarget(t *testing.T) {
	err := validateLocalExecutionGate(
		"staging:password@tcp(127.0.0.1:3306)/FoodTechConsoleStaging",
		"postgres://staging:password@127.0.0.1:5432/FoodTechConsoleStaging?sslmode=disable",
		true,
	)

	if err != nil {
		t.Fatalf("dedicated staging target should pass the local execution gate: %v", err)
	}
}

func TestValidateLocalExecutionGateRejectsProductionNamedTarget(t *testing.T) {
	err := validateLocalExecutionGate(
		"staging:password@tcp(127.0.0.1:3306)/FoodTechConsoleStaging",
		"postgres://staging:password@127.0.0.1:5432/FoodTechConsoleProduction?sslmode=disable",
		true,
	)

	if err == nil {
		t.Fatal("production-named target must remain blocked")
	}
}

func TestIsTestDatabaseRecognizesStagingButNotProduction(t *testing.T) {
	if !isTestDatabase("postgres://staging:password@127.0.0.1:5432/FoodTechConsoleStaging?sslmode=disable") {
		t.Fatal("dedicated staging target must be recognized by the destructive-operation safety gate")
	}

	if isTestDatabase("postgres://production:password@127.0.0.1:5432/FoodTechConsoleProduction?sslmode=disable") {
		t.Fatal("production target must remain outside the staging/test allowlist")
	}
}
