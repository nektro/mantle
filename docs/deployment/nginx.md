# Deploying Mantle with Nginx

1. Download a copy of Mantle from https://github.com/nektro/mantle/releases/latest.
2. Make sure Mantle is visible from http://localhost:8000
3. Download Nginx (https://nginx.org/en/download.html)
4. Configure a `location` context for Mantle such as:

```
location / {
    proxy_pass http://localhost:8000/;
    proxy_set_header Host $host;
}
location /ws {
    proxy_pass http://localhost:8000;
    proxy_set_header Host $host;
    proxy_set_header Connection "Upgrade";
    proxy_set_header Upgrade $http_upgrade;
}
```

## Using HTTPS
Modify your `/` location block and add

```
proxy_set_header X-TLS-Enabled true;
```

## Serving from an HTTP base that is not `/`
- The `--base` option must be sent with the exact text of the nginx location. Ie: `./mantle --root ROOT --base /mantle/`.

## Notes
- The leading slash at the end of `proxy_pass` is critical, particularly if you are serving Mantle from a `location` that isn't `/`.
- If the exposed port is not `80` or `443`, then the `proxy_set_header` value must be `Host $host:$server_port`.
- Your OAuth2 callback URL must the full accessible location of `MANTLE/callback`.
