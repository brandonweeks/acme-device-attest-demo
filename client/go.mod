module github.com/brandonweeks/smolclient

go 1.18

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/fxamacker/cbor/v2 v2.4.0
	github.com/google/go-attestation v0.4.3
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
)

require (
	github.com/google/certificate-transparency-go v1.1.1 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/google/go-tspi v0.2.1-0.20190423175329-115dea689aad // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/sys v0.0.0-20210629170331-7dc0b73dc9fb // indirect
)

replace golang.org/x/crypto => github.com/brandonweeks/golang-crypto v0.0.0-20220408194850-b26392bd189e

replace github.com/google/go-attestation => github.com/brandonweeks/go-attestation v0.0.0-20220408183046-b33219625f8c
