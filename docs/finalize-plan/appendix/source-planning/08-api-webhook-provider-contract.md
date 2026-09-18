# 08 — API, Webhook, dan Provider Contract

## API principles

- HTTPS only.
- JSON request/response.
- Versioned API.
- Explicit idempotency.
- Consistent error envelope.
- Request ID selalu dikembalikan.
- Secret tidak pernah dikembalikan setelah create.
- Public contract memakai `id: string` dari `public_id`.
- Internal BIGSERIAL `id` dan field `public_id` tidak pernah muncul di API response, Inertia props, UI network, atau customer webhook.

## Authentication

```http
Authorization: Bearer sk_live_xxx
Idempotency-Key: order-123-notification
X-Request-Id: req-client-123
```

API key memiliki:

- workspace scope;
- environment (`test`/`live`);
- scopes (`messages:write`, `messages:read`, `contacts:read`);
- created/last used/revoked timestamps;
- optional IP restrictions.

## Send message

```http
POST /v1/messages
Content-Type: application/json
Authorization: Bearer sk_live_xxx
Idempotency-Key: invoice-1001
```

```json
{
  "channel_id": "chn_123",
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
  "metadata": {
    "order_id": "ord_1001"
  }
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

Error class:

- `authentication_error`
- `permission_denied`
- `invalid_request`
- `consent_required`
- `template_not_approved`
- `recipient_invalid`
- `provider_rate_limited`
- `provider_unavailable`
- `quota_exceeded`
- `idempotency_conflict`

### ID mapping rule

```text
internal model: id BIGSERIAL, public_id nanoid(21)
external DTO:  id string = public_id
external DTO:  public_id tidak ada
```

Serializer/mapper wajib digunakan pada semua API, Inertia props, UI payload, webhook, dan event customer. Jangan mengembalikan database model secara langsung.

## Message status endpoint

```http
GET /v1/messages/msg_123
GET /v1/messages/msg_123/events
```

Status response wajib membedakan:

- platform accepted;
- provider accepted;
- sent;
- delivered;
- read;
- failed;
- canceled.

## Webhook endpoint customer

```http
POST https://customer.example.com/webhooks/whatsapp
X-Webhook-Id: evt_123
X-Webhook-Timestamp: 1725430000
X-Webhook-Signature: v1=...
```

Signature proposal:

```text
signed_payload = timestamp + "." + raw_body
signature = HMAC-SHA256(webhook_secret, signed_payload)
```

Rules:

- timestamp tolerance, misalnya 5 menit;
- constant-time compare;
- reject replayed event ID;
- raw body dipakai untuk verification;
- endpoint secret dapat rotate dengan overlap period;
- retry dengan exponential backoff;
- customer dapat replay event dari dashboard.

## Webhook event types

- `channel.connected`
- `channel.disconnected`
- `message.received`
- `message.accepted`
- `message.sent`
- `message.delivered`
- `message.read`
- `message.failed`
- `contact.opted_out`
- `template.updated`
- `campaign.completed`

## Provider adapter model

```go
type Provider string

const (
    ProviderOfficial  Provider = "official"
    ProviderWhatsmeow Provider = "whatsmeow"
)

type Capabilities struct {
    SendText       bool
    SendMedia      bool
    ReceiveMessage bool
    DeliveryStatus  bool
    Templates       bool
    Groups          bool
    Pairing         bool
}

type SendRequest struct {
    ChannelID      string
    To             string
    Content        Content
    IdempotencyKey string
}
```

### Official adapter

- Memanggil Cloud API/partner API.
- Mengelola Phone Number ID dan WABA context.
- Menerima provider webhook.
- Memetakan template/media/status/error.

### Whatsmeow adapter

- Menjalankan session worker terisolasi.
- QR/pairing lifecycle.
- Session persistence terenkripsi.
- Reconnect/backoff.
- Raw protocol event → normalized event.
- Feature flags dan kill switch.
- Risk acknowledgement wajib.

## API documentation experience

Developer Center wajib menyediakan:

- quickstart cURL, Go, JavaScript, PHP, Python;
- Postman collection;
- test credentials/sandbox;
- webhook tester;
- request log yang ter-redact;
- API version/changelog;
- rate limits dan retry guidance;
- copyable error resolution.

## Backward compatibility

- Jangan mengubah arti status secara diam-diam.
- Deprecation notice minimal 90 hari untuk paid API.
- Changelog wajib.
- Contract tests untuk provider adapters.
- Versioned webhook schema.
