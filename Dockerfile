FROM alpine:3.17

RUN apk add --update --no-cache docker-cli git

WORKDIR /var/opts/gitops

COPY start.sh ./

RUN chmod +x start.sh

ENTRYPOINT ["./start.sh"]