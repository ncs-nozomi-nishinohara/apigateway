FROM --platform=$BUILDPLATFORM ncsnozominishinohara/golang:1.13.6-alpine3.11  as build


FROM --platform=$BUILDPLATFORM nginx:1.17.8-alpine
ENV NGINX_CONF=/etc/nginx/conf.d/default.conf
ENV SETTING_FILE_NAME=/app/config/config.yaml
COPY --from=build /var/app/golang/app /app/
COPY --from=build /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY watch.sh /usr/local/bin/watch.sh
COPY supervisord.conf /etc/supervisord.conf
COPY start.sh /usr/local/bin/start.sh
COPY templates/ /app/templates/
RUN chmod +x /usr/local/bin/watch.sh \
    && chmod +x /app/app \
    && chmod +x /usr/local/bin/start.sh \
    && apk add --no-cache supervisor util-linux openssl \
    && rm  -rf /tmp/* /var/cache/apk/*

ENTRYPOINT [ "start.sh" ]

CMD [ "supervisord" ]