package application

import "github.com/camunda/zeebe/clients/go/v8/pkg/pb"

type DeployResourceRequest struct {
	Path string `json:"path"`
}

type DeployResourceResponse struct {
	Key         int64            `json:"key"`
	Deployments []*pb.Deployment `json:"deployments"`
}

type DeployInstanceRequest struct {
	BpmnProcessId string     `json:"bpmnProcessId"`
	Variables     []Variable `json:"variables"`
}

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeployInstanceResponse struct {
	BpmnProcessId string     `json:"bpmnProcessId"`
	Variables     []Variable `json:"variables"`
}
