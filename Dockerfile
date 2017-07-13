FROM alpine
RUN apk --no-cache add ca-certificates
ADD waw /
CMD ["/waw"]