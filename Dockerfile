FROM alpine:3.8

ENV SENDMAIL_AMQP_URL= \
    SENDMAIL_AMQP_QNAME=sendmail \
    SENDMAIL_PROVIDER_NAME=sendgrid \
    SENDMAIL_PROVIDER_URL= \
    SENDMAIL_PROVIDER_KEY= \
    SENDMAIL_TEMPLATE_PATH= \
    SENDMAIL_LOG_LEVEL=

ADD ./sendmail /bin/
ADD ./docker-entrypoint.sh /bin/

RUN adduser -D -s /sbin/nologin sendmail

USER sendmail

ENTRYPOINT ["docker-entrypoint.sh"]