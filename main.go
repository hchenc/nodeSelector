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
	nodeSelectorHandler := handlers.NodeSelectorHandler{Decoder: decoder}
	webhook := admission.Webhook{
		Handler: &nodeSelectorHandler,
	}

	_, err = inject.LoggerInto(klogr.New(), &webhook)
	if err != nil {
		logger.Fatalln(err)
	}

	http.HandleFunc("/nodeselector", webhook.ServeHTTP)

	err = http.ListenAndServeTLS(":443", "/etc/webhook/cert.pem", "/etc/webhook/key.pem", nil)
	if err != nil {
		logger.Fatalln(err)
	}
}
