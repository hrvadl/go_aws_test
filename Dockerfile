FROM golang:1.21.1-bullseye
RUN apt update && apt upgrade -y
WORKDIR /usr/src/app
COPY . .
RUN go mod download
WORKDIR /usr/src/app/cmd
RUN go build 
EXPOSE 8080
WORKDIR /usr/src/app
CMD ["./cmd/cmd"]
