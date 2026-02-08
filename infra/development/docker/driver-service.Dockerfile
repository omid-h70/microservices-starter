FROM alpine

WORKDIR /app
#build is done one my system
COPY ./build/driver-service .

EXPOSE 8081
CMD ["/app/driver-service"]