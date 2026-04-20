package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateLeaveSupportsPayloadSources(t *testing.T) {
	t.Parallel()

	payload := sampleLeavePayload()

	t.Run("payload flag", func(t *testing.T) {
		t.Parallel()

		result := runCommand(t, nil, "create-leave", "--payload", mustJSONArg(t, payload))
		assertActionSuccess(t, result, "leave", "create", "submitted", "LV-")
	})

	t.Run("payload file approved", func(t *testing.T) {
		t.Parallel()

		payloadPath := writeJSONFile(t, payload)
		result := runCommand(t, nil, "create-leave", "--payload-file", payloadPath, "--result", "approved")
		assertActionSuccess(t, result, "leave", "create", "approved", "LV-")
	})

	t.Run("payload stdin rejected", func(t *testing.T) {
		t.Parallel()

		stdin := bytes.NewBufferString(mustJSONArg(t, payload))
		result := runCommand(t, stdin, "create-leave", "--payload-stdin", "--result", "rejected")
		assertActionSuccess(t, result, "leave", "create", "rejected", "LV-")
	})
}

func TestUpdateCommandsSupportNonPayloadSources(t *testing.T) {
	t.Parallel()

	t.Run("update expense with payload file", func(t *testing.T) {
		t.Parallel()

		payloadPath := writeJSONFile(t, sampleExpensePayload())
		result := runCommand(t, nil, "update-expense", "--request-id", "EX-14C0A7B992", "--payload-file", payloadPath, "--result", "approved")
		assertActionSuccess(t, result, "expense", "update", "approved", "EX-")
	})

	t.Run("update procurement with payload stdin", func(t *testing.T) {
		t.Parallel()

		stdin := bytes.NewBufferString(mustJSONArg(t, sampleProcurementPayload()))
		result := runCommand(t, stdin, "update-procurement", "--request-id", "PR-BA08D42C31", "--payload-stdin", "--result", "rejected")
		assertActionSuccess(t, result, "procurement", "update", "rejected", "PR-")
	})
}

func TestPayloadInputValidation(t *testing.T) {
	t.Parallel()

	payloadPath := writeJSONFile(t, sampleLeavePayload())
	result := runCommand(t, nil, "create-leave", "--payload", mustJSONArg(t, sampleLeavePayload()), "--payload-file", payloadPath)

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "only one of --payload, --payload-file, or --payload-stdin may be used") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestCreateExpenseRequiresPayloadSource(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "create-expense")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "exactly one of --payload, --payload-file, or --payload-stdin is required") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestGetDeleteCommandsRequireRequestID(t *testing.T) {
	t.Parallel()

	getResult := runCommand(t, nil, "get-leave")
	if getResult.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, getResult.code)
	}
	if !strings.Contains(getResult.stderr, "--request-id is required") {
		t.Fatalf("unexpected stderr: %q", getResult.stderr)
	}

	deleteResult := runCommand(t, nil, "delete-expense")
	if deleteResult.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, deleteResult.code)
	}
	if !strings.Contains(deleteResult.stderr, "--request-id is required") {
		t.Fatalf("unexpected stderr: %q", deleteResult.stderr)
	}
}

func TestUpdateLeaveRequiresPayloadRequestID(t *testing.T) {
	t.Parallel()

	payload := sampleLeavePayload()
	result := runCommand(t, nil, "update-leave", "--payload", mustJSONArg(t, payload))

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "payload.request_id is required") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestRequestIDPrefixValidation(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "get-procurement", "--request-id", "LV-7B0A3D4F10")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "--request-id must start with \"PR-\"") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestInvalidResultReturnsUsageError(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "create-leave", "--payload", mustJSONArg(t, sampleLeavePayload()), "--result", "oops")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "invalid --result") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestResultBranchesForGetAndDelete(t *testing.T) {
	t.Parallel()

	getResult := runCommand(t, nil, "get-expense", "--request-id", "EX-14C0A7B992", "--result", "not_found")
	assertActionSuccess(t, getResult, "expense", "get", "not_found", "EX-")
	assertNoRecord(t, getResult.stdout)

	deleteResult := runCommand(t, nil, "delete-procurement", "--request-id", "PR-BA08D42C31", "--result", "not_found")
	assertActionSuccess(t, deleteResult, "procurement", "delete", "not_found", "PR-")
}

func TestBusinessRuleFailuresRemain(t *testing.T) {
	t.Parallel()

	t.Run("leave invalid date range", func(t *testing.T) {
		t.Parallel()

		payload := sampleLeavePayload()
		payload.StartDate = "2026-04-22"
		payload.EndDate = "2026-04-20"

		result := runCommand(t, nil, "create-leave", "--payload", mustJSONArg(t, payload))
		assertBusinessFailureContains(t, result, "start_date must be on or before end_date")
	})

	t.Run("leave supports half day", func(t *testing.T) {
		t.Parallel()

		payload := sampleLeavePayload()
		payload.Days = 0.5

		result := runCommand(t, nil, "create-leave", "--payload", mustJSONArg(t, payload))
		assertActionSuccess(t, result, "leave", "create", "submitted", "LV-")
	})

	t.Run("leave rejects non-positive days", func(t *testing.T) {
		t.Parallel()

		payload := sampleLeavePayload()
		payload.Days = 0

		result := runCommand(t, nil, "create-leave", "--payload", mustJSONArg(t, payload))
		assertBusinessFailureContains(t, result, "leave days must be greater than 0")
	})

	t.Run("expense total mismatch", func(t *testing.T) {
		t.Parallel()

		payload := sampleExpensePayload()
		payload.TotalAmount = 900

		result := runCommand(t, nil, "create-expense", "--payload", mustJSONArg(t, payload))
		assertBusinessFailureContains(t, result, "does not match item amount sum")
	})

	t.Run("expense rejects legacy flat payload fields", func(t *testing.T) {
		t.Parallel()

		result := runCommand(t, nil, "create-expense", "--payload", `{"employee_id":"E1001","department":"engineering","expense_type":"travel","currency":"CNY","total_amount":1280.5,"items":[{"category":"transport","amount":800,"invoice_id":"INV-001","occurred_on":"2026-04-10","description":"flight"}],"submitted_at":"2026-04-14T10:30:00+08:00"}`)
		if result.code != ExitUsage {
			t.Fatalf("expected exit %d, got %d stderr=%q", ExitUsage, result.code, result.stderr)
		}
		if !strings.Contains(result.stderr, `unknown field "employee_id"`) {
			t.Fatalf("expected legacy field rejection, got %q", result.stderr)
		}
	})

	t.Run("procurement budget exceeded", func(t *testing.T) {
		t.Parallel()

		payload := sampleProcurementPayload()
		payload.Items[0].Quantity = 3
		payload.Items[0].UnitPrice = 20000

		result := runCommand(t, nil, "create-procurement", "--payload", mustJSONArg(t, payload))
		assertBusinessFailureContains(t, result, "budget exceeded")
	})
}

func TestLeaveRejectsUnknownLegacyField(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "create-leave", "--payload", `{"applicant_id":"E1001","department_id":"engineering","leave_type":"annual","start_date":"2026-04-20","end_date":"2026-04-22","days":1,"reason":"family_trip","employee_name":"Lin"}`)
	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d stderr=%q", ExitUsage, result.code, result.stderr)
	}
	if !strings.Contains(result.stderr, `unknown field "employee_name"`) {
		t.Fatalf("expected unknown employee_name rejection, got %q", result.stderr)
	}
}

func TestCreateLeaveRequestIDIsStable(t *testing.T) {
	t.Parallel()

	payloadArg := mustJSONArg(t, sampleLeavePayload())

	first := runCommand(t, nil, "create-leave", "--payload", payloadArg)
	second := runCommand(t, nil, "create-leave", "--payload", payloadArg)

	if first.code != ExitSuccess || second.code != ExitSuccess {
		t.Fatalf("expected both commands to succeed, first=%d second=%d", first.code, second.code)
	}
	if first.stdout != second.stdout {
		t.Fatalf("expected stable output, first=%q second=%q", first.stdout, second.stdout)
	}
}

func sampleLeavePayload() leavePayload {
	return leavePayload{
		ApplicantID:  "E1001",
		DepartmentID: "engineering",
		LeaveType:    "annual",
		StartDate:    "2026-04-20",
		EndDate:      "2026-04-22",
		Days:         3,
		Reason:       "family_trip",
	}
}

func sampleExpensePayload() expensePayload {
	return expensePayload{
		Employee: expenseEmployee{
			ID:   "E1001",
			Name: "张三",
		},
		Department: expenseDepartment{
			Code: "engineering",
			Name: "工程部",
		},
		ExpenseType: "travel",
		Currency:    "CNY",
		TotalAmount: 1280.5,
		Items: []expenseLineItem{
			{
				Category:    "transport",
				Amount:      800,
				InvoiceID:   "INV-001",
				OccurredOn:  "2026-04-10",
				Description: "flight",
			},
			{
				Category:    "hotel",
				Amount:      480.5,
				InvoiceID:   "INV-002",
				OccurredOn:  "2026-04-11",
				Description: "hotel",
			},
		},
		SubmittedAt: "2026-04-14T10:30:00+08:00",
	}
}

func sampleProcurementPayload() procurementPayload {
	return procurementPayload{
		RequesterID:  "E1001",
		Department:   "engineering",
		BudgetCode:   "RD-2026-001",
		Reason:       "team expansion",
		DeliveryCity: "Shanghai",
		Items: []procurementLineItem{
			{
				Name:      "MacBook Pro",
				Quantity:  2,
				UnitPrice: 18999,
				Vendor:    "Apple",
			},
		},
		Approvers:   []string{"MGR100", "FIN200"},
		RequestedAt: "2026-04-14T11:00:00+08:00",
	}
}

func writeJSONFile(t *testing.T, value any) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "payload.json")
	if err := os.WriteFile(path, []byte(mustJSONArg(t, value)), 0o600); err != nil {
		t.Fatalf("write payload file: %v", err)
	}
	return path
}

func assertActionSuccess(t *testing.T, result commandResult, wantType, wantAction, wantStatus, wantPrefix string) {
	t.Helper()

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, result.code, result.stderr)
	}
	if result.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.stderr)
	}

	var response actionResponse
	if err := json.Unmarshal([]byte(result.stdout), &response); err != nil {
		t.Fatalf("unmarshal action response: %v", err)
	}
	if response.Type != wantType || response.Action != wantAction || response.Status != wantStatus {
		t.Fatalf("unexpected response: %#v", response)
	}
	if !strings.HasPrefix(response.RequestID, wantPrefix) {
		t.Fatalf("expected request id prefix %q, got %q", wantPrefix, response.RequestID)
	}
}

func assertNoRecord(t *testing.T, stdout string) {
	t.Helper()

	var response actionResponse
	if err := json.Unmarshal([]byte(stdout), &response); err != nil {
		t.Fatalf("unmarshal action response: %v", err)
	}
	if response.Record != nil {
		t.Fatalf("expected no record, got %#v", response.Record)
	}
}

func assertBusinessFailureContains(t *testing.T, result commandResult, want string) {
	t.Helper()

	if result.code != ExitFailure {
		t.Fatalf("expected exit %d, got %d", ExitFailure, result.code)
	}
	if !strings.Contains(result.stderr, want) {
		t.Fatalf("expected stderr to contain %q, got %q", want, result.stderr)
	}
}
