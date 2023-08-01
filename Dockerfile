FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o application cmd/main.go

EXPOSE 3000

CMD [ "/app/application" ]