FROM golang:latest as builder
COPY . /srv/xmstep
ENV CGO_ENABLED=0 
# ENV GOPROXY=https://goproxy.cn,direct
RUN cd /srv/xmstep &&\
    go build -ldflags '-s -w' . &&\
    ls -l

#
FROM alpine:latest
LABEL source.url="https://github.com/fimreal/xmstep"
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata ca-certificates &&\
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&\
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /srv/xmstep/xmstep /xmstep

ENTRYPOINT [ "/xmstep" ]

ENV PORT=3000