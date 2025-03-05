# 

## Certificates 

### Create certificates
Use the create.sh script in the ./certs/ folder to create the certificates

## Start the Web Server
Webserver listening on port 8443 and uses the certificates server.pem and server-key.pem from the ./certs folder.

## Run the Client
Client code sends request to server port 8443 and uses ca certificate as server.pem from the ./certs folder.
