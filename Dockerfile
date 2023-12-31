FROM golang:latest AS build

WORKDIR /go/src/app
COPY . .
ENV CGO_ENABLED=0

RUN make build

FROM scratch
USER 1000:1000
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/bin/htpasswd-operator /htpasswd-operator


ENTRYPOINT [ "/htpasswd-operator", "watch" ]