FROM golang:1.12.0-alpine3.9 as builder
RUN apk add --update git binutils upx
RUN mkdir -p /src/build
WORKDIR /src
COPY ./ ./
RUN go build -o build/imup
RUN strip --strip-unneeded build/imup
RUN upx -9 build/imup

FROM alpine:3.9
RUN apk add --update tini ca-certificates \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /app /data \
    && chmod a+w /data
COPY --from=builder --chown=root:root /src/build/imup /app/
VOLUME [ "/data" ]
WORKDIR /data
ENTRYPOINT [ "tini", "--" ]
CMD [ "/app/imup", "-listen", "0.0.0.0:3000", "-storage", "/data" ]
EXPOSE 3000