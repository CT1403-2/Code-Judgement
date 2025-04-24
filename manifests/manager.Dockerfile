FROM golang:1.24.2-bookworm AS builder

WORKDIR /app

COPY manager/go.mod manager/go.sum Makefile ./

RUN make dependencies

COPY . .

ENV PATH="/app/build/bin:$PATH"

RUN make  target=manager

FROM debian:bookworm-slim@sha256:b1211f6d19afd012477bd34fdcabb6b663d680e0f4b0537da6e6b0fd057a3ec3

COPY --from=builder /app/build/bin/manager /usr/local/bin/manager

CMD ["/usr/local/bin/manager"]
