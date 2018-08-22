FROM golang as build

RUN mkdir -p /go/src/github.com/christianalexander/FirestoreRestore
WORKDIR /go/src/github.com/christianalexander/FirestoreRestore

COPY . .
RUN make

FROM alpine as certificates
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=build /go/src/github.com/christianalexander/FirestoreRestore/FirestoreRestore /firestorerestore
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/firestorerestore"]
