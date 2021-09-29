FROM ubuntu

WORKDIR /mongoshake

COPY ./ /

EXPOSE 9100
EXPOSE 9101

ENTRYPOINT [ "/docker-entrypoint.sh" ]

CMD ./collector.linux --conf config/collector.conf