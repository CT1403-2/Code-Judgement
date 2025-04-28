FROM debian:bookworm-slim

RUN apt update && apt upgrade -qqy && apt install -qqy curl jq

RUN useradd -u 1000 -m runner
WORKDIR /playground

RUN chown -R runner:runner /playground

USER runner

RUN curl -LO https://go.dev/dl/go1.24.2.linux-amd64.tar.gz && \
    echo "68097bd680839cbc9d464a0edce4f7c333975e27a90246890e9f1078c7e702ad  go1.24.2.linux-amd64.tar.gz" | sha256sum -c -

RUN rm -rf /playground/go && tar -C /playground -xzf go1.24.2.linux-amd64.tar.gz && rm go1.24.2.linux-amd64.tar.gz
ENV PATH=$PATH:/playground/go/bin

RUN mkdir -p /playground/app/
RUN printf "module main\n\ngo 1.24\n" > /playground/app/go.mod

COPY judge/scripts/run.sh /playground/run.sh
