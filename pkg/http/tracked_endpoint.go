package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keel-hq/keel/types"
)

type trackedImage struct {
	Image        string `json:"image"`
	Trigger      string `json:"trigger"`
	PollSchedule string `json:"pollSchedule"`
	Provider     string `json:"provider"`
	Namespace    string `json:"namespace"`
	Policy       string `json:"policy"`
	Registry     string `json:"registry"`
}

// trackedHandler godoc
// @Summary List tracked images
// @Description Returns a list of all tracked container images with their triggers and policies
// @Tags tracked
// @Produce json
// @Success 200 {array} trackedImage
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /v1/tracked [get]
func (s *TriggerServer) trackedHandler(resp http.ResponseWriter, req *http.Request) {
	trackedImages, err := s.providers.TrackedImages()

	var imgs []trackedImage

	for _, img := range trackedImages {
		imgs = append(imgs, trackedImage{
			Image:        img.Image.Name(),
			Trigger:      img.Trigger.String(),
			PollSchedule: img.PollSchedule,
			Provider:     img.Provider,
			Namespace:    img.Namespace,
			Policy:       img.Policy.Name(),
			Registry:     img.Image.Registry(),
		})
	}

	response(&imgs, 200, err, resp, req)
}

type trackRequest struct {
	Provider   string `json:"provider"`
	Identifier string `json:"identifier"`
	Trigger    string `json:"trigger"`
	Schedule   string `json:"schedule"`
}

// trackSetHandler godoc
// @Summary Update tracking settings
// @Description Updates the trigger type and poll schedule for a tracked resource
// @Tags tracked
// @Accept json
// @Produce json
// @Param request body trackRequest true "Track settings request"
// @Success 200 {object} APIResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Resource not found"
// @Security ApiKeyAuth
// @Router /v1/tracked [put]
func (s *TriggerServer) trackSetHandler(resp http.ResponseWriter, req *http.Request) {

	var trackReq trackRequest
	dec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := dec.Decode(&trackReq)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(resp, "%s", err)
		return
	}

	switch trackReq.Provider {
	case types.ProviderTypeKubernetes.String():
		// ok
	default:
		http.Error(resp, "unsupported provider, supported: 'kubernetes'", http.StatusBadRequest)
		return
	}

	switch trackReq.Trigger {
	case "default", "poll":
		// ok
	default:
		http.Error(resp, "unknown trigger type, supported: 'default', 'poll'", http.StatusBadRequest)
		return
	}

	if trackReq.Schedule != "" {
		_, err = time.ParseDuration(trackReq.Schedule)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "%s", err)
			return
		}
	} else {
		trackReq.Schedule = types.KeelPollDefaultSchedule
	}

	for _, v := range s.grc.Values() {
		if v.Identifier == trackReq.Identifier {

			labels := v.GetLabels()
			delete(labels, types.KeelTriggerLabel)
			v.SetLabels(labels)

			ann := v.GetAnnotations()
			ann[types.KeelTriggerLabel] = trackReq.Trigger
			ann[types.KeelPollScheduleAnnotation] = trackReq.Schedule

			v.SetAnnotations(ann)

			err := s.kubernetesClient.Update(v)

			response(&APIResponse{Status: "updated"}, 200, err, resp, req)
			return
		}
	}

	resp.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(resp, "resource with identifier '%s' not found", trackReq.Identifier)
}
