#! /bin/sh

function term()
{
    kill -15 $child
    wait $child
}

trap term SIGTERM

exec "`sendmail \
	--provider-url ${SENDMAIL_PROVIDER_URL} \
	--provider-key ${SENDMAIL_PROVIDER_KEY} \
	--provider-name ${SENDMAIL_PROVIDER_NAME} \
	--templates-path ${SENDMAIL_TEMPLATES_PATH} \
	--amqp-url ${SENDMAIL_AMQP_URL} \
    --amqp-qname ${SENDMAIL_AMQP_QNAME} \
	--log-level ${SENDMAIL_LOG_LEVEL} \
	`" &

child=$!
wait $child