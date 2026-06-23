package controllers

import (
	"net/http"
	"strconv"

	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
)

func FetchAuditLogs(w http.ResponseWriter, r *http.Request) {
	page := 1
	if rawPage := r.URL.Query().Get("page"); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err != nil || parsedPage < 1 {
			responses.Err(w, http.StatusBadRequest, err)
			return
		}
		page = parsedPage
	}

	pageSize := 50
	if rawPageSize := r.URL.Query().Get("page_size"); rawPageSize != "" {
		parsedPageSize, err := strconv.Atoi(rawPageSize)
		if err != nil || parsedPageSize < 1 {
			responses.Err(w, http.StatusBadRequest, err)
			return
		}
		if parsedPageSize > 50 {
			parsedPageSize = 50
		}
		pageSize = parsedPageSize
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewAuditLogsRepository(db)
	logsPage, err := repo.FetchPage(page, pageSize, repositories.AuditLogFilters{
		Query:    r.URL.Query().Get("q"),
		User:     r.URL.Query().Get("user"),
		Action:   r.URL.Query().Get("action"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	})
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, logsPage)
}
