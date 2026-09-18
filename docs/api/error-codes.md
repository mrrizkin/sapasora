# API error contract

JSON API errors use a stable envelope:

```json
{
  "status": "Bad Request",
  "code": "bad_request",
  "message": "Request failed",
  "request_id": "..."
}
```

`code` is stable for client behavior; `message` is presentation text and may
be sanitized in production. `request_id` matches the `X-Request-ID` response
header and must be included in support reports. The current public categories
are `bad_request`, `validation_failed`, `unauthorized`, `forbidden`,
`not_found`, `conflict`, `provider_unavailable`, `rate_limited`, and
`internal_error`.
