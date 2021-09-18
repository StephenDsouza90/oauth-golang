# oauth-golang
An example of OAuth2.0 written in Golang.

## Run Server

``` bash
$ cd example/server
$ go build server.go
$ ./server
```

## Run Client

```
$ cd example/client
$ go build client.go
$ ./client
```

### Open the browser

[http://localhost:9094](http://localhost:9094)

```
username: test
```

```
{
  "access_token": "GIGXO8XWPQSAUGOYQGTV8Q",
  "token_type": "Bearer",
  "refresh_token": "5FBLXQ47XJ2MGTY8YRZQ8W",
  "expiry": "2019-01-08T01:53:45.868194+08:00"
}
```

