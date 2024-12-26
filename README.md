# Chirpy
Chirpy is a simple platform to chirp, users can post text based contents.

This is the api / backend codebase written in Go for chirpy.

## Features
- Register a new user
- Update user email and password
- Login a user
- Refresh user token
- Revoke user token / Logout a user
- Post a chirp
- Get all chirps
  - Get all chirps by user id
  - Get all chirps sorted by created time
- Get a chirp by id
- Delete a chirp
- Polka Webhook for payments (simulated)
- See Metrics (Admin), how many times the landing page was visited
- Reset (Dev Only), reset the database
- Health Check, check if the server is running
- Landing Page, browser accessible. Chirpy landing page

## Project Setup
1. Clone the repository
```bash
git clone https://github.com/mohits-git/experiments.git
cd experiments/go-server
```

2. Install dependencies
```bash
go mod download
```

3. Run the server
```bash
go run .
```

4. (OR) To run the tests
```bash
go test ./...
```

5. (OR) Build and run the server
```bash
go build .
./go-server
```

## API Documentation

### Users

**Register a new user**

```http
POST /api/users

Headers:
  Content-Type: application/json
```
```json
{
  "email": "your_email@mail.com",
  "password": "your_password"
}
```
Response:
```http
201 Created
```
```json
{
  "id": "user_id",
  "email": "your_mail@mail.com",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "is_chirpy_red": "false"
}
```

**Update user email and password**

```http
PUT /api/users

Headers:
  Content-Type: application/json
  Authorization: Bearer <access_token>
```
```json
{
  "email": "your_email@mail.com",
  "password": "your_password"
}
```
Response:
```http
200 OK
```
```json
{
  "id": "user_id",
  "email": "your_mail@mail.com",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "is_chirpy_red": "false"
}
```


### Auth

**Login a user**

```http
POST /api/login
Headers:
  Content-Type: application/json
```
```json
{
  "email": "your_email@mail.com",
  "password": "your_password"
}
```
Response:
```http
200 OK
```
```json
{
  "id": "user_id",
  "email": "your_mail@mail.com",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "is_chirpy_red": "false",
  "token": "<access_token>",
  "refresh_token": "<refresh_token>"
}
```

**Refresh user token**

```http
POST /api/refresh
Headers:
  Content-Type: application/json
  Authorization: Bearer <refresh_token>
```
Response:
```http
200 OK
```
```json
{
  "token": "<access_token>"
}
```

**Revoke user token / Logout a user**

```http
POST /api/logout
Headers:
  Content-Type: application/json
  Authorization: Bearer <refresh_token>
```
Response:
```http
204 OK
```

### Chirps

**Post a chirp**

```http
POST /api/chirps
Headers:
  Content-Type: application/json
  Authorization: Bearer <access_token>
```
```json
{
  "body": "your_chirp_content"
}
```
Response:
```http
201 Created
```
```json
{
  "id": "chirp_id",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "body": "your_chirp_content",
  "user_id": "user_id"
}
```

**Get all chirps**

```http
GET /api/chirps?[auther_id=<user_id>]&[sort=(asc|desc)]
Headers:
  Content-Type: application/json
```
Response:
```http
200 OK
```
```json
[
  {
    "id": "chirp_id",
    "created_at": "<timestamp>",
    "updated_at": "<timestamp>",
    "body": "your_chirp_content",
    "user_id": "user_id"
  }
]
```

**Get a chirp by id**

```http
GET /api/chirps/:id
Headers:
  Content-Type: application/json
```
Response:
```http
200 OK
```
```json
{
  "id": "chirp_id",
  "created_at": "<timestamp>",
  "updated_at": "<timestamp>",
  "body": "your_chirp_content",
  "user_id": "user_id"
}
```

**Delete a chirp**

```http
DELETE /api/chirps/:id
Headers:
  Content-Type: application/json
  Authorization: Bearer <access_token>
```
Response:
```http
204 OK
```

**Polka Webhook for payments (simulated)**

- Event: `user.upgrade`

```http
POST /api/polka/webhooks
Headers:
  Content-Type: application/json
  Authorization: ApiKey <polka_api_key>
```
```json
{
  "event": "user.upgrade",
  "data": {
    "user_id": "user_id"
  }
}
```
Response:
```http
204 No Content
```

### Admin

**See Metrics**
- Accessible on browser

```http
GET /admin/metrics
```
Response:
```http
200 OK
```
HTML Page with metrics

**Reset (Dev Only)**
- Reset the database

```http
POST /admin/reset
```
Response:
```http
200 No Content
```

### Server

**Health Check**
- Check if the server is running

```http
GET /api/healthz
```
Response:
```http
200 OK
```

**Landing Page**
- Landing page for the platform
- Browser Accessible

```http
GET /app
```
Response:
```http
200 OK
```
HTML Page with landing page
