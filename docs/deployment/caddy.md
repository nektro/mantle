# Deploying Mantle with Caddy

1. Download a copy of Mantle from https://github.com/nektro/mantle/releases/latest.
2. Make sure Mantle is visible from http://localhost:8000
3. Download Caddy (https://caddyserver.com/download)
4. Configure a new site context for Mantle such as:

```caddy
mantle.example.com {
    proxy / http://localhost:8000/ {
        transparent
    }
}
```

## Using HTTPS
Add the following option to your ``proxy`` block:
```caddy
header_upstream X-TLS-Enabled true
```

## Serving from an HTTP base that is not `/`
```caddy
example.com {
    proxy /mantle http://localhost:8000/ {
        transparent
    }
}
```

- The `--base` option must be sent with the exact text of the Caddy location. Ie: `./mantle --root ROOT --base /mantle/`.

## Notes:
- The leading slash at the end of the `proxy` directive is critical, particularly if you are serving Mantle from a URL that isn't `/`.
- If the exposed port is not `80` or `443`, then the `header_upstream` directive value must be `Host $host:$server_port`.
- Your OAuth2 callback URL must the full accessible location of `MANTLE/callback`.
