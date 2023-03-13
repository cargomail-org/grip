curl \
--key ../../cert/client.key \
--cert ../../cert/client.crt \
--url     https://foo.127.0.0.1.nip.io:8080/token.oauth2 \
--request POST \
--header  'Authorization: Basic cnMwODpsb25nLXNlY3VyZS1yYW5kb20tc2VjcmV0' \
--header  'Content-Type: application/x-www-form-urlencoded' \
-d 'grant_type=token-exchange' \
-d 'client_id=myClientID' \
-d 'client_secret=myClientPassword'