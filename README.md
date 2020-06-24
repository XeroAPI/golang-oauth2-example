# golang-oauth2-example
A basic example using golang to complete the OAuth 2 flow on Xero's API without the use of an SDK.

## Running this app

The first thing you'll need is the client ID and client secret from
[your application](https://developer.xero.com/myapps/).

Copy `config.example.yml` to a new file called `config.yml`. The two most important fields we need to change in
`config.yml` are the `client_id`, and `client_secret`.

You can run this in two ways: Docker, or natively with the go binary on your machine.

There are also some environment variables that can be set to alter the behaviour of the application:

* `DEBUG` - If set to `true`, will output more information, including auth tokens and some API response bodies.
* `APP_PORT` - Docker Compose only: Will tell Docker Compose to expose the service on the given port. Defaults to 8000.

### With Docker (recommended)

After [installing Docker](https://docs.docker.com/get-docker/), run `docker-compose up --build` in this directory.

The application should go through its initial build process, then start.

### With Go

This was written with version `1.14.3` - Your mileage may vary with other versions.

## What it Does

This uses the `oauth2` library to handle the client calls.

It also spawns a small HTTP server on `localhost:8000` in order to receive the data back from the Xero API once the user
has authorised it.
