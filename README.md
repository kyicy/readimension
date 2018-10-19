# [readimension](https://www.readimension.com)
readimension is a `epub` web service provides both ***file management*** and ***browser reader***.

## Features
- Book format: epub
- File Explorer like File Management
- Responsive web interface
- Web based epub reader([satorumurmur/bibi](https://github.com/satorumurmur/bibi))


## Usage
Download the binary file from [release page](https://github.com/kyicy/readimension/releases)

or 
`go get -v github.com/kyicy/readimension`

Switch to an ***working directory*** where `readimension` will save data.

Create a configuration file, `config.json`
``` json
{
    "production": {
        "addr": "0.0.0.0",
        "port": "10086",
        "session_secret": "",
        "emails": ["example@example.com"],
        "google_analytics": ""
    },
    "development": {
        "addr": "0.0.0.0",
        "port": "10086",
        "session_secret": "",
        "emails": ["example@example.com"],
        "google_analytics": ""
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