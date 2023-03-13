package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/mail"

	"github.com/mhale/smtpd"
)

const (
	certFile string = "../cert/smtpd.crt"
	keyFile  string = "../cert/smtpd.key"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from: %s for: %s with subject: %s", from, to[0], subject)
	return nil
}

func ListenAndServeTLS(addr string, handler smtpd.Handler) error {
	srv := &smtpd.Server{
		Addr:         addr,
		TLSListener:  false,
		TLSRequired:  true,
		Handler:      handler,
		Appname:      "SMTP-GRIP",
		Hostname:     "",
		AuthRequired: false,
	}
	srv.ConfigureTLS(certFile, keyFile)
	// accept any certificate
	srv.TLSConfig.ClientAuth = tls.RequireAnyClientCert
	return srv.ListenAndServe()
}

func main() {
	ListenAndServeTLS("bar.127.0.0.2.nip.io:2525", mailHandler)
}
