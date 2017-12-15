# Wumpe

Wumpe is a simple auto build & deploy system which listens for GitHub or GitLab webhook events and runs a configured command, e.g. a script triggering a git pull and then a Hugo build or just a Makefile.

## Setup

Make sure [Go is correctly installed and configured](https://golang.org/doc/install) and a `$GOPATH` is set.
In the following we assume that `$GOBIN=/usr/local/bin`.

1. Get source code: `go get -u github.com/stustanet/wumpe`
2. `cd $GOPATH/src/github.com/stustanet/wumpe`
3. Build: `go install`
4. Copy the systemd unit file to `/etc/systemd/system/wumpe.service`: `cp systemd/wumpe.service /etc/systemd/system/wumpe.service`
5. `cp wumpe.toml.sample /etc/wumpe.toml` and adjust it.
6. Setup the build system user and git repos as configured in `wumpe.toml` and `wumpe.service`.
7. Setup the webhooks in your GitHub (Settings > Webhooks) or GitLab (Settings > Integrations) repo.
8. Activate Wumpe by running run `systemctl enable wumpe && systemctl start wumpe`

Wumpe should now be running.

## Wumpe updates

1. `go get -u github.com/stustanet/wumpe`
2. `systemctl restart wumpe`
