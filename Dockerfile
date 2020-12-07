FROM golang:1.13.1 as bd
RUN adduser --disabled-login --gecos "" appuser
WORKDIR /github.com/layer5io/meshery-istio
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-istio .

FROM alpine
RUN apk add --no-cache \
	libc6-compat
COPY --from=bd /meshery-istio /app/
# USER appuser
WORKDIR /app
CMD ./meshery-istio
