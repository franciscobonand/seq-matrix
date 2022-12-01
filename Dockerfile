FROM golang:1.19-alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /seq-val

FROM build

WORKDIR /app

EXPOSE 9001

ENTRYPOINT ["/seq-val"]