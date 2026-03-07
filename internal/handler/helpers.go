package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/falasefemi2/workaround-backend/internal/response"
)

func scanUUIDOrError(w http.ResponseWriter, value any, message string) (pgtype.UUID, bool) {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		response.Error(w, http.StatusBadRequest, message)
		return pgtype.UUID{}, false
	}

	return id, true
}

func scanNumericOrError(w http.ResponseWriter, value float64, message string) (pgtype.Numeric, bool) {
	var n pgtype.Numeric
	if err := n.Scan(fmt.Sprintf("%f", value)); err != nil {
		response.Error(w, http.StatusBadRequest, message)
		return pgtype.Numeric{}, false
	}

	return n, true
}

func scanOptionalUUIDOrError(w http.ResponseWriter, value string, message string) (pgtype.UUID, bool) {
	if value == "" {
		return pgtype.UUID{}, true
	}

	return scanUUIDOrError(w, value, message)
}

func parsePagination(r *http.Request) (int32, int32) {
	limit := int32(10)
	offset := int32(0)

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			limit = int32(val)
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			offset = int32(val)
		}
	}

	return limit, offset
}
