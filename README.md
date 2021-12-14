# readimension
readimension is a `epub` web service provides both ***file management*** and ***browser reader***.

## Features
- Book format: epub
- File Explorer like File Management
- Responsive web interface
- Web based epub reader([satorumurmur/bibi](https://github.com/satorumurmur/bibi))


## Usage
`go install github.com/kyicy/readimension@latest`

Switch to an ***working directory*** where `readimension` will save data.

Create a configuration file, `config.json`
``` json
{
    "production": {
        "addr": "127.0.0.1",
        "port": "10086",
        "serve_static": false,
        "session_secret": "session_secret",
        "emails": ["example@example.com"],
        "google_analytics": "UA-****-*"
    },
    "development": {
        "addr": "127.0.0.1",
        "port": "10086",
        "serve_static": true,
        "session_secret": "session_secret",
        "emails": ["example@example.com"],
        "google_analytics": "UA-****-*"
    }
}
```
Then start the server
``` sh
readimension --env development --path .
```

`addr` and `port` defines which `ip` and `port` the web service shall listen to.
`emails` contains an array of emails are allowed to register users.

`readimension` will generate three folders {`uploads`, `covers`, `books`} and one database file `readimension.db`.

In production environment, it's preferred to set `serve_static` to false and set up a `nginx` instance to serve static files (`covers` and `books`).

Then, just enjoy reading.

## Screenshots

### Desktop
| ![](screenshots/pc_eva.png) |
| --- |
| ![](screenshots/pc_1.jpg) |
| ![](screenshots/pc_2.jpg) |

### Mobile
| ![](screenshots/mobile_eva.jpg) | 
| --- |
|![](screenshots/mobile_opm.jpg) |

| ![](screenshots/mobile_1.jpg) | ![](screenshots/mobile_2.jpg) |
| --- |  --- |
| ![](screenshots/mobile_3.jpg) | ![](screenshots/mobile_4.jpg) |
