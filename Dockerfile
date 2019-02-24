FROM alpine
COPY ./kube-monkey /kube-monkey
ENTRYPOINT /kube-monkey
