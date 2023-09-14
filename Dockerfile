FROM golang:1.19-alpine

WORKDIR /app

COPY . .

EXPOSE 8042

CMD ["go","run","phony.go"]