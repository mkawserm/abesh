# CREATING SELF SIGNED CERTIFICATE

```
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

or

# Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
# https://safecurves.cr.yp.to/
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key


openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

```

# RUN ABESH HTTP SERVER WITH SELF SINGED CERTIFICATE

# manifest.yaml
```

version: "1"

capabilities:
  - contract_id: "abesh:ex_echo"
    values:
      key1: "test1"
      key2: "test2"
  - contract_id: "abesh:httpserver"
    values:
      host: "0.0.0.0"
      port: "9090"
      default_request_timeout: "5s"
      cert_file: "/Volumes/BoxZ/CloudDriveZ/domain/github.com/mkawserm/abesh/server.crt"
      key_file: "/Volumes/BoxZ/CloudDriveZ/domain/github.com/mkawserm/abesh/server.key"

triggers:
  - trigger: "abesh:httpserver"
    trigger_values:
      method: "GET"
      path: "/default"
    service: "abesh:ex_echo"

start:
  - "abesh:httpserver"

```

# RUN
```
go run main/default/main.go run --manifest default.yaml

```