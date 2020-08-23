package main

import (
	"crypto/tls"
	"flag"
)

const defaultPort = 8443

type Option struct {
	Port        int
	SidecarFile string
	CertFile    string
	KeyFile     string
}

func (o *Option) addFlags() {
	flag.Set("logtostderr", "true")
	flag.IntVar(&o.Port, "port", defaultPort, "server port")
	flag.StringVar(&o.SidecarFile, "sidecar-file", o.SidecarFile, "sidecar file")
	flag.StringVar(&o.CertFile, "tls-cert-file", o.CertFile, "tls cert file")
	flag.StringVar(&o.KeyFile, "tls-private-key-file", o.KeyFile, "tls cert file")
}

func (o *Option) configTLS() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(o.CertFile, o.KeyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}
