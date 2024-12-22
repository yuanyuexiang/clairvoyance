FROM registry.cn-zhangjiakou.aliyuncs.com/matrix-net/ubuntu:24.04
LABEL maintainer "matrix@126.com"

ENV TZ='Asia/Shanghai'

WORKDIR /matrix

RUN ls -l
COPY clairvoyance ./
RUN chmod +x ./clairvoyance
CMD ["/matrix/clairvoyance"]