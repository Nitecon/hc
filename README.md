# HTTP/2 Client (hc)

The `hc` package provides a user-friendly HTTP/2 client for Go. It supports all the standard HTTP methods (GET, POST, PUT, DELETE) and has built-in compression and JSON handling.

## Features

- HTTP/2 support.
- Built-in gzip compression for request data.
- Utility functions for common actions.
- Enhanced JSON support for sending and receiving data.

## Installation

```sh
go get your-repo-path/hc
```

## Making Requests

For standard requests:

```go
err := c.Get("https://example.com")
if err != nil {
log.Fatal(err)
}

err = c.Post("https://example.com", []byte("your data here"))
if err != nil {
log.Fatal(err)
}

err = c.Put("https://example.com", []byte("your data here"))
if err != nil {
log.Fatal(err)
}

err = c.Delete("https://example.com")
if err != nil {
log.Fatal(err)
}

```

For JSON requests:

```go
data := map[string]string{"key": "value"}
err = c.PostJson("https://example.com", data)
if err != nil {
log.Fatal(err)
}
```

## Accessing Response Data

Check the status code:

```go

statusCode := c.Status()
statusText := c.StatusText()

```

Read JSON data:

```go
var responseStruct YourStructType
err := c.ReadJson(&responseStruct)
if err != nil {
    log.Fatal(err)
}
```

## Note
Always make sure to check the error return after making a request before you access any other utility functions. This ensures that the request was successful and that the data is available for extraction.

## Contributing
Contributions are welcome! Please fork the repository and open a pull request with your changes, or open an issue with any suggestions, bugs, or feedback.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.