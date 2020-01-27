# Mantle
![loc](https://sloc.xyz/github/nektro/mantle)
[![license](https://img.shields.io/github/license/nektro/mantle.svg)](https://github.com/nektro/mantle/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg)](https://discord.gg/P6Y4zQC)
[![paypal](https://img.shields.io/badge/donate-paypal-009cdf)](https://paypal.me/nektro)
[![circleci](https://circleci.com/gh/nektro/mantle.svg?style=svg)](https://circleci.com/gh/nektro/mantle)
[![release](https://img.shields.io/github/v/release/nektro/mantle)](https://github.com/nektro/mantle/releases/latest)
[![goreportcard](https://goreportcard.com/badge/github.com/nektro/mantle)](https://goreportcard.com/report/github.com/nektro/mantle)
[![codefactor](https://www.codefactor.io/repository/github/nektro/skarn/badge)](https://www.codefactor.io/repository/github/nektro/skarn)

Easy and effective communication is the foundation of any successful team or community. That's where Mantle comes in, providing you the messaging and voice platform that puts you in charge of both the conversation and the data.

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

## Development

### Prerequisites
- The Go Language 1.12+ (https://golang.org/dl/)
- GCC on your PATH (for the https://github.com/mattn/go-sqlite3 installation)

### Installing
Run
```
$ go get -u -v github.com/nektro/mantle
```
and then make your way to `$GOPATH/src/github.com/nektro/mantle/`.

Once there, run:
```
$ ./start.sh
```

## Deployment
Pre-compiled binaries can be obtained from https://github.com/nektro/mantle/releases/latest.

Or you can build from source:
```
$ ./build_all.sh
```

## Built With
- http://github.com/gorilla/sessions
- http://github.com/gorilla/websocket
- http://github.com/nektro/go-util
- http://github.com/nektro/go.etc
- http://github.com/nektro/go.oauth2
- http://github.com/satori/go.uuid
- http://github.com/spf13/pflag

## Contributing
[![issues](https://img.shields.io/github/issues/nektro/mantle.svg)](https://github.com/nektro/mantle/issues)
[![pulls](https://img.shields.io/github/issues-pr/nektro/mantle.svg)](https://github.com/nektro/mantle/pulls)

We listen to issues all the time right here on GitHub. Labels are extensively to show the progress through the fixing process. Question issues are okay but make sure to close the issue when it has been answered! Off-topic and '+1' comments will be deleted. Please use post/comment reactions for this purpose.

When making a pull request, please have it be associated with an issue and make a comment on the issue saying that you're working on it so everyone else knows what's going on :D

## Contact
- hello@nektro.net
- Meghan#2032 on discordapp.com
- https://twitter.com/nektro

## License
Apache 2.0
