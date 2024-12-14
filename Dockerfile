FROM golang:latest

RUN go version
ENV GOPATH=/
ENV ENV=production
ENV DB_PASSWORD=qwerty
ENV SIGNING_KEY=fsadfasgagashdjgfasdf5z12b135afg56
ENV SALT=aihvsop198kgmlk

COPY ./ ./

RUN go mod download
RUN go build -o todo-app-go ./cmd/main.go

CMD ["./todo-app-go"]