FROM alpine:latest

LABEL maintainer="Musker.Chao <aery_mzc9123@163.com>"

WORKDIR /opt/apiTools

ADD dist.tar.gz .

RUN mv dist/* . && rm -rf dist && chmod +x apiTools

EXPOSE 8091

ENTRYPOINT ["./apiTools", "run"]

CMD ["all"]
