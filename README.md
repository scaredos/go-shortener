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
- GET `/api/v1/addUrl` with parameters `?url=` and `password`
- For a custom shortened URL, specify the parameter `short`
- Example Request:

```
Add URL: example.com:1720/api/v1/addUrl?url=https://github.com/scaredos&short=github&password=example
Use URL: example.com:1720/github
```

### Editing URLs 
- To edit a shortened URL, use the API
- GET `/api/v1/editUrl` with parameters `?url=&short=&password=`
- You must specify both URL (the new destination URL) and short (the shortened URL)
```
Add URL: example.com:1720/api/v1/addUrl?url=https://github.com/scaredos&short=github&password=example
Edit URL: example.com:1720/api/v1/editUrl?url=https://bing.com/&short=github&password=example
```

### TODO
- Add password for adding URLs to service
- Add option for logging requests (IP, URL, User-Agent)
