FROM golang:1.17-bullseye AS Build

RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build -o main -mod=mod .

FROM scratch

COPY --from=Build /build/main /app/
WORKDIR /app

CMD ["./main"]