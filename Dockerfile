FROM alpine:3.2
ADD db-srv /db-srv
ENTRYPOINT [ "/db-srv" ]
