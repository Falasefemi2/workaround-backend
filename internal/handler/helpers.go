package handler

import (
	"net/http"
	"strconv"
)

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
