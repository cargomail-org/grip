curl \
--key ../../cert/client.key \
--cert ../../cert/client.crt \
--url "smtp://bar.127.0.0.2.nip.io:2525" --ssl -k \
--mail-from myself@example.com \
--mail-auth eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c \
--mail-rcpt receiver@example.com \
--upload-file email.txt