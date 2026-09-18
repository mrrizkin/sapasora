# 11 — API, Webhook, dan Provider Contract

## API principles

- HTTPS only.
- REST/JSON, versioned.
- Explicit idempotency.
- Consistent error envelope.
- Request ID selalu dikembalikan.
- Secret hanya ditampilkan saat create/rotate.
- Public contract hanya memakai `id: string` dari `public_id`.
- Internal `id BIGSERIAL` dan field `public_id` tidak pernah muncul di response, Inertia, UI network, webhook customer, atau public docs.

## Authentication

```http
Authorization: Bearer sk_live_xxx
Idempotency-Key: invoice-1001
X-Request-Id: req-client-123
```

API key memiliki workspace scope, environment, scopes (`messages:write`, `messages:read`, `contacts:read`), created/last-used/revoked timestamps, dan optional IP restriction.

## Send endpoint

```http
POST /v1/messages
Content-Type: application/json
Authorization: Bearer sk_live_xxx
Idempotency-Key: invoice-1001
```

```json
{
  "channel_id": "chn_public_123",
  "to": "+628123456789",
  "type": "template",
  "template": {
    "name": "invoice_ready",
    "language": "id",
    "variables": {
      "customer_name": "Budi",
      "invoice_number": "INV-1001"
    }
  },
  "metadata": {"order_id": "ord_public_1001"}
}
```

Response:

```json
{
  "id": "msg_public_123",
  "status": "queued",
  "request_id": "req_123",
  "provider": "official",
  "created_at": "2026-09-04T10:00:00Z"
}
```

## Error envelope

```json
{
  "error": {
    "code": "template_not_approved",
    "message": "Template belum disetujui oleh provider.",
    "next_action": "Buka Template Manager dan pilih template approved.",
    "retryable": false,
    "request_id": "req_123"
  }
}
```

Error classes: authentication, permission, invalid request, consent required, template not approved, invalid recipient, provider rate limited/unavailable, quota exceeded, idempotency conflict.

## Status

```http
GET /v1/messages/msg_public_123
GET /v1/messages/msg_public_123/events
POST /v1/messages/msg_public_123/cancel
```

Status internal/external: created, queued, sending, sent, delivered, read, failed, canceled. Bedakan platform accepted dan provider accepted.

## Customer webhook

```http
POST https://customer.example.com/webhooks/whatsapp
X-Webhook-Id: evt_public_123
X-Webhook-Timestamp: 1725430000
X-Webhook-Signature: v1=...
```

Signature proposal:

```text
signed_payload = timestamp + "." + raw_body
signature = HMAC-SHA256(webhook_secret, signed_payload)
```

Rules: timestamp tolerance, constant-time compare, replay protection, raw body verification, secret rotation overlap, exponential backoff, DLQ, controlled replay.

Event types:

- `channel.connected/disconnected`;
- `message.received/accepted/sent/delivered/read/failed`;
- `contact.opted_out`;
- `template.updated`;
- `campaign.completed`.

## Provider contract

Gunakan capability-specific interfaces, bukan satu interface raksasa:

```go
type MessageSender interface {
    Send(ctx context.Context, req SendRequest) (SendResult, error)
}

type ChannelHealth interface {
    Health(ctx context.Context, channelID string) (Health, error)
}
```

### Official adapter

Cloud API/partner API, WABA/Phone Number context, provider webhooks, template/media/status mapping.

### Whatsmeow adapter

Session worker, QR/pairing, encrypted persistence, reconnect/backoff, raw event normalization, feature flags, kill switch, risk acknowledgement.

## Endpoint roadmap

```text
POST /v1/messages
GET  /v1/messages/{id}
GET  /v1/messages/{id}/events
POST /v1/messages/{id}/cancel
POST /v1/contacts/import
GET  /v1/conversations
POST /v1/webhooks/test
POST /v1/api-keys/rotate
```

## Documentation

- OpenAPI versioned.
- Scalar sebagai API reference.
- cURL, Go, JavaScript, PHP, Python examples.
- Postman collection.
- Sandbox/test credentials.
- Webhook tester/log/replay.
- Rate limit, retry, error, deprecation, changelog.
- Contract tests dan breaking-change CI check.
