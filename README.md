# acme-device-attest-demo
This repository contains hosted and local demonstrations of the [draft-bweeks-acme-device-attest](https://brandonweeks.github.io/draft-bweeks-acme-device-attest/draft-bweeks-acme-device-attest.html) specification using a Trusted Platform Module.

The certificate authority is built using a [fork of `step-ca`](https://github.com/brandonweeks/step-ca/tree/acme-device-attest), an open source Go certificate authority that implements the ACME protocol. The client is built using [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20220408190544-5352b0902921/acme) and [google/go-attestation](https://github.com/google/go-attestation).

## Instructions
### Hosted
A hosted instance of the certificate authority is available at `ca.attestation.dev`. To get an ephemeral Cloud Shell environment containing this repository and virtualized Trusted Platform Module (TPM), click the button below.

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://shell.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fbrandonweeks%2Facme-device-attest-demo&cloudshell_print=cloudshell_instructions.txt&cloudshell_open_in_editor=client.go&cloudshell_workspace=client)

Then you can run `go run client.go` from within the Cloud Shell to request a certificate containing the attested TPM identity of the Cloud Shell instance.

### Local
- `cd ca/`
- `docker build -t step-ca .`
- `docker run -it step-ca`

In another shell:
- `cd client/`
- `go run client.go -ca_address="http://localhost:8080"`