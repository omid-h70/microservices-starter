FROM alpine

WORKDIR /app
#build is done one my system
COPY ./build/api-gateway .

EXPOSE 8081
CMD ["/app/api-gateway"]