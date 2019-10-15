FROM golang:1.13.1 as bd
RUN adduser --disabled-login --gecos "" appuser
WORKDIR /github.com/layer5io/meshery-istio
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-istio .
RUN find . -name "*.go" -type f -delete; mv istio /
RUN wget -O /istio.tar.gz https://github.com/istio/istio/releases/download/1.3.0/istio-1.3.0-linux.tar.gz

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=bd /meshery-istio /app/
COPY --from=bd /istio /app/istio
COPY --from=bd /istio.tar.gz /app/
COPY --from=bd /etc/passwd /etc/passwd
ENV ISTIO_VERSION=istio-1.3.0
USER appuser
WORKDIR /app
CMD ./meshery-istio
