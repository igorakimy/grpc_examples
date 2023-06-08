FROM golang AS build

ENV location app
ENV HOST productinfo
ENV PORT 80

ADD ./client ${location}/client
ADD ./protobuf ${location}/protobuf

COPY go.mod go.sum ${location}/

WORKDIR ${location}

RUN go get -d ./...
RUN go install ./...

WORKDIR client

RUN CGO_ENABLED=0 go build -o /bin/grpc-productinfo-client

FROM scratch
COPY --from=build /bin/grpc-productinfo-client /bin/grpc-productinfo-client
ENTRYPOINT ["/bin/grpc-productinfo-client"]
EXPOSE 50502

