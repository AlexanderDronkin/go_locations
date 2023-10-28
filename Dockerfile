FROM alpine:latest
ADD ./src/locations/locations /
WORKDIR /
CMD [ "/locations" ]