FROM alpine:3.11.5

COPY ./run.sh /
COPY ./kubedebug /bin/kubedebug
RUN chmod 755 /run.sh
ENV CONTAINER_ID=1
ENV IMAGE=nicolaka/netshoot:latest
ENV RCMD=bash
EXPOSE 19675

CMD ["/run.sh"]
