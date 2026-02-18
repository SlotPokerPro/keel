package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/keel-hq/keel/types"
)

type resourcePolicyUpdateRequest struct {
	Policy     string `json:"policy"`
	Identifier string `json:"identifier"`
	Provider   string `json:"provider"`
}

// policyUpdateHandler godoc
// @Summary Update resource policy
// @Description Updates the Keel policy for a specific resource
// @Tags policies
// @Accept json
// @Produce json
// @Param request body resourcePolicyUpdateRequest true "Policy update request"
// @Success 200 {object} APIResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Resource not found"
// @Security ApiKeyAuth
// @Router /v1/policies [put]
func (s *TriggerServer) policyUpdateHandler(resp http.ResponseWriter, req *http.Request) {
	var policyRequest resourcePolicyUpdateRequest
	dec := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := dec.Decode(&policyRequest)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(resp, "%s", err)
		return
	}

	if policyRequest.Identifier == "" {
		http.Error(resp, "identifier cannot be empty", http.StatusBadRequest)
		return
	}

	for _, v := range s.grc.Values() {
		if v.Identifier == policyRequest.Identifier {

			labels := v.GetLabels()
			delete(labels, types.KeelPolicyLabel)
			v.SetLabels(labels)

			ann := v.GetAnnotations()
			ann[types.KeelPolicyLabel] = policyRequest.Policy

			v.SetAnnotations(ann)

			err := s.kubernetesClient.Update(v)

			response(&APIResponse{Status: "updated"}, 200, err, resp, req)
			return
		}
	}

	resp.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(resp, "resource with identifier '%s' not found", policyRequest.Identifier)
	return
}
