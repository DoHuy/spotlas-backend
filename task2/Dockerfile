FROM golang:1.19.0

ARG PROJECT=spotlas
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o main .
EXPOSE 3000
CMD ["./main"]