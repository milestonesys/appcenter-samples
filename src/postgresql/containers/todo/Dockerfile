FROM golang:1.24.2 AS build
COPY src /usr/src/todo
WORKDIR /usr/src/todo
RUN go build

FROM ubuntu:24.04
COPY --from=build /usr/src/todo/todo /bin/todo
CMD [ "todo" ]
