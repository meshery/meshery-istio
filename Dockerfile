FROM golang:1.13 as bd
RUN adduser --disabled-login appuser
WORKDIR /github.com/layer5io/meshery-istio
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-istio .
RUN find . -name "*.go" -type f -delete; mv istio /

FROM alpine
RUN apk --update add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=bd /meshery-istio /app/
COPY --from=bd /istio /app/istio
COPY --from=bd /etc/passwd /etc/passwd
USER appuser
WORKDIR /app
CMD ./meshery-istio
