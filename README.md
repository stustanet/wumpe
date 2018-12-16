# Wumpe

Wumpe is a simple auto build & deploy system which listens for GitHub or GitLab webhook events and runs a configured command, e.g. a script triggering a git pull and then a Hugo build or just a Makefile.

## Setup

Make sure [Go is correctly installed and configured](https://golang.org/doc/install) and a `$GOPATH` is set.
In the following we assume that `$GOBIN=/usr/local/bin`.

1. Get source code, build and install Wumpe: `go get -u github.com/stustanet/wumpe`
2. `cd $GOPATH/src/github.com/stustanet/wumpe`
3. `cp systemd/wumpe.service /etc/systemd/system/wumpe.service`
4. `cp wumpe.toml.sample /etc/wumpe.toml` and adjust it.
5. Setup the build system user and git repos as configured in `wumpe.toml` and `wumpe.service`.
6. Setup the webhooks in your GitHub (Settings > Webhooks) or GitLab (Settings > Integrations) repo.
7. Activate Wumpe by running run `systemctl enable --now wumpe`

Wumpe should now be running. You can check the status with `systemctl status wumpe`.

## Update Wumpe

1. `go get -u github.com/stustanet/wumpe`
2. `systemctl restart wumpe`
