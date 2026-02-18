package http

import (
	"net/http"

	"github.com/keel-hq/keel/types"
)

type dailyStats struct {
	Timestamp         int `json:"timestamp"`
	WebhooksReceived  int `json:"webhooksReceived"`
	ApprovalsApproved int `json:"approvalsApproved"`
	ApprovalsRejected int `json:"approvalsRejected"`
	Updates           int `json:"updates"`
}

// statsHandler godoc
// @Summary Get statistics
// @Description Returns aggregated statistics about webhooks, approvals, and updates
// @Tags stats
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /v1/stats [get]
func (s *TriggerServer) statsHandler(resp http.ResponseWriter, req *http.Request) {
	stats, err := s.store.AuditStatistics(&types.AuditLogStatsQuery{})
	response(stats, 200, err, resp, req)
}
