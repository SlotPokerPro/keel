package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keel-hq/keel/types"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

var newNativeWebhooksCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "native_webhook_requests_total",
		Help: "How many /v1/webhooks/native requests processed, partitioned by image.",
	},
	[]string{"image"},
)

func init() {
	prometheus.MustRegister(newNativeWebhooksCounter)
}

// nativeHandler godoc
// @Summary Trigger native webhook
// @Description Triggers a deployment update event directly with repository and tag information
// @Tags webhooks
// @Accept json
// @Produce plain
// @Param payload body types.Repository true "Repository information with name and tag"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Router /v1/webhooks/native [post]
func (s *TriggerServer) nativeHandler(resp http.ResponseWriter, req *http.Request) {
	repo := types.Repository{}
	if err := json.NewDecoder(req.Body).Decode(&repo); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to decode request")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	event := types.Event{}

	if repo.Name == "" {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(resp, "repository name cannot be empty")
		return
	}

	if repo.Tag == "" {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(resp, "repository tag cannot be empty")
		return
	}

	event.Repository = repo
	event.CreatedAt = time.Now()
	event.TriggerName = "native"
	s.trigger(event)

	resp.WriteHeader(http.StatusOK)

	newNativeWebhooksCounter.With(prometheus.Labels{"image": event.Repository.Name}).Inc()
	return
}
