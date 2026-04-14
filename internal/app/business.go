package app

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const procurementBudgetLimit = 50000.0

var (
	createUpdateResults = []string{"submitted", "approved", "rejected"}
	getResults          = []string{"found", "not_found"}
	deleteResults       = []string{"deleted", "not_found"}
)

type payloadInputOptions struct {
	Payload      string
	PayloadFile  string
	PayloadStdin bool
}

type leavePayload struct {
	RequestID     string `json:"request_id,omitempty"`
	EmployeeID    string `json:"employee_id"`
	EmployeeName  string `json:"employee_name"`
	LeaveType     string `json:"leave_type"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Days          int    `json:"days"`
	Reason        string `json:"reason"`
	HandoverTo    string `json:"handover_to"`
	UrgentContact string `json:"urgent_contact"`
}

type expensePayload struct {
	EmployeeID  string            `json:"employee_id"`
	Department  string            `json:"department"`
	ExpenseType string            `json:"expense_type"`
	Currency    string            `json:"currency"`
	TotalAmount float64           `json:"total_amount"`
	Items       []expenseLineItem `json:"items"`
	SubmittedAt string            `json:"submitted_at"`
}

type expenseLineItem struct {
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	InvoiceID   string  `json:"invoice_id"`
	OccurredOn  string  `json:"occurred_on"`
	Description string  `json:"description"`
}

type procurementPayload struct {
	RequesterID  string                `json:"requester_id"`
	Department   string                `json:"department"`
	BudgetCode   string                `json:"budget_code"`
	Reason       string                `json:"reason"`
	DeliveryCity string                `json:"delivery_city"`
	Items        []procurementLineItem `json:"items"`
	Approvers    []string              `json:"approvers"`
	RequestedAt  string                `json:"requested_at"`
}

type procurementLineItem struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Vendor    string  `json:"vendor"`
}

type actionResponse struct {
	Type       string `json:"type"`
	Action     string `json:"action"`
	Status     string `json:"status"`
	RequestID  string `json:"request_id"`
	Summary    any    `json:"summary,omitempty"`
	Record     any    `json:"record,omitempty"`
	Validation string `json:"validation,omitempty"`
}

type leaveSummary struct {
	EmployeeID   string `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	LeaveType    string `json:"leave_type"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Days         int    `json:"days"`
	Reason       string `json:"reason"`
	HandoverTo   string `json:"handover_to"`
}

type expenseSummary struct {
	EmployeeID     string  `json:"employee_id"`
	Department     string  `json:"department"`
	ExpenseType    string  `json:"expense_type"`
	Currency       string  `json:"currency"`
	TotalAmount    float64 `json:"total_amount"`
	ItemCount      int     `json:"item_count"`
	SubmittedAt    string  `json:"submitted_at"`
	PrimaryInvoice string  `json:"primary_invoice"`
}

type procurementSummary struct {
	RequesterID   string   `json:"requester_id"`
	Department    string   `json:"department"`
	BudgetCode    string   `json:"budget_code"`
	DeliveryCity  string   `json:"delivery_city"`
	TotalAmount   float64  `json:"total_amount"`
	ItemCount     int      `json:"item_count"`
	ApproverCount int      `json:"approver_count"`
	Approvers     []string `json:"approvers"`
	RequestedAt   string   `json:"requested_at"`
}

type leaveRecord struct {
	RequestID     string `json:"request_id"`
	EmployeeID    string `json:"employee_id"`
	EmployeeName  string `json:"employee_name"`
	LeaveType     string `json:"leave_type"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Days          int    `json:"days"`
	Reason        string `json:"reason"`
	HandoverTo    string `json:"handover_to"`
	UrgentContact string `json:"urgent_contact"`
	Status        string `json:"status"`
}

type expenseRecord struct {
	RequestID   string            `json:"request_id"`
	EmployeeID  string            `json:"employee_id"`
	Department  string            `json:"department"`
	ExpenseType string            `json:"expense_type"`
	Currency    string            `json:"currency"`
	TotalAmount float64           `json:"total_amount"`
	Items       []expenseLineItem `json:"items"`
	SubmittedAt string            `json:"submitted_at"`
	Status      string            `json:"status"`
}

type procurementRecord struct {
	RequestID    string                `json:"request_id"`
	RequesterID  string                `json:"requester_id"`
	Department   string                `json:"department"`
	BudgetCode   string                `json:"budget_code"`
	Reason       string                `json:"reason"`
	DeliveryCity string                `json:"delivery_city"`
	Items        []procurementLineItem `json:"items"`
	Approvers    []string              `json:"approvers"`
	RequestedAt  string                `json:"requested_at"`
	Status       string                `json:"status"`
}

func newCreateLeaveCommand() *cobra.Command {
	var input payloadInputOptions
	result := "submitted"

	cmd := &cobra.Command{
		Use:         "create-leave",
		Short:       "Create a mock leave application",
		Description: "Validate and create a mock leave application from JSON input.",
		Example: "mock create-leave --payload '{\"employee_id\":\"E1001\",\"employee_name\":\"Lin\",\"leave_type\":\"annual\",\"start_date\":\"2026-04-20\",\"end_date\":\"2026-04-22\",\"days\":3,\"reason\":\"family_trip\",\"handover_to\":\"E2001\",\"urgent_contact\":\"13800138000\"}'\n" +
			"mock create-leave --payload-file ./leave.json --result approved\n" +
			"printf '{\"employee_id\":\"E1001\",\"employee_name\":\"Lin\",\"leave_type\":\"annual\",\"start_date\":\"2026-04-20\",\"end_date\":\"2026-04-22\",\"days\":3,\"reason\":\"family_trip\",\"handover_to\":\"E2001\",\"urgent_contact\":\"13800138000\"}' | mock create-leave --payload-stdin --result rejected",
		Args: cobra.NoArgs,
		ParamFields: append(
			payloadSourceFields("leave request"),
			resultField("submitted", createUpdateResults),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readPayloadInput(cmd, input)
			if err != nil {
				return err
			}
			payload, err := parseLeavePayload(raw, false)
			if err != nil {
				return err
			}
			status, err := parseResult(result, createUpdateResults)
			if err != nil {
				return err
			}
			requestID, err := stableRequestID("LV-", leavePayloadForStableID(payload))
			if err != nil {
				return failuref("build leave request id: %v", err)
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:       "leave",
				Action:     "create",
				Status:     status,
				RequestID:  requestID,
				Summary:    buildLeaveSummary(payload),
				Validation: "passed",
			})
		},
	}

	addPayloadInputFlags(cmd, &input)
	cmd.Flags().StringVar(&result, "result", result, resultUsage(createUpdateResults))
	return cmd
}

func newGetLeaveCommand() *cobra.Command {
	var requestID string
	result := "found"

	cmd := &cobra.Command{
		Use:         "get-leave",
		Short:       "Get a mock leave application",
		Description: "Return a mock leave application record for the given request id.",
		Example:     "mock get-leave --request-id LV-7B0A3D4F10\nmock get-leave --request-id LV-7B0A3D4F10 --result not_found",
		Args:        cobra.NoArgs,
		ParamFields: []cobra.HelpField{
			requiredField("request-id", "string", "Target leave request id"),
			resultField("found", getResults),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "LV-", "--request-id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, getResults)
			if err != nil {
				return err
			}

			response := actionResponse{
				Type:      "leave",
				Action:    "get",
				Status:    status,
				RequestID: id,
			}
			if status == "found" {
				response.Record = mockLeaveRecord(id)
			}
			return writeJSON(cmd.OutOrStdout(), response)
		},
	}

	cmd.Flags().StringVar(&requestID, "request-id", "", "Target leave request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(getResults))
	return cmd
}

func newUpdateLeaveCommand() *cobra.Command {
	var input payloadInputOptions
	result := "submitted"

	cmd := &cobra.Command{
		Use:         "update-leave",
		Short:       "Update a mock leave application",
		Description: "Validate and update a mock leave application from JSON input.",
		Example: "mock update-leave --payload '{\"request_id\":\"LV-7B0A3D4F10\",\"employee_id\":\"E1001\",\"employee_name\":\"Lin\",\"leave_type\":\"annual\",\"start_date\":\"2026-04-21\",\"end_date\":\"2026-04-23\",\"days\":3,\"reason\":\"family_trip\",\"handover_to\":\"E2001\",\"urgent_contact\":\"13800138000\"}'\n" +
			"mock update-leave --payload-file ./leave-update.json --result approved\n" +
			"printf '{\"request_id\":\"LV-7B0A3D4F10\",\"employee_id\":\"E1001\",\"employee_name\":\"Lin\",\"leave_type\":\"personal\",\"start_date\":\"2026-04-21\",\"end_date\":\"2026-04-22\",\"days\":2,\"reason\":\"family_trip\",\"handover_to\":\"E2001\",\"urgent_contact\":\"13800138000\"}' | mock update-leave --payload-stdin --result rejected",
		Args: cobra.NoArgs,
		ParamFields: append(
			payloadSourceFields("leave update payload including request_id"),
			resultField("submitted", createUpdateResults),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readPayloadInput(cmd, input)
			if err != nil {
				return err
			}
			payload, err := parseLeavePayload(raw, true)
			if err != nil {
				return err
			}
			id, err := validateRequestID(payload.RequestID, "LV-", "payload.request_id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, createUpdateResults)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:       "leave",
				Action:     "update",
				Status:     status,
				RequestID:  id,
				Summary:    buildLeaveSummary(payload),
				Validation: "passed",
			})
		},
	}

	addPayloadInputFlags(cmd, &input)
	cmd.Flags().StringVar(&result, "result", result, resultUsage(createUpdateResults))
	return cmd
}

func newDeleteLeaveCommand() *cobra.Command {
	var requestID string
	result := "deleted"

	cmd := &cobra.Command{
		Use:         "delete-leave",
		Short:       "Delete a mock leave application",
		Description: "Return a mock leave deletion receipt for the given request id.",
		Example:     "mock delete-leave --request-id LV-7B0A3D4F10\nmock delete-leave --request-id LV-7B0A3D4F10 --result not_found",
		Args:        cobra.NoArgs,
		ParamFields: []cobra.HelpField{
			requiredField("request-id", "string", "Target leave request id"),
			resultField("deleted", deleteResults),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "LV-", "--request-id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, deleteResults)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:      "leave",
				Action:    "delete",
				Status:    status,
				RequestID: id,
			})
		},
	}

	cmd.Flags().StringVar(&requestID, "request-id", "", "Target leave request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(deleteResults))
	return cmd
}

func newCreateExpenseCommand() *cobra.Command {
	var input payloadInputOptions
	result := "submitted"

	cmd := &cobra.Command{
		Use:         "create-expense",
		Short:       "Create a mock expense reimbursement",
		Description: "Validate and create a mock expense reimbursement from JSON input.",
		Example: "mock create-expense --payload '{\"employee_id\":\"E1001\",\"department\":\"engineering\",\"expense_type\":\"travel\",\"currency\":\"CNY\",\"total_amount\":1280.5,\"items\":[{\"category\":\"transport\",\"amount\":800,\"invoice_id\":\"INV-001\",\"occurred_on\":\"2026-04-10\",\"description\":\"flight\"},{\"category\":\"hotel\",\"amount\":480.5,\"invoice_id\":\"INV-002\",\"occurred_on\":\"2026-04-11\",\"description\":\"hotel\"}],\"submitted_at\":\"2026-04-14T10:30:00+08:00\"}'\n" +
			"mock create-expense --payload-file ./expense.json --result approved\n" +
			"printf '{\"employee_id\":\"E1001\",\"department\":\"engineering\",\"expense_type\":\"travel\",\"currency\":\"CNY\",\"total_amount\":1280.5,\"items\":[{\"category\":\"transport\",\"amount\":800,\"invoice_id\":\"INV-001\",\"occurred_on\":\"2026-04-10\",\"description\":\"flight\"},{\"category\":\"hotel\",\"amount\":480.5,\"invoice_id\":\"INV-002\",\"occurred_on\":\"2026-04-11\",\"description\":\"hotel\"}],\"submitted_at\":\"2026-04-14T10:30:00+08:00\"}' | mock create-expense --payload-stdin --result rejected",
		Args: cobra.NoArgs,
		ParamFields: append(
			payloadSourceFields("expense request"),
			resultField("submitted", createUpdateResults),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readPayloadInput(cmd, input)
			if err != nil {
				return err
			}
			payload, err := parseExpensePayload(raw)
			if err != nil {
				return err
			}
			status, err := parseResult(result, createUpdateResults)
			if err != nil {
				return err
			}
			requestID, err := stableRequestID("EX-", payload)
			if err != nil {
				return failuref("build expense request id: %v", err)
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:       "expense",
				Action:     "create",
				Status:     status,
				RequestID:  requestID,
				Summary:    buildExpenseSummary(payload),
				Validation: "passed",
			})
		},
	}

	addPayloadInputFlags(cmd, &input)
	cmd.Flags().StringVar(&result, "result", result, resultUsage(createUpdateResults))
	return cmd
}

func newGetExpenseCommand() *cobra.Command {
	var requestID string
	result := "found"

	cmd := &cobra.Command{
		Use:         "get-expense",
		Short:       "Get a mock expense reimbursement",
		Description: "Return a mock expense reimbursement record for the given request id.",
		Example:     "mock get-expense --request-id EX-14C0A7B992\nmock get-expense --request-id EX-14C0A7B992 --result not_found",
		Args:        cobra.NoArgs,
		ParamFields: []cobra.HelpField{
			requiredField("request-id", "string", "Target expense request id"),
			resultField("found", getResults),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "EX-", "--request-id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, getResults)
			if err != nil {
				return err
			}

			response := actionResponse{
				Type:      "expense",
				Action:    "get",
				Status:    status,
				RequestID: id,
			}
			if status == "found" {
				response.Record = mockExpenseRecord(id)
			}
			return writeJSON(cmd.OutOrStdout(), response)
		},
	}

	cmd.Flags().StringVar(&requestID, "request-id", "", "Target expense request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(getResults))
	return cmd
}

func newUpdateExpenseCommand() *cobra.Command {
	var input payloadInputOptions
	var requestID string
	result := "submitted"

	cmd := &cobra.Command{
		Use:         "update-expense",
		Short:       "Update a mock expense reimbursement",
		Description: "Validate and update a mock expense reimbursement from JSON input.",
		Example: "mock update-expense --request-id EX-14C0A7B992 --payload '{\"employee_id\":\"E1001\",\"department\":\"engineering\",\"expense_type\":\"travel\",\"currency\":\"CNY\",\"total_amount\":1280.5,\"items\":[{\"category\":\"transport\",\"amount\":800,\"invoice_id\":\"INV-001\",\"occurred_on\":\"2026-04-10\",\"description\":\"flight\"},{\"category\":\"hotel\",\"amount\":480.5,\"invoice_id\":\"INV-002\",\"occurred_on\":\"2026-04-11\",\"description\":\"hotel\"}],\"submitted_at\":\"2026-04-14T10:30:00+08:00\"}'\n" +
			"mock update-expense --request-id EX-14C0A7B992 --payload-file ./expense-update.json --result approved\n" +
			"printf '{\"employee_id\":\"E1001\",\"department\":\"engineering\",\"expense_type\":\"travel\",\"currency\":\"CNY\",\"total_amount\":1280.5,\"items\":[{\"category\":\"transport\",\"amount\":800,\"invoice_id\":\"INV-001\",\"occurred_on\":\"2026-04-10\",\"description\":\"flight\"},{\"category\":\"hotel\",\"amount\":480.5,\"invoice_id\":\"INV-002\",\"occurred_on\":\"2026-04-11\",\"description\":\"hotel\"}],\"submitted_at\":\"2026-04-14T10:30:00+08:00\"}' | mock update-expense --request-id EX-14C0A7B992 --payload-stdin --result rejected",
		Args: cobra.NoArgs,
		ParamFields: append(
			[]cobra.HelpField{
				requiredField("request-id", "string", "Target expense request id"),
			},
			append(
				payloadSourceFields("expense update payload"),
				resultField("submitted", createUpdateResults),
			)...,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "EX-", "--request-id")
			if err != nil {
				return err
			}
			raw, err := readPayloadInput(cmd, input)
			if err != nil {
				return err
			}
			payload, err := parseExpensePayload(raw)
			if err != nil {
				return err
			}
			status, err := parseResult(result, createUpdateResults)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:       "expense",
				Action:     "update",
				Status:     status,
				RequestID:  id,
				Summary:    buildExpenseSummary(payload),
				Validation: "passed",
			})
		},
	}

	addPayloadInputFlags(cmd, &input)
	cmd.Flags().StringVar(&requestID, "request-id", "", "Target expense request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(createUpdateResults))
	return cmd
}

func newDeleteExpenseCommand() *cobra.Command {
	var requestID string
	result := "deleted"

	cmd := &cobra.Command{
		Use:         "delete-expense",
		Short:       "Delete a mock expense reimbursement",
		Description: "Return a mock expense deletion receipt for the given request id.",
		Example:     "mock delete-expense --request-id EX-14C0A7B992\nmock delete-expense --request-id EX-14C0A7B992 --result not_found",
		Args:        cobra.NoArgs,
		ParamFields: []cobra.HelpField{
			requiredField("request-id", "string", "Target expense request id"),
			resultField("deleted", deleteResults),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "EX-", "--request-id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, deleteResults)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:      "expense",
				Action:    "delete",
				Status:    status,
				RequestID: id,
			})
		},
	}

	cmd.Flags().StringVar(&requestID, "request-id", "", "Target expense request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(deleteResults))
	return cmd
}

func newCreateProcurementCommand() *cobra.Command {
	var input payloadInputOptions
	result := "submitted"

	cmd := &cobra.Command{
		Use:         "create-procurement",
		Short:       "Create a mock procurement request",
		Description: "Validate and create a mock procurement request from JSON input.",
		Example: "mock create-procurement --payload '{\"requester_id\":\"E1001\",\"department\":\"engineering\",\"budget_code\":\"RD-2026-001\",\"reason\":\"team expansion\",\"delivery_city\":\"Shanghai\",\"items\":[{\"name\":\"MacBook Pro\",\"quantity\":2,\"unit_price\":18999,\"vendor\":\"Apple\"}],\"approvers\":[\"MGR100\",\"FIN200\"],\"requested_at\":\"2026-04-14T11:00:00+08:00\"}'\n" +
			"mock create-procurement --payload-file ./procurement.json --result approved\n" +
			"printf '{\"requester_id\":\"E1001\",\"department\":\"engineering\",\"budget_code\":\"RD-2026-001\",\"reason\":\"team expansion\",\"delivery_city\":\"Shanghai\",\"items\":[{\"name\":\"MacBook Pro\",\"quantity\":2,\"unit_price\":18999,\"vendor\":\"Apple\"}],\"approvers\":[\"MGR100\",\"FIN200\"],\"requested_at\":\"2026-04-14T11:00:00+08:00\"}' | mock create-procurement --payload-stdin --result rejected",
		Args: cobra.NoArgs,
		ParamFields: append(
			payloadSourceFields("procurement request"),
			resultField("submitted", createUpdateResults),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readPayloadInput(cmd, input)
			if err != nil {
				return err
			}
			payload, err := parseProcurementPayload(raw)
			if err != nil {
				return err
			}
			status, err := parseResult(result, createUpdateResults)
			if err != nil {
				return err
			}
			requestID, err := stableRequestID("PR-", payload)
			if err != nil {
				return failuref("build procurement request id: %v", err)
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:       "procurement",
				Action:     "create",
				Status:     status,
				RequestID:  requestID,
				Summary:    buildProcurementSummary(payload),
				Validation: "passed",
			})
		},
	}

	addPayloadInputFlags(cmd, &input)
	cmd.Flags().StringVar(&result, "result", result, resultUsage(createUpdateResults))
	return cmd
}

func newGetProcurementCommand() *cobra.Command {
	var requestID string
	result := "found"

	cmd := &cobra.Command{
		Use:         "get-procurement",
		Short:       "Get a mock procurement request",
		Description: "Return a mock procurement request record for the given request id.",
		Example:     "mock get-procurement --request-id PR-BA08D42C31\nmock get-procurement --request-id PR-BA08D42C31 --result not_found",
		Args:        cobra.NoArgs,
		ParamFields: []cobra.HelpField{
			requiredField("request-id", "string", "Target procurement request id"),
			resultField("found", getResults),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "PR-", "--request-id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, getResults)
			if err != nil {
				return err
			}

			response := actionResponse{
				Type:      "procurement",
				Action:    "get",
				Status:    status,
				RequestID: id,
			}
			if status == "found" {
				response.Record = mockProcurementRecord(id)
			}
			return writeJSON(cmd.OutOrStdout(), response)
		},
	}

	cmd.Flags().StringVar(&requestID, "request-id", "", "Target procurement request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(getResults))
	return cmd
}

func newUpdateProcurementCommand() *cobra.Command {
	var input payloadInputOptions
	var requestID string
	result := "submitted"

	cmd := &cobra.Command{
		Use:         "update-procurement",
		Short:       "Update a mock procurement request",
		Description: "Validate and update a mock procurement request from JSON input.",
		Example: "mock update-procurement --request-id PR-BA08D42C31 --payload '{\"requester_id\":\"E1001\",\"department\":\"engineering\",\"budget_code\":\"RD-2026-001\",\"reason\":\"team expansion\",\"delivery_city\":\"Shanghai\",\"items\":[{\"name\":\"MacBook Pro\",\"quantity\":2,\"unit_price\":18999,\"vendor\":\"Apple\"}],\"approvers\":[\"MGR100\",\"FIN200\"],\"requested_at\":\"2026-04-14T11:00:00+08:00\"}'\n" +
			"mock update-procurement --request-id PR-BA08D42C31 --payload-file ./procurement-update.json --result approved\n" +
			"printf '{\"requester_id\":\"E1001\",\"department\":\"engineering\",\"budget_code\":\"RD-2026-001\",\"reason\":\"team expansion\",\"delivery_city\":\"Shanghai\",\"items\":[{\"name\":\"MacBook Pro\",\"quantity\":2,\"unit_price\":18999,\"vendor\":\"Apple\"}],\"approvers\":[\"MGR100\",\"FIN200\"],\"requested_at\":\"2026-04-14T11:00:00+08:00\"}' | mock update-procurement --request-id PR-BA08D42C31 --payload-stdin --result rejected",
		Args: cobra.NoArgs,
		ParamFields: append(
			[]cobra.HelpField{
				requiredField("request-id", "string", "Target procurement request id"),
			},
			append(
				payloadSourceFields("procurement update payload"),
				resultField("submitted", createUpdateResults),
			)...,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "PR-", "--request-id")
			if err != nil {
				return err
			}
			raw, err := readPayloadInput(cmd, input)
			if err != nil {
				return err
			}
			payload, err := parseProcurementPayload(raw)
			if err != nil {
				return err
			}
			status, err := parseResult(result, createUpdateResults)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:       "procurement",
				Action:     "update",
				Status:     status,
				RequestID:  id,
				Summary:    buildProcurementSummary(payload),
				Validation: "passed",
			})
		},
	}

	addPayloadInputFlags(cmd, &input)
	cmd.Flags().StringVar(&requestID, "request-id", "", "Target procurement request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(createUpdateResults))
	return cmd
}

func newDeleteProcurementCommand() *cobra.Command {
	var requestID string
	result := "deleted"

	cmd := &cobra.Command{
		Use:         "delete-procurement",
		Short:       "Delete a mock procurement request",
		Description: "Return a mock procurement deletion receipt for the given request id.",
		Example:     "mock delete-procurement --request-id PR-BA08D42C31\nmock delete-procurement --request-id PR-BA08D42C31 --result not_found",
		Args:        cobra.NoArgs,
		ParamFields: []cobra.HelpField{
			requiredField("request-id", "string", "Target procurement request id"),
			resultField("deleted", deleteResults),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := validateRequestID(requestID, "PR-", "--request-id")
			if err != nil {
				return err
			}
			status, err := parseResult(result, deleteResults)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), actionResponse{
				Type:      "procurement",
				Action:    "delete",
				Status:    status,
				RequestID: id,
			})
		},
	}

	cmd.Flags().StringVar(&requestID, "request-id", "", "Target procurement request id")
	cmd.Flags().StringVar(&result, "result", result, resultUsage(deleteResults))
	return cmd
}

func addPayloadInputFlags(cmd *cobra.Command, input *payloadInputOptions) {
	cmd.Flags().StringVar(&input.Payload, "payload", "", "Raw JSON payload")
	cmd.Flags().StringVar(&input.PayloadFile, "payload-file", "", "Path to a JSON payload file")
	cmd.Flags().BoolVar(&input.PayloadStdin, "payload-stdin", false, "Read the JSON payload from stdin")
}

func payloadSourceFields(subject string) []cobra.HelpField {
	return []cobra.HelpField{
		oneOfField("payload", "string", fmt.Sprintf("Raw JSON payload for the %s", subject)),
		oneOfField("payload-file", "string", fmt.Sprintf("Path to a JSON payload file for the %s", subject)),
		oneOfField("payload-stdin", "bool", fmt.Sprintf("Read the %s JSON payload from stdin", subject)),
	}
}

func oneOfField(name, fieldType, description string) cobra.HelpField {
	return cobra.HelpField{
		Name:        name,
		Type:        fieldType,
		Required:    "one-of",
		Default:     "-",
		Description: description,
	}
}

func resultField(defaultValue string, results []string) cobra.HelpField {
	return optionalField("result", "string", defaultValue, resultUsage(results))
}

func resultUsage(results []string) string {
	return "Mock result: " + strings.Join(results, "|")
}

func readPayloadInput(cmd *cobra.Command, input payloadInputOptions) (string, error) {
	provided := 0
	if strings.TrimSpace(input.Payload) != "" {
		provided++
	}
	if strings.TrimSpace(input.PayloadFile) != "" {
		provided++
	}
	if input.PayloadStdin {
		provided++
	}
	if provided == 0 {
		return "", fmt.Errorf("exactly one of --payload, --payload-file, or --payload-stdin is required")
	}
	if provided > 1 {
		return "", fmt.Errorf("only one of --payload, --payload-file, or --payload-stdin may be used")
	}

	switch {
	case strings.TrimSpace(input.Payload) != "":
		return input.Payload, nil
	case strings.TrimSpace(input.PayloadFile) != "":
		data, err := os.ReadFile(strings.TrimSpace(input.PayloadFile))
		if err != nil {
			return "", failuref("read --payload-file %q: %v", input.PayloadFile, err)
		}
		return string(data), nil
	default:
		data, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return "", failuref("read --payload-stdin: %v", err)
		}
		return string(data), nil
	}
}

func parseLeavePayload(raw string, requireRequestID bool) (leavePayload, error) {
	var payload leavePayload
	if err := decodePayload(raw, &payload); err != nil {
		return leavePayload{}, err
	}
	if err := validateLeaveRequired(payload, requireRequestID); err != nil {
		return leavePayload{}, err
	}
	if err := validateLeaveBusiness(payload); err != nil {
		return leavePayload{}, err
	}
	return payload, nil
}

func parseExpensePayload(raw string) (expensePayload, error) {
	var payload expensePayload
	if err := decodePayload(raw, &payload); err != nil {
		return expensePayload{}, err
	}
	if err := validateExpenseRequired(payload); err != nil {
		return expensePayload{}, err
	}
	if err := validateExpenseBusiness(payload); err != nil {
		return expensePayload{}, err
	}
	return payload, nil
}

func parseProcurementPayload(raw string) (procurementPayload, error) {
	var payload procurementPayload
	if err := decodePayload(raw, &payload); err != nil {
		return procurementPayload{}, err
	}
	if err := validateProcurementRequired(payload); err != nil {
		return procurementPayload{}, err
	}
	if err := validateProcurementBusiness(payload); err != nil {
		return procurementPayload{}, err
	}
	return payload, nil
}

func decodePayload(raw string, target any) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fmt.Errorf("payload input is empty")
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &object); err != nil {
		return fmt.Errorf("invalid payload JSON: %w", err)
	}

	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("invalid payload fields: %w", err)
	}
	if dec.More() {
		return fmt.Errorf("invalid payload JSON: multiple JSON values")
	}
	return nil
}

func validateLeaveRequired(payload leavePayload, requireRequestID bool) error {
	if requireRequestID && strings.TrimSpace(payload.RequestID) == "" {
		return fmt.Errorf("payload.request_id is required")
	}

	required := []struct {
		value string
		name  string
	}{
		{payload.EmployeeID, "employee_id"},
		{payload.EmployeeName, "employee_name"},
		{payload.LeaveType, "leave_type"},
		{payload.StartDate, "start_date"},
		{payload.EndDate, "end_date"},
		{payload.Reason, "reason"},
		{payload.HandoverTo, "handover_to"},
		{payload.UrgentContact, "urgent_contact"},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("payload.%s is required", field.name)
		}
	}
	if _, err := time.Parse("2006-01-02", payload.StartDate); err != nil {
		return fmt.Errorf("payload.start_date must use YYYY-MM-DD")
	}
	if _, err := time.Parse("2006-01-02", payload.EndDate); err != nil {
		return fmt.Errorf("payload.end_date must use YYYY-MM-DD")
	}
	return nil
}

func validateLeaveBusiness(payload leavePayload) error {
	switch payload.LeaveType {
	case "annual", "sick", "personal":
	default:
		return failuref("unsupported leave_type %q", payload.LeaveType)
	}
	if payload.Days <= 0 {
		return failuref("leave days must be greater than 0")
	}

	startDate, _ := time.Parse("2006-01-02", payload.StartDate)
	endDate, _ := time.Parse("2006-01-02", payload.EndDate)
	if startDate.After(endDate) {
		return failuref("leave start_date must be on or before end_date")
	}
	return nil
}

func validateExpenseRequired(payload expensePayload) error {
	required := []struct {
		value string
		name  string
	}{
		{payload.EmployeeID, "employee_id"},
		{payload.Department, "department"},
		{payload.ExpenseType, "expense_type"},
		{payload.Currency, "currency"},
		{payload.SubmittedAt, "submitted_at"},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("payload.%s is required", field.name)
		}
	}
	if _, err := time.Parse(time.RFC3339, payload.SubmittedAt); err != nil {
		return fmt.Errorf("payload.submitted_at must use RFC3339")
	}
	for i, item := range payload.Items {
		if err := validateExpenseItemRequired(i, item); err != nil {
			return err
		}
	}
	return nil
}

func validateExpenseItemRequired(index int, item expenseLineItem) error {
	required := []struct {
		value string
		name  string
	}{
		{item.Category, "category"},
		{item.InvoiceID, "invoice_id"},
		{item.OccurredOn, "occurred_on"},
		{item.Description, "description"},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("payload.items[%d].%s is required", index, field.name)
		}
	}
	if _, err := time.Parse("2006-01-02", item.OccurredOn); err != nil {
		return fmt.Errorf("payload.items[%d].occurred_on must use YYYY-MM-DD", index)
	}
	return nil
}

func validateExpenseBusiness(payload expensePayload) error {
	if len(payload.Items) == 0 {
		return failuref("expense items must contain at least one entry")
	}
	switch payload.Currency {
	case "CNY", "USD":
	default:
		return failuref("unsupported currency %q", payload.Currency)
	}
	if payload.TotalAmount <= 0 {
		return failuref("expense total_amount must be greater than 0")
	}

	var totalCents int64
	for _, item := range payload.Items {
		if item.Amount <= 0 {
			return failuref("expense item amount must be greater than 0")
		}
		totalCents += amountToCents(item.Amount)
	}
	if totalCents != amountToCents(payload.TotalAmount) {
		return failuref("expense total_amount does not match item amount sum")
	}
	return nil
}

func validateProcurementRequired(payload procurementPayload) error {
	required := []struct {
		value string
		name  string
	}{
		{payload.RequesterID, "requester_id"},
		{payload.Department, "department"},
		{payload.BudgetCode, "budget_code"},
		{payload.Reason, "reason"},
		{payload.DeliveryCity, "delivery_city"},
		{payload.RequestedAt, "requested_at"},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("payload.%s is required", field.name)
		}
	}
	if _, err := time.Parse(time.RFC3339, payload.RequestedAt); err != nil {
		return fmt.Errorf("payload.requested_at must use RFC3339")
	}
	for i, item := range payload.Items {
		if err := validateProcurementItemRequired(i, item); err != nil {
			return err
		}
	}
	for i, approver := range payload.Approvers {
		if strings.TrimSpace(approver) == "" {
			return fmt.Errorf("payload.approvers[%d] is required", i)
		}
	}
	return nil
}

func validateProcurementItemRequired(index int, item procurementLineItem) error {
	required := []struct {
		value string
		name  string
	}{
		{item.Name, "name"},
		{item.Vendor, "vendor"},
	}
	for _, field := range required {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("payload.items[%d].%s is required", index, field.name)
		}
	}
	return nil
}

func validateProcurementBusiness(payload procurementPayload) error {
	if len(payload.Items) == 0 {
		return failuref("procurement items must contain at least one entry")
	}
	if len(payload.Approvers) < 2 {
		return failuref("procurement approvers must contain at least two entries")
	}

	totalAmount := 0.0
	for _, item := range payload.Items {
		if item.Quantity <= 0 {
			return failuref("procurement item quantity must be greater than 0")
		}
		if item.UnitPrice <= 0 {
			return failuref("procurement item unit_price must be greater than 0")
		}
		totalAmount += float64(item.Quantity) * item.UnitPrice
	}
	if totalAmount > procurementBudgetLimit {
		return failuref("budget exceeded: total_amount %.2f exceeds %.2f", totalAmount, procurementBudgetLimit)
	}
	return nil
}

func parseResult(raw string, allowed []string) (string, error) {
	raw = strings.TrimSpace(raw)
	for _, next := range allowed {
		if raw == next {
			return raw, nil
		}
	}
	return "", fmt.Errorf("invalid --result %q: must be one of %s", raw, strings.Join(allowed, ", "))
}

func validateRequestID(raw, prefix, field string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	if !strings.HasPrefix(raw, prefix) {
		return "", fmt.Errorf("%s must start with %q", field, prefix)
	}
	suffix := strings.TrimPrefix(raw, prefix)
	if len(suffix) < 6 {
		return "", fmt.Errorf("%s must include an identifier suffix", field)
	}
	return raw, nil
}

func buildLeaveSummary(payload leavePayload) leaveSummary {
	return leaveSummary{
		EmployeeID:   payload.EmployeeID,
		EmployeeName: payload.EmployeeName,
		LeaveType:    payload.LeaveType,
		StartDate:    payload.StartDate,
		EndDate:      payload.EndDate,
		Days:         payload.Days,
		Reason:       payload.Reason,
		HandoverTo:   payload.HandoverTo,
	}
}

func buildExpenseSummary(payload expensePayload) expenseSummary {
	return expenseSummary{
		EmployeeID:     payload.EmployeeID,
		Department:     payload.Department,
		ExpenseType:    payload.ExpenseType,
		Currency:       payload.Currency,
		TotalAmount:    payload.TotalAmount,
		ItemCount:      len(payload.Items),
		SubmittedAt:    payload.SubmittedAt,
		PrimaryInvoice: payload.Items[0].InvoiceID,
	}
}

func buildProcurementSummary(payload procurementPayload) procurementSummary {
	return procurementSummary{
		RequesterID:   payload.RequesterID,
		Department:    payload.Department,
		BudgetCode:    payload.BudgetCode,
		DeliveryCity:  payload.DeliveryCity,
		TotalAmount:   procurementTotalAmount(payload.Items),
		ItemCount:     len(payload.Items),
		ApproverCount: len(payload.Approvers),
		Approvers:     append([]string(nil), payload.Approvers...),
		RequestedAt:   payload.RequestedAt,
	}
}

func mockLeaveRecord(requestID string) leaveRecord {
	return leaveRecord{
		RequestID:     requestID,
		EmployeeID:    "E1001",
		EmployeeName:  "Lin",
		LeaveType:     "annual",
		StartDate:     "2026-04-20",
		EndDate:       "2026-04-22",
		Days:          3,
		Reason:        "family_trip",
		HandoverTo:    "E2001",
		UrgentContact: "13800138000",
		Status:        "submitted",
	}
}

func mockExpenseRecord(requestID string) expenseRecord {
	return expenseRecord{
		RequestID:   requestID,
		EmployeeID:  "E1001",
		Department:  "engineering",
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
		Status:      "submitted",
	}
}

func mockProcurementRecord(requestID string) procurementRecord {
	return procurementRecord{
		RequestID:    requestID,
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
		Status:      "submitted",
	}
}

func leavePayloadForStableID(payload leavePayload) leavePayload {
	payload.RequestID = ""
	return payload
}

func procurementTotalAmount(items []procurementLineItem) float64 {
	total := 0.0
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

func stableRequestID(prefix string, payload any) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha1.Sum(data)
	return prefix + strings.ToUpper(hex.EncodeToString(sum[:]))[:10], nil
}

func writeJSON(out io.Writer, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}

func amountToCents(value float64) int64 {
	return int64(math.Round(value * 100))
}
