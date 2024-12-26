# AdRouter API Documentation

The AdRouter API enables users to manage and deliver ad campaigns through a RESTful interface. The API supports CRUD operations on campaigns, as well as targeting and delivery features. Each endpoint is detailed below with input and output specifications.

## Base URL

The API is hosted at:

```
https://adrouter.site
```

---

## Endpoints

### 1. **Health Check**

#### `GET /ping`

Checks the API server's health and latency.

**Response:**

```json
{
  "ping": "<time-elapsed>"
}
```

---

### 2. **Campaign Management**

#### `GET /v1/get_campaign/:cid`

Fetches details of a specific campaign.

**Path Parameters:**

- `cid`: Campaign ID (string, required)

**Response:**

- `200 OK`: Campaign details.
- `404 Not Found`: Campaign not found.

---

#### `POST /v1/create_campaign`

Creates a new campaign.

**Request Body:**

```json
{
  "cid": "string",
  "name": "string (6-32 characters)",
  "img": "string",
  "cta": "string",
  "app": "string (optional)",
  "app_rule": "include | exclude (needed only if app is given)",
  "country": "string (optional)",
  "country_rule": "include | exclude (needed only if country is given)",
  "os": "string (optional)",
  "os_rule": "include | exclude (needed only if os is given)"
}
```

**Response:**

- `201 Created`: Campaign created.
- `400 Bad Request`: Validation errors.

---

#### `POST /v1/add_campaign`

Adds an new campaign without any rules.

**Request Body:**

```json
{
  "cid": "string",
  "name": "string",
  "img": "string",
  "cta": "string"
}
```

**Response:**

- `201 Created`: Campaign added.
- `400 Bad Request`: Validation errors.

---

#### `PATCH /v1/toggle_status/:cid`

Toggles the active status of a campaign.

**Path Parameters:**

- `cid`: Campaign ID (string, required)

**Response:**

- `200 OK`: Status toggled successfully.
- `404 Not Found`: Campaign not found.

---

#### `PATCH /v1/update_campaign_name`

Updates the name of a campaign.

**Request Body:**

```json
{
  "cid": "string",
  "name": "string (6-32 characters)"
}
```

**Response:**

- `200 OK`: Name updated successfully.
- `400 Bad Request`: Validation errors.

---

#### `PATCH /v1/update_campaign_image`

Updates the image of a campaign.

**Request Body:**

```json
{
  "cid": "string",
  "img": "string"
}
```

**Response:**

- `200 OK`: Image updated successfully.
- `400 Bad Request`: Validation errors.

---

#### `PATCH /v1/update_campaign_cta`

Updates the call-to-action (CTA) of a campaign.

**Request Body:**

```json
{
  "cid": "string",
  "cta": "string"
}
```

**Response:**

- `200 OK`: CTA updated successfully.
- `400 Bad Request`: Validation errors.

---

### 3. **Targeting Management**

#### `POST /v1/add_target_app`

Adds targeting by application ID.

**Request Body:**

```json
{
  "cid": "string",
  "app": "string",
  "rule": "include | exclude"
}
```

**Response:**

- `201 Created`: Target app added successfully.
- `400 Bad Request`: Validation errors.

---

#### `POST /v1/add_target_country`

Adds targeting by country.

**Request Body:**

```json
{
  "cid": "string",
  "country": "string",
  "rule": "include | exclude"
}
```

**Response:**

- `201 Created`: Target country added successfully.
- `400 Bad Request`: Validation errors.

---

#### `POST /v1/add_target_os`

Adds targeting by operating system.

**Request Body:**

```json
{
  "cid": "string",
  "os": "string",
  "rule": "include | exclude"
}
```

**Response:**

- `201 Created`: Target OS added successfully.
- `400 Bad Request`: Validation errors.

---

#### `PATCH /v1/update_target_app`

Updates targeting by application ID.

**Request Body:**

```json
{
  "cid": "string",
  "app": "string",
  "rule": "include | exclude"
}
```

**Response:**

- `200 OK`: Target app updated successfully.
- `400 Bad Request`: Validation errors.

---

#### `PATCH /v1/update_target_country`

Updates targeting by country.

**Request Body:**

```json
{
  "cid": "string",
  "country": "string",
  "rule": "include | exclude"
}
```

**Response:**

- `200 OK`: Target country updated successfully.
- `400 Bad Request`: Validation errors.

---

#### `PATCH /v1/update_target_os`

Updates targeting by operating system.

**Request Body:**

```json
{
  "cid": "string",
  "os": "string",
  "rule": "include | exclude"
}
```

**Response:**

- `200 OK`: Target OS updated successfully.
- `400 Bad Request`: Validation errors.

---

### 4. **Delivery**

#### `GET /v1/delivery`

Fetches available campaigns based on targeting criteria.

**Query Parameters:**

- `app`: Application ID (string, required)
- `country`: Country (string, required)
- `os`: Operating System (string, required)

**Response:**

- `200 OK`: List of campaigns matching the criteria.
- `204 No Content`: No campaigns available.

---

### 5. **Deletion**

#### `DELETE /v1/delete_campaign/:cid`

Deletes a campaign by ID.

**Path Parameters:**

- `cid`: Campaign ID (string, required)

**Response:**

- `200 OK`: Campaign deleted successfully.
- `404 Not Found`: Campaign not found.

---

#### `DELETE /v1/delete_target_app/:cid`

Deletes a target app from a campaign.

**Path Parameters:**

- `cid`: Campaign ID (string, required)

**Response:**

- `200 OK`: Target app deleted successfully.
- `404 Not Found`: Resource not found.

---

#### `DELETE /v1/delete_target_country/:cid`

Deletes a target country from a campaign.

**Path Parameters:**

- `cid`: Campaign ID (string, required)

**Response:**

- `200 OK`: Target country deleted successfully.
- `404 Not Found`: Resource not found.

---

#### `DELETE /v1/delete_target_os/:cid`

Deletes a target OS from a campaign.

**Path Parameters:**

- `cid`: Campaign ID (string, required)

**Response:**

- `200 OK`: Target OS deleted successfully.
- `404 Not Found`: Resource not found.

---

### 6. **Error Handling**

All error responses include the following format:

```json
{
  "error": "<error-message>"
}
```

## Notes

- All endpoints return responses in JSON format.
- Use appropriate HTTP methods (GET, POST, PATCH, DELETE) to interact with the API.
- Ensure valid inputs to avoid `400 Bad Request` errors.
