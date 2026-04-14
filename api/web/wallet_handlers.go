package web

import (
	"net/http"

	"github.com/masudur-rahman/expense-tracker-bot/services/all"
)

// HandleListWallets handles GET /wallets.
func HandleListWallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := UserFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "missing claims")
		return
	}

	wallets, err := all.GetServices().Wallet.ListWallets(claims.UserID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, wallets)
}
