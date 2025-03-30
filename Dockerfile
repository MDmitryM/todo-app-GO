FROM golang:1.23 AS build

ENV GOPATH=/
WORKDIR /src/
COPY ./ /src

RUN go mod download; CGO_ENABLED=0 go build -o /todo-app-go ./cmd/main.go

FROM alpine:3.17


COPY --from=build /todo-app-go /todo-app-go
COPY ./configs/ /configs/

CMD ["./todo-app-go"]