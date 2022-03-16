FROM ssi-registry.teda.th/ssi/ssi-core-api/core:1.0.0

ADD go.mod go.sum /app/
RUN go mod download
ADD . /app/
CMD fresh -c other_runner.conf
