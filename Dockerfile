FROM golang as build

WORKDIR /not-gopath

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make

FROM alpine as certificates
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=build /not-gopath/firestorerestore /
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/firestorerestore"]
