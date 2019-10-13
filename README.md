# Instagram Photos Resizer

A Command Line Application to resize and split photos for Instagram and upload those to Google Photos.

## TL;DR

1. Create a new Google Cloud project https://console.cloud.google.com/projectcreate
1. Enable the Photos Library API
1. Open https://console.cloud.google.com/apis/credentials
1. Create an OAuth client ID where the application type is other
1. Set the following environment variables:

```sh
export GOOGLE_CLIENT_ID=
export GOOGLE_CLIENT_SECRET=
```

```
./instagram-photos-resizer photo1.jpg photo2.jpg
```

The photos will be uploaded to your Google Photos Library. 

## What it does

Instagram limits the width of the uploaded photos to 1080px.
It is recommended to resize the image prior to uploading to avoid heavy compression.

This tool resizes the specified photo(s), splitting them into multiple photos when it detects the width is a multiple of the height.

For example, assume the photo specified is 4000px by 2500px, the tool will split the original photo into 2, both at 1080px by 1350px.

## How to build

* `make` to build the binary
* `make cross` to cross build the binary
* `make image` to build a Docker image