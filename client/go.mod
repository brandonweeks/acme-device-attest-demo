module github.com/brandonweeks/acme-device-attest-demo/client

go 1.18

require (
	github.com/fxamacker/cbor/v2 v2.7.0
	github.com/google/go-attestation v0.4.4-0.20220404204839-8820d49b18d9
	github.com/google/go-tpm-tools v0.3.8
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
)

require (
	github.com/google/certificate-transparency-go v1.1.2 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/google/go-tspi v0.2.1-0.20190423175329-115dea689aad // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace golang.org/x/crypto => github.com/brandonweeks/golang-crypto v0.0.0-20220601020110-5663e12aa0bb

replace github.com/google/go-attestation => github.com/brandonweeks/go-attestation v0.0.0-20220602235615-164122a1d59b
