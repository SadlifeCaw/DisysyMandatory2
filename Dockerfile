FROM golang:1.17

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV PORT $port

CMD cd DisysyMandatory2
CMD cd Node
CMD go run . $PORT