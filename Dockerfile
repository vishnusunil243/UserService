#Build stage
FROM golang:1.21-bullseye AS builder

RUN apt-get update

WORKDIR /user

COPY . .
 
COPY .env .env

RUN go mod download

RUN go build -o ./out/dist ./cmd

#production stage

FROM busybox

RUN mkdir -p /user/out/dist

COPY --from=builder /user/out/dist /user/out/dist

COPY --from=builder /user/.env /user/out

WORKDIR /user/out/dist

EXPOSE 8082

CMD ["./dist"]