FROM golang:1.20 AS builder

WORKDIR /go/src/
COPY ./ .
RUN export GOPROXY=https://proxy.golang.com.cn,direct && go build -o prometheus-webhook-wechat prometheus-webhook-wechat/cmd/

FROM golang:1.20
RUN apt-get install  bash -y
COPY --from=builder /go/src/prometheus-webhook-wechat /opt/prometheus-webhook-wechat
COPY config.yml /opt/
COPY template.tmpl /opt/
WORKDIR  /opt/

CMD ["./prometheus-webhook-wechat"]