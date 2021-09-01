FROM alpine:3.11.5

COPY ./run.sh /
COPY ./kubedebug /bin/kubedebug
RUN chmod 777 /bin/kubedebug
RUN chmod 755 /run.sh
ENV CONTAINERID=1
EXPOSE 19675

CMD ["/run.sh"]
