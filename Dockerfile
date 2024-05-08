FROM alpine:3.14
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -g 990 ayoconnect && \
    adduser -S -D -H -u 990 -G ayoconnect ayoconnect
USER ayoconnect
COPY ./of-bulk-disbursement /usr/local/bin/
EXPOSE 9000
CMD ["of-bulk-disbursement"]