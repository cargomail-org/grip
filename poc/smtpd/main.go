package main

import (
	"bytes"
	"log"
	"net"
	"net/mail"

	"github.com/mhale/smtpd"
)

const (
	certificate string = "../cert/smtpd.crt"
	key         string = "../cert/smtpd.key"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from: %s for: %s with subject: %s", from, to[0], subject)
	return nil
}

func main() {
	// smtpd.ListenAndServe("bar.127.0.0.2.nip.io:2525", mailHandler, "SMTP-GRIP", "")
	smtpd.ListenAndServeTLS("bar.127.0.0.2.nip.io:2525", certificate, key, mailHandler, "SMTP-GRIP", "")
}
