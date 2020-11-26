FROM golang:alpine
RUN apk add openssl
COPY src src
WORKDIR src
RUN go build -o /bin/godrunk
ENV PORT=8080
EXPOSE 8080
WORKDIR /bin
COPY start.sh .
CMD start.sh
