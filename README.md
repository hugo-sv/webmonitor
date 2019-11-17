# webmonitor

Webmonitor is a CLI web monitoring and alerting tool written in Go.

## Usage

### Build the project

From the project directory, run

```shell
go install
```

### Run the project

#### Command Line Interface

To run the project, if the gopath is setup, you can simply run the following command :

```
webmonitor [options] [URL]...
```

Otherwise, `webmonitor` can be replaced by `go run main.go` from the project directory.

#### Options

```
-interval INT (default : 5)
    Interval (in seconds) at which to check the websites
-timeout INT (default : 10)
    Timeout (in seconds) for http requests
```

#### Examples

```shell
webmonitor -interval=5 -timeout=10 "https://www.google.com" "https://httpstat.us/300" "https://httpstat.us/200?sleep=120" "https://httpstat.us/200?sleep=10000"
```

#### User Interface

You can press :

- q to quit
- c to clear alerts
- s to switch the statistics timeframe

#### Usage

The API `https://httpstat.us/{statusCode}?sleep={sleepTime}` can be use to test any response code and response time.

Disabling internet connection may as well emulate a website going down.

## Project Archtecture

### Main

In the `main` function, using go channels and tickers, operations are executed as they go. Usually :

- Every `interval` seconds, each websites' stats are requested. If, in a **2 min** timeframe, an alert is triggered, it is added in the UI.
- Every **10 sec**, the stats view is refreshed if the user is looking at a **10 min** timeframe.
- Every **1 min**, the stats view is refreshed if the user is looking at a **1h** timeframe.
- Every time a UI input is detected, the associated action is executed.
- Every time requested website's stats are recieved, they are processed.

### CLI

The `cli` module process the flags from the commande executed.

### Display

The `display` module handles every UI related actions :

- Generating the layouts
- Updating the panels
- String formating in the `format.go` script

### Monitor

The `monitor` module handles the HTTP get request, and compute a response time.

### Statistics

The `statistic` module compute the main statistics form the records of status code and response time in the considered timeframe :

- **Max** : Maximal response time
- **Avg** : Average response time
- **Availability** : Percent of succesfull requests (Status code 200)

## Alerting logic test

Run

```shell
webmonitor -interval=2 -timeout=5 "https://www.google.com" "https://httpstat.us/400" "https://httpstat.us/200" "https://httpstat.us/200?sleep=10000"
```

As the program starts, two alerts should be raised :

- httpstat.us/400, not returning a 200 status code
- httpstat.us/200?sleep=10000, triggering timeouts

Once the alerts are raised, stop your internet connection.
After less than 30 seconds, alerts should be raised from every websites.

Restart your internet connection, alerts from google.com and httpstat.us/200 being up again should be raised within 2 min.

## Notes

### Possible improvements

#### Max Response Time

If needed, it might be possible to improve the time performances in a trade-off with space complexity.

#### Tickers start time

Tickers currently start too late in this project. With a 10 sec interval, the user has to wait for 10 sec before the websites are requested. It might be possible to force an early stard of the FetchTicker.

#### Alert system

The alert system could be ugraded, so that active alerts are highlighted, with data on their duration and an updated availability.

#### Status Code

For a given Website, displaying the status code repartiion with a pie chart might make the information more readable.

#### Timeout managment

Timeout detection and invalid URL currently produce the same error (408).

#### Unit tests

With more time, I would have liked to follow a test driven devlopment.
I only implemented a few tests for the `monitoring` module.
