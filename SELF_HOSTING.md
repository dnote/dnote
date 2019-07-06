# Hosting Dnote On Your Machine

This guide documents the steps for installing the Dnote server on your own machine.

**Please note that self hosted version of Dnote server is currently in beta.**

## Installation

1. Install Postgres 10+.
2. Create a `dnote` database by running `createdb dnote`
2. Download the official Dnote server release.
3. Extract the archive and move the `dnote-server` executable to `/usr/local/bin`.

```bash
tar -xzf dnote-server-$version-$os.tar.gz
chmod +x ./dnote-server/dnote-server
mv ./dnote-server/dnote-server /usr/local/bin
```

4. Run Dnote

```bash
GO_ENV=PRODUCTION \
DBHost=localhost \
DBPort=5432 \
DBName=dnote \
DBUser=$user \
DBPassword=$password \
  dnote-server start
```

Replace $user and $password with the credentials of the Postgres user that owns the `dnote` database.

By default, dnote server will run on the port 8080.

## Configuration

By now, Dnote is fully functional in your machine. The API, frontend app, and the background tasks are all in the single binary. Let's take a few more steps to configure Dnote.

### Configure Nginx

To make it accessible from the Internet, you need to configure Nginx.

1. Create a new file in `/etc/nginx/sites-enabled/dnote` with the following contents:

```
server {
	server_name my-dnote-server.com;

	location / {
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $remote_addr;
		proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:8080; 
	}
}
```

2. Reload the nginx configuration by running the following:

```
sudo service nginx reload
```

Now you can access the Dnote frontend application on `http://my-dnote-server.com`, and the API on `http://my-dnote-server.com/api`.

### Configure TLS by using LetsEncrypt

TODO

### Run Dnote As a Daemon

We can use `systemctl` to run Dnote in the background as a Daemon, and automatically start it on system reboot.

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
ExecStart=/home/$user/dnote-server start
Environment=GO_ENV=PRODUCTION DBHost=localhost DBPort=5432 DBName=dnote DBUser=$DBUser
Environment=DBPassword=$DBPassword

[Install]
WantedBy=multi-user.target
```

Replace `$user`, `$DBUser`, and `$DBPassword` with the actual values.

2. Reload the change by running `sudo systemctl daemon-reload`.
3. Enable the Daemon  by running `sudo systemctl enable dnote`.`
4. Start the Daemon by running `sudo systemctl start dnote`
