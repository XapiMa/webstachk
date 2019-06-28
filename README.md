# webStatusCheck
This is a simple tool to check endpoint status.
Access the endpoint every specified time and check if it matches the specified status.

## Installation
```
$ go get github.com/XapiMa/webStatusChecker
```

or

```
$ git clone https://github.com/XapiMa/webStatusChecker.git
$ go build main.go
```

If you need a different Architecture executable file:

```
 $ git clone https://github.com/XapiMa/webStatusChecker.git
 $ GOOS=linux GOARCH=amd64 go build main.go -o webStatusChecker
```
Please refer to the official document for details of available environment.
https://golang.org/doc/install/source#environment

## Usage
### Write target.csv 
Write target.csv with the url, status, and access_time_interval  column in the same directory as the executable file.

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
$ webStatusChecker
```

If you want to specify the location of target.csv :
```
$ webStatusChecker -t path/to/target.csv
```

If you want to write the result to a file:
```
$ webStatusChecker -o path/to/output/file
```

```
Usage of webStatusChecker:
  -l int
    	Monitoring time (second). In the case of 0, it is infinite. (default 0)
  -n int
    	Parallel number. (default 200)
  -o string
    	output file path. If not set, it will be output to standard output
  -t string
    	path to target.csv (default "In the same directory as the executable file")
  -v	show verbose
```
