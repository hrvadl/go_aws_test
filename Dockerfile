FROM golang:1.21.1-bullseye
RUN apt update && apt upgrade -y
WORKDIR /usr/src/app
COPY . .
RUN go mod download
WORKDIR /usr/src/gw/cmd
RUN go build -o /app
RUN rm -rf /usr/src/app
CMD ["/app"]
