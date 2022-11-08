package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2/klogr"
	"net/http"
	"nodeSelector/handlers"
	"nodeSelector/utils"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	logger = utils.GetLogger()
)

func main() {

	scheme := runtime.NewScheme()
	decoder, err := admission.NewDecoder(scheme)
	if err != nil {
		logger.Fatalln(err)
	}

	deploymentHandler := handlers.WebHookDeploymentHandler{Decoder: decoder}
	deployWebhook := admission.Webhook{
		Handler: &deploymentHandler,
	}
	_, err = inject.LoggerInto(klogr.New(), &deployWebhook)
	if err != nil {
		logger.Fatalln(err)
	}
	http.HandleFunc("/dp", deployWebhook.ServeHTTP)

	statefulSetHandler := handlers.WebHookStatefulSetHandler{Decoder: decoder}
	statefulSetWebhook := admission.Webhook{
		Handler: &statefulSetHandler,
	}

	_, err = inject.LoggerInto(klogr.New(), &statefulSetWebhook)
	if err != nil {
		logger.Fatalln(err)
	}
	http.HandleFunc("/sts", statefulSetWebhook.ServeHTTP)

	err = http.ListenAndServeTLS(":443", "/etc/webhook/cert.pem", "/etc/webhook/key.pem", nil)
	if err != nil {
		logger.Fatalln(err)
	}
}
