FROM golang:alpine

RUN apk add --no-cache curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
COPY ./challenge ./challenge
COPY ./openapi.json ./openapi.json
WORKDIR /go/challenge/
RUN task build --force
CMD ["./main_server"]
