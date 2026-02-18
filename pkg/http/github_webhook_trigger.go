package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/keel-hq/keel/types"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

var newGithubWebhooksCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "Github_webhook_requests_total",
		Help: "How many /v1/webhooks/github requests processed, partitioned by image.",
	},
	[]string{"image"},
)

func init() {
	prometheus.MustRegister(newGithubWebhooksCounter)
}

type githubRegistryPackageWebhook struct {
	Action          string `json:"action"`
	RegistryPackage struct {
		Name           string `json:"name"`
		PackageType    string `json:"package_type"`
		PackageVersion struct {
			Version string `json:"version"`
		} `json:"package_version"`
		UpdatedAt string `json:"updated_at"`
	} `json:"registry_package"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type githubPackageV2Webhook struct {
	Action  string `json:"action"`
	Package struct {
		Id             int    `json:"id"`
		Name           string `json:"name"`
		Namespace      string `json:"namespace"`
		Ecosystem      string `json:"ecosystem"`
		PackageVersion struct {
			Name              string `json:"name"`
			ContainerMetadata struct {
				Tag struct {
					Name   string `json:"name"`
					Digest string `json:"digest"`
				} `json:"tag"`
			} `json:"container_metadata"`
		} `json:"package_version"`
	} `json:"package"`
}

// githubHandler godoc
// @Summary Trigger GitHub webhook
// @Description Receives and processes GitHub webhook notifications for container registry events (both GitHub Packages and GitHub Container Registry)
// @Tags webhooks
// @Accept json
// @Produce plain
// @Param X-GitHub-Event header string true "GitHub event type (package or registry_package)"
// @Param payload body object true "GitHub webhook payload"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad request"
// @Router /v1/webhooks/github [post]
func (s *TriggerServer) githubHandler(resp http.ResponseWriter, req *http.Request) {
	// GitHub provides different webhook events for each registry.
	// Github Package uses 'registry_package'
	// Github Container Registry uses 'package_v2'
	// events can be classified as 'X-GitHub-Event' in Request Header.
	hookEvent := req.Header.Get("X-GitHub-Event")

	var imageName, imageTag string

	switch hookEvent {
	case "package":
		payload := new(githubPackageV2Webhook)
		if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("trigger.githubHandler: failed to decode request")
			resp.WriteHeader(http.StatusBadRequest)
			return
		}

		if payload.Package.Ecosystem != "CONTAINER" {
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "registry package type was not container")
		}

		if payload.Package.Name == "" { // github package name
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "repository name cannot be empty")
			return
		}

		if payload.Package.Namespace == "" { // github package org
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "repository namespace cannot be empty")
			return
		}

		if payload.Package.PackageVersion.ContainerMetadata.Tag.Name == "" { // tag
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "repository tag cannot be empty")
			return
		}

		imageName = strings.Join(
			[]string{"ghcr.io", payload.Package.Namespace, payload.Package.Name},
			"/",
		)
		imageTag = payload.Package.PackageVersion.ContainerMetadata.Tag.Name

		break

	case "registry_package":
		payload := new(githubRegistryPackageWebhook)
		if err := json.NewDecoder(req.Body).Decode(payload); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("trigger.githubHandler: failed to decode request")
			resp.WriteHeader(http.StatusBadRequest)
			return
		}

		if payload.RegistryPackage.PackageType != "CONTAINER" {
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "registry package type was not CONTAINER")
		}

		if payload.Repository.FullName == "" { // github package name
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "repository full name cannot be empty")
			return
		}

		if payload.RegistryPackage.Name == "" { // github package name
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "repository package name cannot be empty")
			return
		}

		if payload.RegistryPackage.PackageVersion.Version == "" { // tag
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(resp, "repository tag cannot be empty")
			return
		}

		// XXX <jsonroot>.registry_package.package_version.package_url could work too but it ends with colon
		imageName = strings.Join(
			[]string{"ghcr.io", payload.Repository.FullName},
			"/",
		)
		imageTag = payload.RegistryPackage.PackageVersion.Version

		break
	}

	event := types.Event{}
	event.CreatedAt = time.Now()
	event.TriggerName = "github"
	event.Repository.Name = imageName
	event.Repository.Tag = imageTag

	s.trigger(event)

	resp.WriteHeader(http.StatusOK)

	newGithubWebhooksCounter.With(prometheus.Labels{"image": event.Repository.Name}).Inc()
}
