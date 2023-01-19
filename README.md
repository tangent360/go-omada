# go-omada

golang library for the TP-Link Omada SDN. This is work in progress so is not complete.

This is currently being tested on controller version `5.7.6` on a hardware `OC200` controller , and may not work correctly on other versions as TP-Link do not publish up to date API documentation with new releases.

# Features:
- login
- get networks
- get devices
- get clients

# Example usage
See [example/main.go](example/main.go)

# Authentication

Authentication is handled via a username and password, which you can create in the controller admin section. Permissions are not very granular with only `admin` or `reader` roles available. Currently this package only needs `reader`.

The provided [example](example/main.go) shows how to provide credentials via environment variables if that is your thing.

Once logged in the sessions appear to be valid for 14 days and you will need your application to periodically re-login if the session expires by calling `omada.Login()`.

# HTTPS Vertification

HTTPS verification is enabled by default and recommended for good security. This can be disabled by setting the environment variable `OMADA_DISABLE_HTTPS_VERIFICATION` to `true`
