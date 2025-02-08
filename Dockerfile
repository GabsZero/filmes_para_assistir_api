FROM golang:1.23.6-alpine3.21 

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .


RUN go build -v -o ./build

RUN chown 1000:1000 -R .

USER 1000

CMD [ "build/assistir_filmes" ]
