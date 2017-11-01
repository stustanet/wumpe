# Wumpe

Wumpe is a tiny static website builder which listens for GitLab webhooks and runs e.g. Hugo (configurable) when applicable.

## Setup

1. Checkout the git repo e.g. to `/usr/local/src/`.
2. From `/usr/local/src/wumpe` run `go build -o /usr/local/bin/wumpe`
3. Copy `/usr/local/src/wumpe/wumpe.toml.sample` to `/usr/local/src/wumpe/wumpe.toml` and adjust it.
4. Setup the build system user and git repos as configured in the `wumpe.toml`.
5. Setup the webhooks in your GitLab repo (Settings > Integrations)
4. Symlink / copy the systemd unit in `systemd/wumpe.service` to `/etc/systemd/system/wumpe.service`
5. Activate Wumpe by running run `systemctl enable wumpe && systemctl start wumpe`

Wumpe should now be running.

## Wumpe updates

1. `cd /usr/local/src/wumpe`
2. `git pull`
3. `go build -o /usr/local/bin/wumpe`
4. `systemctl restart wumpe`
