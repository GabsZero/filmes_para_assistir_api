FROM golang:1.23.6-alpine3.21 AS builder

RUN adduser -D -g '' johnny


WORKDIR /usr/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN mkdir build
RUN go build -v -o ./build/

RUN chown johnny:johnny -R .

USER johnny

CMD [ "build/assistir_filmes" ]
