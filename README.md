# Mantle
![loc](https://tokei.rs/b1/github/nektro/mantle)
[![license](https://img.shields.io/github/license/nektro/mantle.svg)](https://github.com/nektro/mantle/blob/master/LICENSE)
[![discord](https://img.shields.io/discord/551971034593755159.svg)](https://discord.gg/P6Y4zQC)
[![paypal](https://img.shields.io/badge/donate-paypal-009cdf)](https://paypal.me/nektro)
[![circleci](https://circleci.com/gh/nektro/mantle.svg?style=svg)](https://circleci.com/gh/nektro/mantle)
[![goreportcard](https://goreportcard.com/badge/github.com/nektro/mantle)](https://goreportcard.com/report/github.com/nektro/mantle)
[![astronomer](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Fastronomer.ullaakut.eu%2Fshields%3Fowner%3Dnektro%26name%3Dmantle)](https://github.com/Ullaakut/astronomer)
Easy and effective communication is the foundation of any successful team or community. That's where Mantle comes in, providing you the messaging and voice platform that puts you in charge of both the conversation and the data.

## Getting Started
These instructions will help you get the project up and running and are required before moving on.

### Creating External Auth Credentials
In order to allow users to log in to Mantle, you will need to create an app on your Identity Provider(s) of choice. See the [nektro/go.oauth2](https://github.com/nektro/go.oauth2#readme) docs for more detailed info on this process on where to go and what data you'll need.

Here you can also fill out a picture and description that will be displayed during the authorization of users on your chosen Identity Provider. When prompted for the "Redirect URI" during the app setup process, the URL to use will be `http://mantle/callback`, replacing `mantle` with any origins you wish Mantle to be usable from, such as `example.com` or `localhost:800`.

Once you have finished the app creation process and obtained the Client ID and Client Secret, create a folder in your home directory at the path of `~/.config/mantle/`. All of Mantle's config and local save files will go here. This directory will be referred to as `.mantle` going forward.

In the `.mantle` folder make a `config.json` file and put the following data inside, replacing `AUTH` with whichever Identity Provider you chose, such as `discord`, `reddit`, etc. And `CLIENT_ID` and `CLIENT_SECRET` with their respective values. Do not worry, this folder will remain entirely private, even to users with full access.

The current config version is `1`. See [docs/config](./docs/config/) for more info on setting up this info to be read by Mantle.

Once fully configured, you're ready to obtain a binary to run Mantle from either the [Development](#development) or [Deployment](#deployment) sections depending on your needs.

## Development

### Prerequisites
- A directory you wish to proxy through Mantle
- The Go Language 1.12+ (https://golang.org/dl/)

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
