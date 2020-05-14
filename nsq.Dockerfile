FROM nsqio/nsq

COPY ./entrypoint.sh .
RUN apk add curl jq bash

CMD [ "./entrypoint.sh", "nsqd", "--lookupd-tcp-address=nsqlookupd:4160", "--tls-root-ca-file=/etc/ssl/certs/root.crt", "--tls-cert=/etc/ssl/certs/dev.internal.crt", "--tls-key=/etc/ssl/certs/dev.internal.key", "--tls-required=true", "--data-path=/tmp/nsq" ]
