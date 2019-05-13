FROM golang:1.12.0-alpine3.9 as builder
RUN apk add --update git binutils upx
RUN mkdir -p /src/build
WORKDIR /src
COPY ./ ./
RUN chmod 0755 thumb100.sh
RUN go build -o build/imup
RUN strip --strip-unneeded build/imup
RUN upx -9 build/imup

FROM alpine:3.9
RUN apk add --update tini ca-certificates imagemagick6 \
    && rm -rf /var/cache/apk/* \
    && ln -s /usr/bin/convert-6 /usr/local/bin/convert \
    && mkdir -p /data && chmod a+w /data
COPY --from=builder --chown=root:root /src/build/imup /usr/local/sbin/imup
COPY --from=builder --chown=root:root /src/thumb100.sh /usr/local/bin/thumb100.sh
VOLUME [ "/data" ]
WORKDIR /data
ENTRYPOINT [ "tini", "--" ]
CMD [ "/usr/local/sbin/imup", "-listen", "0.0.0.0:3000"]
    # "-storage", "/data" "-thumbcmd", "/usr/local/bin/thumb100.sh"]
EXPOSE 3000