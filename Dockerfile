FROM golang:alpine
COPY src src
WORKDIR src
RUN go build -o /bin/godrunk
ENV PORT=8080
EXPOSE 8080
CMD /bin/godrunk
