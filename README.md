# golang-oauth2-example
A basic example using golang to complete the OAuth 2 flow on Xero's API without the use of an SDK.

## What it Does

This uses the `oauth2` library to handle the client calls.

It also spawns a small HTTP server on `localhost:8000` in order to receive the data back from the Xero API once the user
has authorised it.
