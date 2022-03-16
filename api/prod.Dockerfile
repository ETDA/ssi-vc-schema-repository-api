FROM ssi-registry.teda.th/ssi/ssi-core-api/core:1.0.0

ADD go.mod go.sum /app/
RUN go mod download
ADD . /app/
RUN go build -o main

FROM alpine:3.13.1
COPY --from=0 /app/main /main
CMD ./main
