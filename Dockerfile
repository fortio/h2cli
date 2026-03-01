FROM scratch
COPY h2cli /usr/bin/h2cli
ENV HOME=/home/user
ENTRYPOINT ["/usr/bin/h2cli"]
