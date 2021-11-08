# go-gateway

An easy way to make multiple backend accessible via one domain. Each service gets its own URL prefix and is then proxied without that prefix to the appropriate server.

Services look like this:
```json
[
        { 
            "name" : "Service 1",
            "prefix" : "/service/1",
            "port" : "3000" 
        },
        ...
]
```

Simply run:
```
go build
```
