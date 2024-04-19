FROM golang:latest
LABEL maintainer="Eliot Scott <eliotvscott@gmail.com>"

WORKDIR /

COPY go.mod .
COPY cmd /cmd

CMD ["go", "run", "cmd/main.go"]