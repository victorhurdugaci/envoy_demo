Run once: `go get`

1. Start the echo (backend) service: `cd echo && go run .`. Leave it running
1. Start the auth service: `cd auth && go run .`. Leave it running
1. Start redis `docker run -d -p=6379:6379 redis`
1. Start the rate limit service: `cd rate-limit && go run .` Leave it running
1. Update the IP address in the endpoint discovery service (EDS): ``discovery/main.go` (there might be and easier way...)
1. Start the EDS: `cd discovery && go run .`
1. Start envoy: `docker run --rm -v=$PWD/envoy/:/envoy -p=8080:8080 envoyproxy/envoy-alpine:v1.17-latest -c /envoy/envoy-config.yaml`

Valid auth request:
```
curl --location --request GET 'localhost:8080' --header 'Authorization: my-secret'
```

Change or remove the secret to not be authorized
