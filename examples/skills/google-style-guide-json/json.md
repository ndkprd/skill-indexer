# Google JSON Style Guide — REST API Reference

Extended reference with concrete examples for REST API response/request shapes.

## Golden Rules

1. **camelCase** for property names
2. **Consistent structure** — predictable response format
3. **Include metadata** — pagination inside `data`
4. **Use arrays for collections** — wrapped under a named key (e.g. `items`)
5. **Omit null fields** — use `omitempty`
6. **ISO 8601 for dates** — e.g. `2026-06-10T14:30:00Z`

## Response Structure

### Single Resource

Fields sit directly under `data` — no `type`, `id`, or `attributes` nesting.

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Critical Alert to Slack",
    "channelType": "webhook",
    "channelName": "ops-slack",
    "subject": "[CRITICAL] {{.title}}",
    "messageBody": "Alert: {{.message}}",
    "config": "{\"url\":\"https://hooks.slack.com/...\"}",
    "createdAt": "2026-06-10T14:30:00Z",
    "updatedAt": "2026-06-10T14:30:00Z"
  }
}
```

### Collection (paginated)

Inspired by YouTube Data API v2.0 — pagination fields and `items` array live inside `data`.

```
GET /api/v1/templates?page=1&perPage=20
```

```json
{
  "data": {
    "totalItems": 42,
    "startIndex": 1,
    "itemsPerPage": 20,
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Critical Alert to Slack",
        "channelType": "webhook",
        "channelName": "ops-slack",
        "createdAt": "2026-06-10T14:30:00Z",
        "updatedAt": "2026-06-10T14:30:00Z"
      }
    ]
  }
}
```

### Create/Update Request

Fields sent directly at top level — no `data` wrapper, no `type`/`attributes` nesting.

```
POST /api/v1/templates
```

```json
{
  "name": "Critical Alert to Slack",
  "channelType": "webhook",
  "channelName": "ops-slack",
  "subject": "[CRITICAL] {{.title}}",
  "messageBody": "Alert: {{.message}}",
  "config": "{\"url\":\"https://hooks.slack.com/...\"}"
}
```

### Error Response

Google-style `error` object containing code, message, status text, and optional `errors` array with per-field details.

```json
{
  "error": {
    "code": 404,
    "message": "Template not found",
    "status": "NOT_FOUND",
    "errors": [
      {
        "reason": "templateNotFound",
        "message": "Template with id 'abc-123' does not exist"
      }
    ]
  }
}
```

### Common HTTP status → status string mapping

| Code | Status |
|------|--------|
| 400  | `BAD_REQUEST` |
| 401  | `UNAUTHORIZED` |
| 403  | `FORBIDDEN` |
| 404  | `NOT_FOUND` |
| 409  | `CONFLICT` |
| 422  | `UNPROCESSABLE_ENTITY` |
| 500  | `INTERNAL_SERVER_ERROR` |

## Property Naming

JSON keys are camelCase. Go struct field names stay PascalCase.

| DB column | Go field | JSON key |
|-----------|----------|----------|
| `api_key` | `ApiKey` | `apiKey` |
| `channel_type` | `ChannelType` | `channelType` |
| `channel_name` | `ChannelName` | `channelName` |
| `message_body` | `MessageBody` | `messageBody` |
| `template_ids` | `TemplateIDs` | `templateIds` |
| `created_at` | `CreatedAt` | `createdAt` |
| `updated_at` | `UpdatedAt` | `updatedAt` |

## Null Handling

Optional nullable fields use `omitempty` — absent when empty, never `null`.

Examples of nullable fields:
- `description` — omit if empty
- `config` — omit if empty
- `telegram_chat_id` — omit if empty
- `subject` — omit if empty

## Date Format

All timestamps use ISO 8601 / RFC 3339 format:

```
2026-06-10T14:30:00Z
2026-06-10T14:30:00+07:00
```

## Edge Cases

### Empty collection
```json
{
  "data": {
    "totalItems": 0,
    "startIndex": 1,
    "itemsPerPage": 20,
    "items": []
  }
}
```

### Resource not found
```json
{
  "error": {
    "code": 404,
    "message": "Template not found",
    "status": "NOT_FOUND",
    "errors": []
  }
}
```

### Validation error
```json
{
  "error": {
    "code": 400,
    "message": "Validation failed",
    "status": "BAD_REQUEST",
    "errors": [
      {
        "reason": "invalidField",
        "message": "name is required"
      },
      {
        "reason": "invalidField",
        "message": "channelType must be one of: smtp, webhook, ntfy, netbox, grafana"
      }
    ]
  }
}
```
