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

A database with the following structure:
        
        DROP TABLE posts;
        CREATE TABLE `posts` (
            `id` INTEGER PRIMARY KEY AUTOINCREMENT,
            `slug` VARCHAR(255) NULL,
            `title` VARCHAR(255) NULL,
            `description` VARCHAR(400) NULL,
            `published` DATETIME NULL,
            `body` TEXT,
            `transcript` TEXT NULL,
            `topresult` TEXT NULL,
            `click_to_tweet` VARCHAR(20)
         );
        
        
        DROP TABLE `posts_keywords_xref`;
        CREATE TABLE `posts_keywords_xref` (
          `post_id` INT,
          `keyword_id` VARCHAR(100),
          PRIMARY KEY (post_id, keyword_id)
        );
        
        
        
        DROP TABLE `youtube`;
        CREATE TABLE `youtube` (
          `id` VARCHAR(255) PRIMARY KEY,
          `post_id` INT NOT NULL,
          `body` TEXT NULL
        );
        
        
        DROP TABLE `youtube_music_xref`;
        CREATE TABLE `youtube_music_xref` (
            `youtube_id` INT,
            `music_id` VARCHAR(255)
        );





Usage
-----

Note: on running for the first time, a browser instance will be spawned and you 
will be requested to grant permission for the application to you YouTube channel.

Run the tests:

        go test

Install:

        go install

Back-up (title and body only) from YouTube to the ```backup``` directory:

        go-youtube-admin backup

Update YouTube:

        go-youtube-admin update 139

If things go wrong, delete this file:

        rm ~/.credentials/youtube-go-quickstart.json 

