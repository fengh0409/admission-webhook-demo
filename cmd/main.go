package main

import (
	"admission-webhook-demo/pkg/config"
	"admission-webhook-demo/pkg/webhook"
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

func main() {
	var option Option
	option.addFlags()
	flag.Parse()
	defer glog.Flush()

	tlsConfig, err := option.configTLS()
	if err != nil {
		panic(err)
	}

	if err := config.NewConfig(option.SidecarFile); err != nil {
		panic(err)
	}

	http.HandleFunc("/mutate", webhook.Mutate)
	http.HandleFunc("/validate", webhook.Validate)
	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", option.Port),
		TLSConfig: tlsConfig,
	}

	glog.Infof("Server listening on %v", option.Port)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		glog.Error(err)
	}

}
