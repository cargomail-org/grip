########################## Self-Signed ##########################
# ------------------------ cargomail.org Cargomail Server --------------------------
# cargomail.org-client.key
openssl genrsa -out cargomail.org-server.key 2048
# cargomail.org-client.crt
openssl req -x509 -days 3650 -nodes -key cargomail.org-server.key -out cargomail.org-server.crt -subj /O="Cargomail Provider"/CN="_cargomail-dev-server.cargomail.org" -addext "subjectAltName=IP:127.0.0.1"

# ------------------------ cargomail.org Cargomail Client --------------------------
# cargomail.org-client.key
openssl genrsa -out cargomail.org-client.key 2048
# cargomail.org-client.crt
openssl req -x509 -days 3650 -nodes -key cargomail.org-client.key -out cargomail.org-client.crt -subj /O="Cargomail Provider"/CN="_cargomail-dev-client.cargomail.org"

# ------------------------ federizer.org Cargomail Server --------------------------
# federizer.org-client.key
openssl genrsa -out federizer.org-server.key 2048
# federizer.org-client.crt
openssl req -x509 -days 3650 -nodes -key federizer.org-server.key -out federizer.org-server.crt -subj /O="Cargomail Provider"/CN="_cargomail-dev-server.federizer.org" -addext "subjectAltName=IP:127.0.0.3"

# ------------------------ federizer.org Cargomail Client --------------------------
# federizer.org-client.key
openssl genrsa -out federizer.org-client.key 2048
# federizer.org-client.crt
openssl req -x509 -days 3650 -nodes -key federizer.org-client.key -out federizer.org-client.crt -subj /O="Cargomail Provider"/CN="_cargomail-dev-client.federizer.org"

# ------------------------ umabox.org Cargomail Server --------------------------
# umabox.org-client.key
openssl genrsa -out umabox.org-server.key 2048
# umabox.org-client.crt
openssl req -x509 -days 3650 -nodes -key umabox.org-server.key -out umabox.org-server.crt -subj /O="Cargomail Provider"/CN="_cargomail-dev-server.umabox.org" -addext "subjectAltName=IP:127.0.0.2"

# ------------------------ umabox.org Cargomail Client --------------------------
# umabox.org-client.key
openssl genrsa -out umabox.org-client.key 2048
# umabox.org-client.crt
openssl req -x509 -days 3650 -nodes -key umabox.org-client.key -out umabox.org-client.crt -subj /O="Cargomail Provider"/CN="_cargomail-dev-client.umabox.org"
