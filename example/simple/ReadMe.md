# Simple Example

This is a simple example showing how to use the framework.

An _HTTP_ server is listening on a standard 80-th port on _localhost_ at the 
root path (`http://localhost:80/`).

Two RPC functions are available: `Sum` and `Crash`. 

This simple example can also be used for primitive automatic and manual unit 
testing of the framework.

## Simple Testing 

Below you will find description of some use cases of this example and 
instructions about how to test the feature. Instructions are written in the 
following format: 

* _HTTP_ method and URL of a request
* _HTTP_ headers of the request
* _HTTP_ body of the request
* Response received from the RPC server

### 1. Wrong _HTTP_ method

* `GET localhost:80`
* Headers do not matter
* Body does not matter
* `405 Method Not Allowed`

### 2. Wrong content type in _HTTP_ request

* `POST localhost:80`
* Headers
  * `Content-Type` = `abrakadabra`
* Body does not matter
* `415 Unsupported Media Type`

### 3. Client does not accept JSON content type

* `POST localhost:80`
* Headers
  * `Content-Type` = `application/json`
  * `Accept` = `abrakadabra`
* Body does not matter
* `406 Not Acceptable`

### 4. Unparsable JSON in request

* `POST localhost:80`
* Headers
    * `Content-Type` = `application/json`
    * `Accept` = `*/*`
* Request:
```json
abrakadabra
```
* Response:
```json
{
  "jsonrpc": "M1",
  "id": null,
  "result": null,
  "error": {
    "code": -1,
    "message": "Request is not readable",
    "data": null
  },
  "ok": false
}
```

### 5. Invalid request

Notes: The fourth request field is absent.

* `POST localhost:80`
* Headers
    * `Content-Type` = `application/json`
    * `Accept` = `*/*`
* Request:
```json
{
  "id": "123",
  "method": "m",
  "params": {}
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "123",
  "result": null,
  "error": {
    "code": -2,
    "message": "Invalid request",
    "data": null
  },
  "ok": false
}
```

### 6. Unsupported protocol

* `POST localhost:80`
* Headers
    * `Content-Type` = `application/json`
    * `Accept` = `*/*`
* Request:
```json
{
  "jsonrpc": "x",
  "id": "123",
  "method": "m",
  "params": {}
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "123",
  "result": null,
  "error": {
    "code": -4,
    "message": "Unsupported protocol",
    "data": null
  },
  "ok": false
}
```

### 7. Unknown RPC method

* `POST localhost:80`
* Headers
  * `Content-Type` = `application/json`
  * `Accept` = `*/*`
* Request:
```json
{
  "jsonrpc": "M1",
  "id": "123",
  "method": "x",
  "params": {}
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "123",
  "result": null,
  "error": {
    "code": -8,
    "message": "Unknown method",
    "data": null
  },
  "ok": false
}
```

### 8. Invalid parameters

* `POST localhost:80`
* Headers
  * `Content-Type` = `application/json`
  * `Accept` = `*/*`
* Request:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "method": "Sum",
  "params": {
    "a": 245,
    "b": "10"
  }
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "result": null,
  "error": {
    "code": -16,
    "message": "Invalid parameters",
    "data": null
  },
  "meta": {
    "dur": 0
  },
  "ok": false
}
```

### 9. User error

* `POST localhost:80`
* Headers
  * `Content-Type` = `application/json`
  * `Accept` = `*/*`
* Request:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "method": "Sum",
  "params": {
    "a": 245,
    "b": 11
  }
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "result": null,
  "error": {
    "code": 1,
    "message": "overflow",
    "data": {
      "a": 245,
      "b": 11
    }
  },
  "meta": {
    "dur": 0
  },
  "ok": false
}
```

### 10. Internal RPC error

* `POST localhost:80`
* Headers
  * `Content-Type` = `application/json`
  * `Accept` = `*/*`
* Request:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "method": "Crash",
  "params": {}
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "result": null,
  "error": {
    "code": -32,
    "message": "Internal RPC error",
    "data": null
  },
  "meta": {
    "Laboratory": "L3",
    "dur": 1013
  },
  "ok": false
}
```

### 11. Normal result

* `POST localhost:80`
* Headers
  * `Content-Type` = `application/json`
  * `Accept` = `*/*`
* Request:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "method": "Sum",
  "params": {
    "a": 245,
    "b": 10
  }
}
```

* Response:
```json
{
  "jsonrpc": "M1",
  "id": "1",
  "result": {
    "c": 255
  },
  "error": null,
  "meta": {
    "dur": 0
  },
  "ok": true
}
```
