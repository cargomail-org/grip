package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/mail"
)

const (
	certFile string = "../cert/smtpd.crt"
	keyFile  string = "../cert/smtpd.key"
)

func authHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	return true, nil
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("%sReceived mail from: %s for: %s with subject: %s%s", ColorYellow, from, to[0], subject, ColorReset)
	return nil
}

func listenAndServeTLS(addr string, handler Handler, authHandler AuthHandler) error {
	srv := &Server{
		Addr:         addr,
		TLSListener:  false,
		TLSRequired:  true,
		Handler:      handler,
		Appname:      "SMTP-GRIP",
		Hostname:     "",
		AuthHandler:  authHandler,
		AuthRequired: false,
	}
	srv.ConfigureTLS(certFile, keyFile)
	srv.TLSConfig.ClientAuth = tls.RequireAnyClientCert

	mechs := map[string]bool{"OAUTHBEARER": true}
	srv.AuthMechs = mechs

	return srv.ListenAndServe()
}

func main() {
	listenAndServeTLS("bar.127.0.0.2.nip.io:2525", mailHandler, authHandler)
}
