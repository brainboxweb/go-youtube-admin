YouTube Admin
=============





Requirements
------------

A ```client_secrets.json``` file with the following format. Note the value of the redirect url

```
{
  "installed": {
    "client_id": "3707DDDD4238-m5q92m9at0s0rdm9hc2pb59tc66u47fh.apps.googleusercontent.com",
    "project_id": "youtubegolang-9999",
    "auth_uri": "https:\/\/accounts.google.com\/o\/oauth2\/auth",
    "token_uri": "https:\/\/accounts.google.com\/o\/oauth2\/token",
    "auth_provider_x509_cert_url": "https:\/\/www.googleapis.com\/oauth2\/v1\/certs",
    "client_secret": "ZHhhaXat2gv3xIu3GIDDDD",
    "redirect_uris": [
      "http:\/\/localhost:8080\/oauth2callback"
    ]
  }
}
```

A ```posts.yml``` file. See ```data/posts-test.yml``` for the required structure.






Usage
-----

Note: on running for the first time, a browser instance will be spawned and you 
will be requested to grant permission for the application to you YouTube channel.

Run the tests:

```go test``

Install:

```go install```

Back-up (title and body only) from YouTube to the ```backup``` directory.

```go-youtube-admin backup```

Update YouTube with data from the ```data/posts.yml``` file

```go-youtube-admin update```
