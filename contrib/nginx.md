### Run behind nginx

If you already have an existing nginx webserver running, you can proxy your inlets server through nginx. Simply create a new site configuration file. This example assumes inlets is running port 8000, the default port, but you can use any port number you'd like via `--port=` when starting the server:

```
server {
    listen 80;

    # Replace *.inlets.example.com with your own wildcard domain.
    server_name *.inlets.example.com;

    #Inlets proxy
    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_http_version 1.1;
        proxy_read_timeout 120s;
        proxy_connect_timeout 120s;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```
You can even secure your connection with SSL via your own certificate, or via Letsencrypt!
