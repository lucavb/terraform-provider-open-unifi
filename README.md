# terraform-provider-open-unifi

Terraform provider for [open-unifi](https://github.com/lucavb/open-unifi) — registry address **`registry.terraform.io/lucavb/open-unifi`**.

This repository holds the provider implementation and release artifacts. The UniFi controller (inform, adoption, admin API) lives in the **open-unifi** repo.

## Build

```sh
go build -o dist/terraform-provider-open-unifi ./cmd/tfprovider
# or
make provider-build
```

## Local development

See [examples/terraform/README.md](examples/terraform/README.md) for filesystem-mirror `dev_overrides` against a running controller.

## Acceptance tests

Black-box Terraform tests live under `acceptance/` (require `terraform` CLI and `TF_ACC=1`). They build the provider from this repo and the controller from a sibling checkout:

```sh
export OPEN_UNIFI_ROOT=../open-unifi   # default when repos are siblings
make acceptance
```

## Release

Tags matching `v*` trigger [GoReleaser](.goreleaser.yml) via `.github/workflows/release.yml`.

GitHub Actions secrets (same GPG key as other `lucavb` provider repos):

- `GPG_PRIVATE_KEY`
- `PASSPHRASE`

GoReleaser uses `GPG_FINGERPRINT` from the import-gpg step.

After a green release, publish the provider on the Terraform Registry from this repository (`terraform-provider-open-unifi`).

**Note:** Provider binaries released from `lucavb/open-unifi` (e.g. v0.1.0) are obsolete for registry publish; use releases from this repo.

## CI

See [AGENTS.md](AGENTS.md).
