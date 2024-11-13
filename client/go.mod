module github.com/brandonweeks/acme-device-attest-demo/client

go 1.22.7

toolchain go1.23.3

require (
	github.com/fxamacker/cbor/v2 v2.7.0
	github.com/google/go-attestation v0.5.1
	github.com/google/go-tpm-tools v0.4.4
	golang.org/x/crypto v0.29.0
)

require (
	github.com/google/certificate-transparency-go v1.2.2 // indirect
	github.com/google/go-tpm v0.9.1 // indirect
	github.com/google/go-tspi v0.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sys v0.27.0 // indirect
)

replace golang.org/x/crypto => github.com/brandonweeks/golang-crypto v0.0.0-20241107225453-6018723c7405

replace github.com/google/go-attestation => github.com/brandonweeks/go-attestation v0.0.0-20241111153239-62f7ad0785b8
