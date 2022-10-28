package handlers

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

type NodeSelectorHandler struct {
	Client  client.Client
	Decoder *admission.Decoder
}

func (handler *NodeSelectorHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	deployment := &appsv1.Deployment{}
	err := handler.Decoder.Decode(req, deployment)
	logger.Info(deployment.Name)
	if err != nil {
		logger.Info(err.Error())
		return admission.Errored(http.StatusBadRequest, err)
	}
	logger.Info("start!!!!!!!!")
	if strings.Contains(deployment.Name, "ftdb") && strings.Contains(deployment.Name, "bobft") {
		msg, success := tradHandler(deployment)
		logger.Info(msg)
		if !success {
			admission.Allowed(msg)
		}
	}

	marshaledDeploy, err := json.Marshal(deployment)
	if err != nil {
		logger.Info(err.Error())
		return admission.Errored(http.StatusInternalServerError, err)
	}
	logger.Info("end!!!!!!!!")

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledDeploy)
}
