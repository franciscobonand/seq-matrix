FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /seq-val

FROM alpine:3.17

COPY --from=build /seq-val /seq-val
COPY --from=build /app/.env /.env

EXPOSE 9001

ENTRYPOINT ["/seq-val"]