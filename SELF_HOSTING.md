# Installing Dnote Server

This guide documents the steps for installing the Dnote server on your own machine. If you prefer Docker, please see [the Docker guide](https://github.com/dnote/dnote/blob/master/host/docker/README.md).

## Overview

Dnote server comes as a single binary file that you can simply download and run. It uses Postgres as the database.

## Installation

1. Install Postgres 11+.
2. Create a `dnote` database by running `createdb dnote`
3. Download the official Dnote server release from the [release page](https://github.com/dnote/dnote/releases).
4. Extract the archive and move the `dnote-server` executable to `/usr/local/bin`.

```bash
tar -xzf dnote-server-$version-$os.tar.gz
mv ./dnote-server /usr/local/bin
```

4. Run Dnote

```bash
GO_ENV=PRODUCTION \
OnPremise=true \
DBHost=localhost \
DBPort=5432 \
DBName=dnote \
DBUser=$user \
DBPassword=$password \
WebURL=$webURL \
SmtpHost=$SmtpHost \
SmtpPort=$SmtpPort \
SmtpUsername=$SmtpUsername \
SmtpPassword=$SmtpPassword \
DisableRegistration=false
  dnote-server start
```

Replace `$user`, `$password` with the credentials of the Postgres user that owns the `dnote` database.

Replace `$webURL` with the full URL to your server, without a trailing slash (e.g. `https://your.server`).

Replace `$SmtpHost`, `SmtpPort`, `$SmtpUsername`, `$SmtpPassword` with actual values, if you would like to receive spaced repetition through email.

Replace `DisableRegistration` to `true` if you would like to disable user registrations.

By default, dnote server will run on the port 3000.

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
ExecStart=/usr/local/bin/dnote-server start
Environment=GO_ENV=PRODUCTION
Environment=OnPremise=true
Environment=DBHost=localhost
Environment=DBPort=5432
Environment=DBName=dnote
Environment=DBUser=$DBUser
Environment=DBPassword=$DBPassword
Environment=DBSkipSSL=true
Environment=WebURL=$WebURL
Environment=SmtpHost=
Environment=SmtpPort=
Environment=SmtpUsername=
Environment=SmtpPassword=

[Install]
WantedBy=multi-user.target
```

Replace `$user`, `$WebURL`, `$DBUser`, and `$DBPassword` with the actual values.

Optionally, if you would like to send spaced repetitions throught email, populate `SmtpHost`, `SmtpPort`, `SmtpUsername`, and `SmtpPassword`.

2. Reload the change by running `sudo systemctl daemon-reload`.
3. Enable the Daemon  by running `sudo systemctl enable dnote`.`
4. Start the Daemon by running `sudo systemctl start dnote`

### Configure clients

Let's configure Dnote clients to connect to the self-hosted web API endpoint.

#### CLI

We need to modify the configuration file for the CLI. It should have been generated at `~/.dnote/dnoterc` upon running the CLI for the first time.

The following is an example configuration:

```yaml
editor: nvim
apiEndpoint: https://api.getdnote.com
```

Simply change the value for `apiEndpoint` to a full URL to the self-hosted instance, followed by '/api', and save the configuration file.

e.g.

```yaml
editor: nvim
apiEndpoint: my-dnote-server.com/api
```

#### Browser extension

Navigate into the 'Settings' tab and set the values for 'API URL', and 'Web URL'.
