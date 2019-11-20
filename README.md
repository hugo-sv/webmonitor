# webmonitor

Webmonitor is a CLI web monitoring and alerting tool written in Go.

## Usage

### Build the project

To download the project, run

```
go get "github.com/hugo-sv/webmonitor"
```

From the project directory, run

```shell
go install
```

### Run the project

#### Command Line Interface

Before running the project, you should use a full screen console window.

To run the project, if the gopath is setup, you can simply run the following command :

```
webmonitor [options] [JSON path]
```

Otherwise, `webmonitor` can be replaced by `go run main.go` from the project directory.

```
go run main.go [options] [JSON path]
```

#### Options

```
-ui BOOL (default : true)
    Whether or not to display the UI
```

`JSON path` Is the relative path to the configuration file. Some example configuration paths are located in the `data` folder.

#### Examples

From the project directory :

```shell
webmonitor "data/test1.json"
```

Or, without UI.

```shell
webmonitor -ui=false "data/test1.json"
```

You can as well use custom JSON files. The format is

```json
{
  "timeout": 5,
  "websites": [
    {
      "url": "https://google.com",
      "interval": 2
    }
  ]
}
```

With `interval` being the website's interval check, in seconds, and `timeout` being the timeout limit for get requests, in seconds.

#### User Interface

With the UI, you can press :

- **q** to quit
- **up** and **down** to scroll through alerts
- **s** to switch the statistics timeframe
- Any website ID's key, to view it details

#### Usage

In the `server` folder, there is a go script that can be built and run in another window by using

```shell
cd server
go build testServer.go
testServer
```

It will serve a local server answering to request to http://localhost:8080/. The `test2.json` contain all the relevant routes :

- **up** : Always return 200 status code
- **down** : Always return 500 status code
- **random** : Returns 200 status code with a 80% chance, 500 status code otherwise
- **alert** : Return 200 status code for 120sec then 500 status code for 42sec, making an availability varying from 60% to 100% every two minutes

The API `https://httpstat.us/{statusCode}?sleep={sleepTime}` can be use to test any response code and response time.

Disabling internet connection may as well emulate a website going down.

## Project Architecture

### Main

In the `main` function, using go channels and tickers, operations are executed as they go. Usually :

- Each website, at each of their `interval` seconds, are requested. A response time and a status code will be returned later.
- Each time a response time and a status code is returned, it is processed. If, in a **2 min** timeframe, an alert is triggered, it is added in the UI.
- Every **10 sec**, the stats view is refreshed if the user is looking at a **10 min** timeframe.
- Every **1 min**, the stats view is refreshed if the user is looking at a **1h** timeframe.
- Every time a UI input is detected, the associated action is executed.

### CLI

The `cli` module process the flags from the command executed, parse and check the JSON file.

### Display

The `display` module handles every UI related actions :

- Generating the layouts
- Updating the panels
- String formating in the `format.go` script

### Monitor

The `monitor` module handles the HTTP get request, and compute a response time.

It launches and stop goroutines that periodically fetch and send back data.

### Statistics

The `statistic` module compute the main statistics form the records of status code and response time in the considered timeframe :

- **Max** : Maximal response time
- **Avg** : Average response time
- **Availability** : Percent of successful requests (Status code 200)

## Alerting logic test

Run the `testServer` script.

```shell
cd server
testServer
```

Run the webmonitor app with `test2.json`

```shell
webmonitor "test2.json"
```

As the program starts, one alerts should be raised :

- localhost:8080/down route, not returning a 200 status code

As the program goes :

- localhost:8080/random will randomly go up and down.
- localhost:8080/alert will go up and down every two minutes

Stopping the server will also trigger 404 status error, as it is unreachable.

## Notes

### Possible improvements

With more time on this project, I would have loved to work on the following improvements.

#### Computing Max Response Time

If needed, it might be possible to improve the time performances in a trade-off with space complexity using a segment tree.

#### Tickers start time

Tickers currently start too late in this project. With a 10 sec interval, the user has to wait for 10 sec before the websites are requested.

#### Alert system

The alert system could be upgraded, so that active alerts are highlighted, with data on their duration and an updated availability.

#### Unit tests

With more time, a test driven development could have been followed.
A few tests were implemented for the `monitoring` module.

#### Responsiveness

Displaying the program require a full screen window. Some screen size may not be large enough.
In case of issue, a non-UI mode is available, using the flag `ui=false`.
