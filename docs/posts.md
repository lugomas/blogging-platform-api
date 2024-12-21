# Posts API

## Endpoint: Create Post

**Method:** `POST`  
**URL:** `/posts`

### Request Body
```json
{
  "title": "string",
  "content": "string",
  "category": "string",
  "tags": ["string"]
}
```

### Response
```json
{
"id": "uuid",
"title": "string",
"content": "string",
"category": "string",
"tags": ["string"],
"createdAt": "string",
"updatedAt": "string"
}
```