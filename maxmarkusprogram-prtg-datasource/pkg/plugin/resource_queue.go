package plugin

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const MaxQueueSize = 100

var (
	requestQueue = make([]*ResourceRequest, 0)
	queueLock    sync.Mutex
)

type ResourceRequest struct {
	Request *backend.CallResourceRequest
	Sender  backend.CallResourceResponseSender
}

func (d *Datasource) processQueuedRequests() error {
	queueLock.Lock()
	defer queueLock.Unlock()

	if len(requestQueue) == 0 {
		return nil
	}

	if len(requestQueue) > MaxQueueSize {
		d.logger.Warn("Request queue overflow, dropping old requests")
		requestQueue = requestQueue[len(requestQueue)-MaxQueueSize:]
	}

	var lastError error
	for _, req := range requestQueue {
		err := d.processRequest(req)
		if err != nil {
			d.logger.Error("Failed to process request", "error", err)
			lastError = err
		}
	}

	requestQueue = requestQueue[:0]
	return lastError
}

func (d *Datasource) processRequest(req *ResourceRequest) error {
	path := req.Request.Path
	d.logger.Debug("Processing request", "path", path)

	switch {
	case strings.HasPrefix(path, "groups"):
		return d.handleGetGroups(req.Sender)

	case strings.HasPrefix(path, "devices/"):
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return sendErrorResponse(req.Sender, "group parameter is required", http.StatusBadRequest)
		}
		return d.handleGetDevices(req.Sender, pathParts[1])

	case strings.HasPrefix(path, "sensors/"):
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return sendErrorResponse(req.Sender, "device parameter is required", http.StatusBadRequest)
		}
		return d.handleGetSensors(req.Sender, pathParts[1])

	case strings.HasPrefix(path, "channels/"):
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			return sendErrorResponse(req.Sender, "sensor parameter is required", http.StatusBadRequest)
		}
		return d.handleGetChannel(req.Sender, pathParts[1])

	default:
		return sendErrorResponse(req.Sender, "invalid API endpoint", http.StatusNotFound)
	}
}

func sendErrorResponse(sender backend.CallResourceResponseSender, message string, statusCode int) error {
	errorResponse := map[string]string{"error": message}
	errorJSON, _ := json.Marshal(errorResponse)
	return sender.Send(&backend.CallResourceResponse{
		Status:  statusCode,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    errorJSON,
	})
}
