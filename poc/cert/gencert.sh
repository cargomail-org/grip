###################### CA1 127.0.0.1.nip.io #####################
# ca1.key
openssl genrsa -out ./ca/ca1.key 2048
# ca1.crt
openssl req -new -key ./ca/ca1.key -x509 -days 3650 -out ./ca/ca1.crt -utf8 -subj /O="127.0.0.1.nip.io"/CN="127.0.0.1.nip.io Root"
# ---------------------------- AS -------------------------------
# as.key
openssl genrsa -out as.key 2048
# as.csr
openssl req -new -nodes -key as.key -out as.csr -utf8 -subj /O="Foo Email Provider"/CN="Authorization Server"
# as.crt
openssl x509 -req -extfile <(printf "subjectAltName=DNS:foo.127.0.0.1.nip.io") -days 365 -in as.csr -CA ./ca/ca1.crt -CAkey ./ca/ca1.key -CAcreateserial -out as.crt

########################## Self-Signed ##########################
# ------------------------ SMTP Client --------------------------
# client.key
openssl genrsa -out client.key 2048
# client.crt
openssl req -x509 -days 365 -nodes -key client.key -out client.crt -subj /O="Foo Email Provider"/CN="SMTP Client" -addext "subjectAltName=email:client@foo.127.0.0.1.nip.io"

###################### CA2 127.0.0.2.nip.io #####################
# ca2.key
openssl genrsa -out ./ca/ca2.key 2048
# ca2.crt
openssl req -new -key ./ca/ca2.key -x509 -days 3650 -out ./ca/ca2.crt -utf8 -subj /O="127.0.0.2.nip.io"/CN="127.0.0.2.nip.io Root"
# ------------------------ SMTP Server --------------------------
# smtpd.key
openssl genrsa -out smtpd.key 2048
# smtpd.csr
openssl req -new -nodes -key smtpd.key -out smtpd.csr -utf8 -subj /O="Bar Email Provider"/CN="SMTP Server"
# smtpd.crt
openssl x509 -req -extfile <(printf "subjectAltName=DNS:bar.127.0.0.2.nip.io") -days 365 -in smtpd.csr -CA ./ca/ca2.crt -CAkey ./ca/ca2.key -CAcreateserial -out smtpd.crt
