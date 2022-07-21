package application

import (
	"context"
	"encoding/json"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/dulguundd/logError-lib/logger"
	"net/http"
)

type Handler struct {
	zeebeClient *zbc.ClientConfig
}

func (h Handler) DeployResource(w http.ResponseWriter, r *http.Request) {
	var req DeployResourceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		zbClient, err := zbc.NewClient(h.zeebeClient)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, err)
			logger.Error("Cant connect zeebe")
		} else {
			ctx := context.Background()
			response, err := zbClient.NewDeployResourceCommand().AddResourceFile(bpmnFolderLocation + req.Path).Send(ctx)
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err)
				logger.Error("Deploy error")
			} else {
				apiResponse := DeployResourceResponse{
					Key:         response.GetKey(),
					Deployments: response.GetDeployments(),
				}
				writeResponse(w, http.StatusOK, apiResponse)
			}
		}
	}
}

func (h Handler) CreateInstance(w http.ResponseWriter, r *http.Request) {
	var req DeployInstanceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		zbClient, err := zbc.NewClient(h.zeebeClient)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, err)
			logger.Error("Cant connect zeebe")
		} else {
			variables := make(map[string]interface{})
			for variableCount, variablesValue := range req.Variables {
				variables[req.Variables[variableCount].Name] = variablesValue.Value
			}
			ctx := context.Background()
			request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(req.BpmnProcessId).LatestVersion().VariablesFromMap(variables)
			_, err = request.Send(ctx)
			if err != nil {
				writeResponse(w, http.StatusInternalServerError, err)
				logger.Error("Deploy error")
			} else {
				apiResponse := DeployInstanceResponse{
					BpmnProcessId: req.BpmnProcessId,
					Variables:     req.Variables,
				}
				writeResponse(w, http.StatusOK, apiResponse)
			}
		}
	}
}

func newHandler(zeebeConfig ZeebeConfig) *zbc.ClientConfig {
	return &zbc.ClientConfig{
		GatewayAddress:         zeebeConfig.zeebeAddress,
		UsePlaintextConnection: true,
	}
}
