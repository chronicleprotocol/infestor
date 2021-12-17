FROM golang:1-alpine

ARG SETZER_BRANCH="sc-413/setzer-configurable-endpoints"

RUN apk add --no-cache -X http://dl-cdn.alpinelinux.org/alpine/edge/testing \
  bash git datamash jq curl make jshon perl build-base ca-certificates

# Installing setzer
RUN git clone -b ${SETZER_BRANCH} https://github.com/makerdao/setzer-mcd.git \
  && cp -R setzer-mcd/libexec/setzer/ /root/setzer \
  && rm -rf setzer-mcd

ENV PATH="/root/setzer:${PATH}"

# TMP
ENV SMOCKER="http://172.17.0.2:8080"
ENV ORIGIN_KRAKEN_URL=$SMOCKER
ENV ORIGIN_POLONIEX_URL=$SMOCKER
ENV ORIGIN_HUOBI_URL=$SMOCKER
ENV ORIGIN_COINBASE_URL=$SMOCKER
ENV ORIGIN_BITFINEX_URL=$SMOCKER
ENV ORIGIN_BINANCE_URL=$SMOCKER

WORKDIR /app

COPY . /app

CMD ["/bin/bash"]