FROM alpine

COPY ./bin/admission-webhook-demo /admission-webhook-demo

ENTRYPOINT ["/admission-webhook-demo"]
