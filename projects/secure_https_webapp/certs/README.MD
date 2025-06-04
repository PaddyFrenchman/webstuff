#CERTS
###to gnerate cert.pem and key.pem

openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
  -keyout key.pem -out cert.pem \
  -subj "/C=US/ST=Local/L=Local/O=Dev/OU=Dev/CN=localhost"

  the two files must be in the server folders.