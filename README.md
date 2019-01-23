# Timeout

Server and client timeouts can be a bit tricky to understand, though examples can help tremendously. This repository includes a simple server and a client configured to hit that server. It allows the timeouts of each to be configured through command line flags.

In two separate terminals, one would start the server and then client and observe the logs and number of goroutines:
```
cd server
go run main.go -readHeader 10s
```

```
cd client
go run main.go -before 5s -after 5s -bad
```
