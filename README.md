# GoHooked

GoHooked is an application to handle time-consuming requests externally and receive the responses through a webhook

## Running the application

GoHooked has two params to define on start:
- _n_ is the number of workers to run concurrently (1 by default)
- _m_ is the maximum number of tasks to be scheduled (1 by default)

The maximum number of tasks will always be, at least, equal to the number of workers.

```bash
./gohooked -n 300 -m 1000
```

## Scheduling tasks

To schedule a task send a POST request to `http://localhost:4000/tasks`. It expects the following body:

- id: string with some information to identify the request on the webhook
- url: string with the url to be called
- payload: JSON string with the payload to be sent to the url
- headers: list of objects with key and value (both strings) to be included on the request
- callbackUrl: string with the webhook url

```json
{
    "id": "request-01",
    "url": "https://api.example.com/v1/users",
    "payload": "{\"name\":\"John Doe\",\"age\":33,\"isActive\":true}",
    "headers": [
        {
            "key": "Authorization",
            "value": "Bearer secret-token"
        },
        {
            "key": "Custom-Header",
            "value": "xz-123"
        }
    ],
    "callbackUrl": "https://api.mydomain.com/webhooks"
}
```

## Webhook response

When the task is done (request finished), the response payload is sent to the callback url with the id:

- id: string, same information sent when scheduling the task
- payload: JSON string with the url response

```json
{
    "id": "request-01",
    "payload": "{\"id\":99,\"name\":\"John Doe\"}"
}
```
