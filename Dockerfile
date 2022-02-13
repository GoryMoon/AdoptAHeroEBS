
# -- Stage 1 - Compile the app -- #
FROM golang:1.17-alpine as builder
WORKDIR /app

COPY . .
RUN go build -o bin/server cmd/server/main.go

# -- Stage 2 - Create final image -- #
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/server /usr/local/bin/
VOLUME /run/db
ENV DB_LOC=/run/db
CMD ["server"]
