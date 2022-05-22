FROM golang:1.18-alpine AS build

RUN mkdir /src
ADD . /src
WORKDIR /src
ENV CGO_ENABLED=0
RUN go build ./cmd/tstr


FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=build /src/tstr /tstr

EXPOSE 9090 9000
USER nonroot:nonroot
ENTRYPOINT ["/tstr"]
