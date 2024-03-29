# webstachk
This is a simple tool to check endpoint status.
Access the endpoint every specified time and check if it matches the specified status.

## Installation
```
$ go get github.com/XapiMa/webstachk
```

or

```
$ git clone https://github.com/XapiMa/cmd/webstachk.git
$ go build ./webstachk
```

If you need a different Architecture executable file:

```
 $ git clone https://github.com/XapiMa/webstachk.git
 $ cd webstachk
 $ GOOS=linux GOARCH=amd64 go build -o webstachk ./cmd/webstachk
```
Please refer to the official document for details of available environment.
https://golang.org/doc/install/source#environment

## Usage
### Write config.csv 
Write config.csv with the url, status, and access_time_interval  column in the same directory as the executable file.

```
url,status_code,access_time_interval
url,status_code,access_time_interval
url,status_code,access_time_interval
url,status_code,access_time_interval
url,status_code,access_time_interval
```

ex.
```
http://example.com,200,20s
https://sample.com,301|302|303,1m
```

- Multiple status codes can be specified using a pipe('|').
- Access time interval can be selected from day('d'), hour('h'), minute('m') and second('s')

## Execution
```
$ webstachk
```

If you want to specify the location of config.csv :
```
$ webstachk -t path/to/config.csv
```

If you want to write the result to a file:
```
$ webstachk -t path/to/config.csv -o path/to/output/file
```

```
Usage of webstachk:
  -n int
    	Parallel number. (default 200)
  -o string
    	output file path. If not set, it will be output to standard output
  -t string
    	path to config.csv (required)
  -i int
      interval to self health check(second). In case of 0 it does not check
  -j bool
      change output format to json
```
