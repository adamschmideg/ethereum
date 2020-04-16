FROM offcode/ethspam AS builder

FROM alpine:latest AS production
RUN apk update && apk add curl
COPY --from=builder /usr/local/bin/ethspam /usr/local/bin/
COPY ./spam.sh /usr/local/bin/spam.sh
ENTRYPOINT ["spam.sh"]
