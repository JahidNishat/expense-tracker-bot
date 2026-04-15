package web

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/masudur-rahman/expense-tracker-bot/api/handlers"
	"github.com/masudur-rahman/expense-tracker-bot/models/gqtypes"
	"github.com/masudur-rahman/expense-tracker-bot/modules/convert"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"
)

// HandleGetReport handles GET /summary/report.
func HandleGetReport(w http.ResponseWriter, r *http.Request) {
	claims, ok := UserFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "missing claims")
		return
	}

	duration := handlers.SummaryDuration(r.URL.Query().Get("duration"))
	if duration == "" {
		duration = handlers.DurationThisMonth
	}

	now := time.Now()
	startTime := handlers.CalculateStartTime(duration)

	svc := all.GetServices()
	user, err := svc.User.GetUserByID(claims.UserID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "profile_failed", err.Error())
		return
	}

	txns, err := svc.Txn.ListTransactionsByTime(user.ID, "", startTime.Unix(), now.Unix())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}

	report := gqtypes.Report{
		Name:      fmt.Sprintf("%v %v", user.FirstName, user.LastName),
		StartDate: startTime,
		EndDate:   now,
	}
	txnApis := make([]gqtypes.Transaction, 0, len(txns))
	for _, txn := range txns {
		txnApis = append(txnApis, convert.ToTransactionAPIFormat(txn))
	}
	report.Transactions = txnApis

	summary, err := handlers.BuildSummary(svc, txns)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "summary_failed", err.Error())
		return
	}
	report.Summary = summary

	report.TypeSummary = gqtypes.SortMapToSlice(summary.Type)
	report.CategorySummary, err = handlers.BuildTypeSeparatedSummary(svc, txns, handlers.CategoryKeyFn, svc.Txn.GetTxnCategoryName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "cat_summary_failed", err.Error())
		return
	}
	report.SubcategorySummary, err = handlers.BuildTypeSeparatedSummary(svc, txns, handlers.SubcategoryKeyFn, svc.Txn.GetTxnSubcategoryName)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "subcat_summary_failed", err.Error())
		return
	}

	report.TotalAmount, report.NetBalance = handlers.ComputeTotals(txns)

	pdfFile, err := handlers.GenerateTransactionReportFromTemplate(report, "")
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "pdf_failed", err.Error())
		return
	}
	defer os.Remove(pdfFile)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"financial_transaction_report_%s.pdf\"", duration))
	http.ServeFile(w, r, pdfFile)
}
