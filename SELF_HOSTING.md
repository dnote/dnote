# Installing Dnote Server

This guide documents the steps for installing the Dnote server on your own machine. If you prefer Docker, please see [the Docker guide](https://github.com/dnote/dnote/blob/master/host/docker/README.md).

## Overview

Dnote server comes as a single binary file that you can simply download and run. It uses SQLite as the database.

## Installation

1. Download the official Dnote server release from the [release page](https://github.com/dnote/dnote/releases).
2. Extract the archive and move the `dnote-server` executable to `/usr/local/bin`.

```bash
tar -xzf dnote-server-$version-$os.tar.gz
mv ./dnote-server /usr/local/bin
```

3. Run Dnote

```bash
dnote-server start --webUrl=$webURL
```

Replace `$webURL` with the full URL to your server, without a trailing slash (e.g. `https://your.server`).

Additional flags:
- `--port`: Server port (default: `3000`)
- `--disableRegistration`: Disable user registration (default: `false`)
- `--logLevel`: Log level: `debug`, `info`, `warn`, or `error` (default: `info`)
- `--appEnv`: environment (default: `PRODUCTION`)

You can also use environment variables: `PORT`, `WebURL`, `DisableRegistration`, `LOG_LEVEL`, `APP_ENV`.

## Configuration

By now, Dnote is fully functional in your machine. The API, frontend app, and the background tasks are all in the single binary. Let's take a few more steps to configure Dnote.

### Configure Nginx

To make it accessible from the Internet, you need to configure Nginx.

1. Install nginx.
2. Create a new file in `/etc/nginx/sites-enabled/dnote` with the following contents:

```
server {
	server_name my-dnote-server.com;

	location / {
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $remote_addr;
		proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:3000;
	}
}
```
3. Replace `my-dnote-server.com` with the URL for your server.
4. Reload the nginx configuration by running the following:

```
sudo service nginx reload
```

### Configure Apache2

1. Install Apache2 and install/enable mod_proxy.
2. Create a new file in `/etc/apache2/sites-available/dnote.conf` with the following contents:

```
<VirtualHost *:80>
    ServerName notes.example.com

    ProxyRequests Off
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:3000/ keepalive=On
    ProxyPassReverse / http://127.0.0.1:3000/
    RequestHeader set X-Forwarded-HTTPS "0"
</VirtualHost>
```

3. Enable the dnote site and restart the Apache2 service by running the following:

```
a2ensite dnote
sudo service apache2 restart
```

Now you can access the Dnote frontend application on `/`, and the API on `/api`.

### Configure TLS by using LetsEncrypt

It is recommended to use HTTPS. Obtain a certificate using LetsEncrypt and configure TLS in Nginx.

In the future versions of the Dnote Server, HTTPS will be required at all times.

### Run Dnote As a Daemon

We can use `systemd` to run Dnote in the background as a Daemon, and automatically start it on system reboot.

1. Create a new file at `/etc/systemd/system/dnote.service` with the following content:

```
[Unit]
Description=Starts the dnote server
Requires=network.target
After=network.target

[Service]
Type=simple
User=$user
Restart=always
RestartSec=3
WorkingDirectory=/home/$user
ExecStart=/usr/local/bin/dnote-server start --webUrl=$WebURL

[Install]
WantedBy=multi-user.target
```

Replace `$user` and `$WebURL` with the actual values.

By default, the database will be stored at `$XDG_DATA_HOME/dnote/server.db` (typically `~/.local/share/dnote/server.db`). To use a custom location, add `--dbPath=/path/to/database.db` to the `ExecStart` command.

2. Reload the change by running `sudo systemctl daemon-reload`.
3. Enable the Daemon  by running `sudo systemctl enable dnote`.`
4. Start the Daemon by running `sudo systemctl start dnote`

### Optional: Email Support

To enable sending emails, add the following environment variables to your configuration. But they are not required.

- `SmtpHost` - SMTP server hostname
- `SmtpPort` - SMTP server port
- `SmtpUsername` - SMTP username
- `SmtpPassword` - SMTP password

For systemd, add these as additional `Environment=` lines in `/etc/systemd/system/dnote.service`.

### Configure clients

Let's configure Dnote clients to connect to the self-hosted web API endpoint.

#### CLI

We need to modify the configuration file for the CLI. It should have been generated at `~/.config/dnote/dnoterc` upon running the CLI for the first time.

The following is an example configuration:

```yaml
editor: nvim
apiEndpoint: https://localhost:3000/api
```

Simply change the value for `apiEndpoint` to a full URL to the self-hosted instance, followed by '/api', and save the configuration file.

e.g.

```yaml
editor: nvim
apiEndpoint: my-dnote-server.com/api
```
