# Prometheus Metrics

Metrics are not visible to users by default as it exposes info about your server that can be hidden through the application. Because of this the Prometheus metrics API endpoint, `/metrics`, may only be accessed if you are either logged in as the server owner or pass a bearer token through the `Authorization` HTTP header.

You can see an example of this in the default [`prometheus.yml`](./../data/prometheus/prometheus.yml) packaged with the source for development. However the `<BEARER_TOKEN>` value must be replaced.

The value for this token is currently only available to server owners by visiting the `/api/about?all=1` endpoint.

You can learn more about Prometheus here: https://prometheus.io/.
