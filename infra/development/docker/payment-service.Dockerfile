FROM alpine

WORKDIR /app
#build is done one my system
COPY ./build/payment-service .

EXPOSE 8081
CMD ["/app/payment-service"]