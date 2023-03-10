curl \
--key ../../cert/client.key \
--cert ../../cert/client.crt \
--url "smtp://bar.127.0.0.2.nip.io:2525" --ssl -k \
--mail-from myself@example.com \
--mail-rcpt receiver@example.com \
--upload-file email.txt