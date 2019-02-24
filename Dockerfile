FROM alpine:latest
COPY ./kube-monkey /kube-monkey
ENTRYPOINT /kube-monkey
