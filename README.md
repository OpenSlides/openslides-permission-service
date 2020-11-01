# Permission Service

The permisson service is a service and a libary for OpenSlides4 that tells, if a
user can see specific content.

## Build and run

go build ./cmd/permission/ && ./permission

## Example Request

```
curl http://localhost:9005/internal/permission/is_allowed -d '{"name":"topic.create","user_id":1}'
```

## Test

go test ./...


## Environment Variables

* `PERMISSION_HOST`: Host where the http service listens to. Default is an empty
  string which means all devices.
* `PERMISSION_PORT`: Port where the http services listens to. Default is 9005.
* `DATASTORE`: Sets the datastore service. `fake` (default) or `service`.
* `DATASTORE_READER_HOST`: Host of the datastore reader. The default is
  `localhost`.
* `DATASTORE_READER_PORT`: Port of the datastore reader. The default is `9010`.
* `DATASTORE_READER_PROTOCOL`: Protocol of the datastore reader. The default is
  `http`.