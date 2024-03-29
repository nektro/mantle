# Mantle
![loc](https://sloc.xyz/github/nektro/mantle)
[![license](https://img.shields.io/github/license/nektro/mantle.svg)](https://github.com/nektro/mantle/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg?logo=discord)](https://discord.gg/P6Y4zQC)
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
| `--oauth2-client` | `string` | none. | OAuth2 Client Config in the form `For|ID|Secret`. |

The Identity Provider IDs can be found from the table in the [nektro/go.oauth2](https://github.com/nektro/go.oauth2#readme) documentation.

## All Flags
```
      --base string                  The path to mount all listeners on (default "/")
      --config string                 (default "~/.config/mantle/config.json")
      --dbstorage-debug-sql          Enable this flag to print all executed SQL statements.
      --dbstorage-debug-verbose      Enabled this flag to inlcude binded values in logs.
      --jwt-secret string            Private secret to sign and verify JWT auth tokens with. (default "Random")
      --max-member-count int         
      --mysql-database string        
      --mysql-password string        
      --mysql-url string             
      --mysql-user string            
      --oauth2-client stringArray    Custom client config. Pass in the form: for|id|secret
      --oauth2-default-auth string   A default auth to use when multiple appconf's are enabled.
      --port int                     The port to bind the web server to. (default 8000)
      --postgres-database string     
      --postgres-password string     
      --postgres-sslmode string       (default "verify-full")
      --postgres-url string          
      --postgres-user string         
      --redis-url string             
      --skip-translation-fetch       Enable this flag to only read native translation data.
      --theme stringArray            A CLI way to add config themes.
```

## Deployment
Pre-compiled binaries can be obtained from https://github.com/nektro/mantle/releases/latest.

Or you may run Mantle using the official Docker image. Visit https://hub.docker.com/r/nektro/mantle for the full list of tags.

```
$ docker run -p 80:8000 nektro/mantle
```

## Development
Local development is done with [Docker](https://docs.docker.com/get-docker/) and [`docker-compose`](https://docs.docker.com/compose/install/)

To launch a local instance, add the `OAUTH2_CLIENT_N` env var into `docker-compose.yml` where `N` is 1, 2, 3, etc with your OAuth2 client info and then run the following and visit http://localhost/.

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

<!-- - [![](https://mantle.trademark.cat/api/etc/badges/members_total.svg)](https://mantle.trademark.cat/) -->
<!-- - [![](https://mantle.varelus.com/api/etc/badges/members_total.svg)](https://mantle.varelus.com/) -->

Run an instance and want your site here? Contact me below!

## Contact
- hello@nektro.net
- https://twitter.com/nektro

## License
AGPL-3.0
