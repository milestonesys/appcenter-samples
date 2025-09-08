FROM golang:1.24.2 AS build
COPY src /usr/src/kafka
WORKDIR /usr/src/kafka
RUN go build

FROM ubuntu:24.04
COPY --from=build /usr/src/kafka/producer /bin
CMD [ "producer" ]
