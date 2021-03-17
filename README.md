Run once: `go get`

Start the echo (backend) service: `cd echo && go run .`. Leave it running
Start the auth service: `cd auth && go run .`. Leave it running
Start envoy: `docker run --rm -v=$PWD/envoy/:/envoy -p=8080:8080 envoyproxy/envoy-alpine:v1.17-latest -c /envoy/envoy-config.yaml`

Valid auth request:
```
curl --location --request GET 'localhost:8080' --header 'Authorization: my-secret'
```

Change or remove the secret to not be authorized