FROM postgres:18-alpine
# Apply distribution fixes even before the upstream image is rebuilt. The
# entrypoint needs only a user switch; su-exec avoids the bundled Go runtime.
RUN apk upgrade --no-cache && apk add --no-cache su-exec && rm /usr/local/bin/gosu
COPY docker/gosu.sh /usr/local/bin/gosu
RUN chmod 0755 /usr/local/bin/gosu
