package server

import (
	"crypto/tls"
	"forum/config"
	"log"
	"net/http"
	"time"
)

const (
	CertFilePath = "./tls/certificate.crt"
	KeyFilePath  = "./tls/kamal.key"
)

func Runserver(config config.Config, mux http.Handler) error {
	serverTLSCert, err := tls.LoadX509KeyPair(CertFilePath, KeyFilePath)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverTLSCert},
		MinVersion:   tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:         config.Host + config.Port,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	log.Printf("server is listening on https://%s%s", config.Host, config.Port)

	err = server.ListenAndServeTLS("", "")

	return err
}
