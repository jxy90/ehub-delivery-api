FROM golang:1.13 AS builder

RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /go/src/github.com/hublabs/ehub-delivery-api
ADD go.mod go.sum ./
RUN go mod download
ADD . /go/src/github.com/hublabs/ehub-delivery-api
ENV CGO_ENABLED=0
RUN go build -o ehub-delivery-api

FROM pangpanglabs/alpine-ssl
WORKDIR /go/src/github.com/hublabs/ehub-delivery-api
COPY --from=builder /go/src/github.com/hublabs/ehub-delivery-api/*.yml /go/src/github.com/hublabs/ehub-delivery-api/
COPY --from=builder /go/src/github.com/hublabs/ehub-delivery-api/ehub-delivery-api /go/src/github.com/hublabs/ehub-delivery-api/
COPY --from=builder /go/src/github.com/hublabs/ehub-delivery-api/run.sh /go/src/github.com/hublabs/ehub-delivery-api/
RUN chmod +x ./run.sh

EXPOSE 5000

CMD ["/bin/sh","./run.sh"]