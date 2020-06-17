# Mantle
![loc](https://sloc.xyz/github/nektro/mantle)
[![license](https://img.shields.io/github/license/nektro/mantle.svg)](https://github.com/nektro/mantle/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg?logo=discord)](https://discord.gg/P6Y4zQC)
[![paypal](https://img.shields.io/badge/donate-paypal-009cdf?logo=paypal)](https://paypal.me/nektro)
[![circleci](https://circleci.com/gh/nektro/mantle.svg?style=svg)](https://circleci.com/gh/nektro/mantle)
[![release](https://img.shields.io/github/v/release/nektro/mantle)](https://github.com/nektro/mantle/releases/latest)
[![goreportcard](https://goreportcard.com/badge/github.com/nektro/mantle)](https://goreportcard.com/report/github.com/nektro/mantle)
[![codefactor](https://www.codefactor.io/repository/github/nektro/mantle/badge)](https://www.codefactor.io/repository/github/nektro/mantle)
[![downloads](https://img.shields.io/github/downloads/nektro/mantle/total.svg)](https://github.com/nektro/mantle/releases)
[![crowdin](https://badges.crowdin.net/mantle/localized.svg)](https://crowdin.com/project/mantle)
[![issuehunt](https://img.shields.io/badge/issuehunt-mantle-38d39f)](https://issuehunt.io/r/nektro/mantle)
[![docker_pulls](https://img.shields.io/docker/pulls/nektro/mantle)](https://hub.docker.com/r/nektro/mantle)
[![docker_stars](https://img.shields.io/docker/stars/nektro/mantle)](https://hub.docker.com/r/nektro/mantle)

Easy and effective communication is the foundation of any successful team or community. That's where Mantle comes in, providing you the messaging platform that puts you in charge of both the conversation and the data.

## Getting Started
These instructions will help you get the project up and running and are required before moving on.

### Creating External Auth Credentials
In order to allow users to log in to Mantle, you will need to create an app on your Identity Provider(s) of choice. See the [nektro/go.oauth2](https://github.com/nektro/go.oauth2#readme) docs for more detailed info on this process on where to go and what data you'll need.

Here you can also fill out a picture and description that will be displayed during the authorization of users on your chosen Identity Provider. When prompted for the "Redirect URI" during the app setup process, the URL to use will be `http://mantle/callback`, replacing `mantle` with any origins you wish Mantle to be usable from, such as `example.com` or `localhost:800`.

Once you have finished the app creation process you should now have a Client ID and Client Secret. These are passed into Mantle through flags as well.

| Name | Type | Default | Description |
|------|------|---------|-------------|
| `--auth-{IDP-ID}-id` | `string` | none. | Client ID. |
| `--auth-{IDP-ID}-secret` | `string` | none. | Client Secret. |

The Identity Provider IDs can be found from the table in the [nektro/go.oauth2](https://github.com/nektro/go.oauth2#readme) documentation.

### Other Flags
These flags are optional but offer more customization and enable debugging, or assist in running in containerized environments. They also may all be passed as environment variables in the form that `--oauth2-default-auth` may also be passed in with the `OAUTH2_DEFAULT_AUTH` variable.

| Name | Type | Default | Description |
|------|------|---------|-------------|
| `--dbstorage-debug-sql` | `bool` | `false` | Enable this flag to print all executed SQL statements. |
| `--oauth2-default-auth` | `string` | none | Use this if you'd like to have a default auth when multiple are enabled. |
| `--port` | `int` | `8000` | The port to bind the web server to. |
| `--jwt-secret` | `string` | random | HMAC JWT signing secret. |
| `--redis-url` | `string` | none | Host of Redis instance to use. |
| `--max-member-count` | `int`| none | The maximum number of users that may be a "member" of the server at one time. Overrides setting in "Overview". Not retroactive. |

## Deployment
Pre-compiled binaries can be obtained from https://github.com/nektro/mantle/releases/latest.

Or you may run Mantle using the official Docker image. Visit https://hub.docker.com/r/nektro/mantle for the full list of tags.

```
$ docker run -p 80:8000 nektro/mantle
```

## Development
Local development is done with [Docker](https://docs.docker.com/get-docker/) and [`docker-compose`](https://docs.docker.com/compose/install/)

To launch a local instance, edit `./data/config.json` with your OAuth2 client info and then run the following and visit http://localhost/

```
$ docker-compose up
```

## Built With
- https://github.com/PuerkitoBio/goquery
- https://github.com/aymerick/raymond
- https://github.com/dgrijalva/jwt-go
- https://github.com/gorilla/mux
- https://github.com/gorilla/websocket
- https://github.com/mattn/go-sqlite3
- https://github.com/mitchellh/go-homedir
- https://github.com/nektro/go-util
- https://github.com/nektro/go.dbstorage
- https://github.com/nektro/go.etc
- https://github.com/nektro/go.oauth2
- https://github.com/nektro/go.sdrie
- https://github.com/oklog/ulid
- https://github.com/prometheus/client_golang
- https://github.com/rakyll/statik
- https://github.com/spf13/pflag
- https://github.com/valyala/fastjson

## Contributing
[![issues](https://img.shields.io/github/issues/nektro/mantle.svg)](https://github.com/nektro/mantle/issues)

We listen to issues all the time right here on GitHub. Labels are extensively to show the progress through the fixing process. Question issues are okay but make sure to close the issue when it has been answered! Off-topic and '+1' comments will be deleted. Please use post/comment reactions for this purpose.

## Public Instances
Want to try out a live server before installing your own, or get involved in an already existing community? Try these servers below:

- [![](https://mantle.trademark.cat/api/etc/badges/members_total.svg)](https://mantle.trademark.cat/)
- [![](https://mantle.varelus.com/api/etc/badges/members_total.svg)](https://mantle.varelus.com/)

Run an instance and want your site here? Contact me below!

## Contact
- hello@nektro.net
- Meghan#2032 on discordapp.com
- https://twitter.com/nektro

## License
Apache 2.0
