FROM golang:1.16 AS build

WORKDIR /barnbridge

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM scratch
COPY --from=build /barnbridge/farmingsimulator /farmingsimulator
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/farmingsimulator", "run", "--config=/config/config.yml"]
