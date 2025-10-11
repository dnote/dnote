# Self-Hosting Dnote Server

For Docker installation, see [the Docker guide](https://github.com/dnote/dnote/blob/master/host/docker/README.md).

## Quick Start

Download from [releases](https://github.com/dnote/dnote/releases), extract, and run:

```bash
tar -xzf dnote-server-$version-$os.tar.gz
mv ./dnote-server /usr/local/bin
dnote-server start --webUrl=https://your.server
```

You're up and running. Database: `~/.local/share/dnote/server.db` (customize with `--dbPath`). Run `dnote-server start --help` for options.

Set `apiEndpoint: https://your.server/api` in `~/.config/dnote/dnoterc` to connect your CLI to the server.

## Optional guide

### Nginx

Create `/etc/nginx/sites-enabled/dnote`:

```
server {
	server_name my-dnote-server.com;

	location / {
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $remote_addr;
		proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:3001;
	}
}
```

Replace `my-dnote-server.com` with your domain, then reload:

```bash
sudo service nginx reload
```

### Apache2

Enable `mod_proxy`, then create `/etc/apache2/sites-available/dnote.conf`:

```
<VirtualHost *:80>
    ServerName notes.example.com

    ProxyRequests Off
    ProxyPreserveHost On
    ProxyPass / http://127.0.0.1:3001/ keepalive=On
    ProxyPassReverse / http://127.0.0.1:3001/
    RequestHeader set X-Forwarded-HTTPS "0"
</VirtualHost>
```

Enable and restart:

```bash
a2ensite dnote
sudo service apache2 restart
```

### TLS

Use LetsEncrypt to obtain a certificate and configure HTTPS in your reverse proxy.

### systemd Daemon

Create `/etc/systemd/system/dnote.service`:

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

Replace `$user` and `$WebURL`. Add `--dbPath` to `ExecStart` if you want a custom database location.

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable dnote
sudo systemctl start dnote
```

### Email Support

If you want emails, add these environment variables:

- `SmtpHost` - SMTP hostname
- `SmtpPort` - SMTP port
- `SmtpUsername` - SMTP username
- `SmtpPassword` - SMTP password

For systemd, add as `Environment=` lines in the service file.
