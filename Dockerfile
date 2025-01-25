FROM golang:latest AS build

ENV GOPATH=/
WORKDIR /src/
COPY ./ /src

RUN go mod download; CGO_ENABLED=0 go build -o /todo-app-go ./cmd/main.go

FROM alpine:3.17


COPY --from=build /todo-app-go /todo-app-go
COPY ./configs/ /configs/
#ENV ENV=production
#ENV DB_PASSWORD=qwerty
#ENV SIGNING_KEY=fsadfasgagashdjgfasdf5z12b135afg56
#ENV SALT=aihvsop198kgmlk

CMD ["./todo-app-go"]