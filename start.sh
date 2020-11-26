openssl req -newkey rsa:2048 -sha256 -nodes -keyout key.pem -x509 -days 365 -out cert.pem -subj "/C=RU/ST=Saint-Petersburg/L=A2 Green/O=A2 Inc/CN=${CN}"

./godrunk