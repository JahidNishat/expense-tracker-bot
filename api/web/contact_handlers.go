package web

import (
	"net/http"

	"github.com/masudur-rahman/expense-tracker-bot/services/all"
)

// HandleListContacts handles GET /contacts.
func HandleListContacts(w http.ResponseWriter, r *http.Request) {
	claims, ok := UserFromContext(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "unauthorized", "missing claims")
		return
	}

	contacts, err := all.GetServices().Contact.ListContacts(claims.UserID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, contacts)
}
