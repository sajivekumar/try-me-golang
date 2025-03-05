cd ./certs
rm -rf *.pem 

openssl req -newkey rsa:2048 \
   -new -nodes -x509 \
 -days 3650 \
  -out server.pem \
  -keyout server-key.pem \
  -addext "subjectAltName = DNS:localhost" \
  -subj "/C=US/ST=Bangalore/L=Mahadevpura/O=MyOrg/OU=My Unit/CN=localhost"


openssl req -newkey rsa:2048 \
   -new -nodes -x509 \
 -days 3650 \
  -out client.pem \
  -keyout client-key.pem \
  -addext "subjectAltName = DNS:localhost" \
  -subj "/C=US/ST=Bangalore/L=Mahadevpura/O=MyOrg/OU=My Unit/CN=localhost"
