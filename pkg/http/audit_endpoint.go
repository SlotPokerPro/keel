package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/keel-hq/keel/types"
)

// adminAuditLogHandler godoc
// @Summary Get audit logs
// @Description Returns paginated audit logs with optional filtering
// @Tags audit
// @Produce json
// @Param limit query int false "Maximum number of results"
// @Param offset query int false "Offset for pagination"
// @Param filter query string false "Comma-separated list of resource kinds to filter"
// @Param email query string false "Filter by email"
// @Success 200 {object} auditLogsResponse
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /v1/audit [get]
func (s *TriggerServer) adminAuditLogHandler(resp http.ResponseWriter, req *http.Request) {

	query := &types.AuditLogQuery{}
	limitS := req.URL.Query().Get("limit")
	if limitS != "" {
		l, err := strconv.Atoi(limitS)
		if err == nil {
			query.Limit = l
		}
	}

	offsetS := req.URL.Query().Get("offset")
	if offsetS != "" {
		o, err := strconv.Atoi(offsetS)
		if err == nil {
			query.Offset = o
		}
	}

	kindFilter := req.URL.Query().Get("filter")
	if kindFilter != "" {
		kinds := strings.Split(kindFilter, ",")
		query.ResourceKindFilter = kinds
	}

	emailFilter := req.URL.Query().Get("email")
	if emailFilter != "" {
		query.Email = strings.TrimSpace(emailFilter)
	}

	entries, err := s.store.GetAuditLogs(query)
	if err != nil {
		response(nil, 500, err, resp, req)
		return
	}

	result := auditLogsResponse{
		Data:   entries,
		Offset: query.Offset,
		Limit:  query.Limit,
	}

	count, err := s.store.AuditLogsCount(query)
	if err == nil {
		result.Total = count
	}

	response(result, http.StatusOK, err, resp, req)
}

type auditLogsResponse struct {
	Data   []*types.AuditLog `json:"data"`
	Total  int               `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}
