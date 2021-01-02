# go-shortener
A Lightweight URL Shortener Written in Go

### Usage
- You must specify the port which to run the service on
```
git clone https://github.com/scaredos/go-shortener
cd go-shortener
go build server.go
./server.go --port <port>
```


### Adding URLs
- To add a URL to shorten, use the API
- GET `/api/v1/addUrl` with parameters `?url=`
- For a custom shortened URL, specify the parameter `short`
- Example Request:

```
Add URL: example.com:1720/api/v1/addUrl?url=https://github.com/scaredos&short=github
Use URL: example.com:1720/github
```

### TODO
- Add password for adding URLs to service
- Add option for logging requests (IP, URL, User-Agent)
