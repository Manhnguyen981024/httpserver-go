# Introduce 
We're going to build an Chirpy and It is a social network similar to Twitter.
It's a restful API server that allows users to:
- Add user `POST /api/users`
- Create Chirp `POST /api/chirps`
- GET ALL chirps `GET /api/chirps`
- Get chirps by chirp id `GET /api/chirps/{chirpID}`
- Login `POST /api/login`
- Refresh token `POST /api/refresh`
- Update user `PUT /api/users`
- Delete chirp by id `DELETE /api/chirps/{chirpID}`
- Webhook `POST /api/polka/webhooks`

# Learning Goals
- Understand what web servers are and how they power real-world web applications
- Build a production-style HTTP server in Go, without the use of a framework
- Use JSON, headers, and status codes to communicate with clients via a RESTful API
- Learn what makes Go a great language for building fast web servers
- Use type safe SQL to store and retrieve data from a Postgres database
- Implement a secure authentication/authorization system with well-tested cryptography libraries
- Build and understand webhooks and API keys

# install
## Tech Stack
- Go 1.25+
- PostgreSQL
- github.com/lib/pq
- github.com/google/uuid
- goose 
- sqlc

# Config
We'll use a single JSON file to keep track of two things:
- Who is currently logged in
- The connection credentials for the PostgreSQL database

- Manually create a config file in your home directory, ~/.gatorconfig.json, with the following content:
`json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
`

# install postgresql
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

# Goose Migrations
`go install github.com/pressly/goose/v3/cmd/goose@latest`

- Up one migration
`goose postgres <connection_string> up` 

- Down one migration
`goose postgres <connection_string> down`

# Install sqlc
`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
use `sqlc generate` to generate sql

# API details
## Add user `POST /api/users`
- Create a User for our server
*URL* : `POST /api/users`
*Request body* :
`json {
  "email": "mloneusk@example.co",
  "password": "abc123",
}`
*Response body* :
*HTTP status code success* : `HTTP 201 Created`
`json
{
  "id": "50746277-23c6-4d85-a890-564c0044c2fb",
  "created_at": "2021-07-07T00:00:00Z",
  "updated_at": "2021-07-07T00:00:00Z",
  "email": "user@example.com"
}`
*HTTP status code error* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Create Chirp `POST /api/chirps`
- Post a chirp to our social media
*URL* : `POST /api/chirps`
*Request header* :
`
Authorization: Bearer <token>
`
*Request body* :
`json {
  "body": "message",
}`
*Response body* :
*HTTP success status code* : `HTTP 201 Created`
`json
{
  "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
  "created_at": "2021-01-01T00:00:00Z",
  "updated_at": "2021-01-01T00:00:00Z",
  "body": "Hello, world!",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## GET ALL chirps `GET /api/chirps`
*URL* : `GET /api/chirps`
*Optional parameters* :
`
  ?author_id=id  
  ?sort=asc or desc 
`
- ?author_id=id  : get chirps by author id
- ?sort=asc or desc sort chirps

*Response body* :
*HTTP success status code* : `HTTP 200 OK`
`json
[
  {
    "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
    "created_at": "2021-01-01T00:00:00Z",
    "updated_at": "2021-01-01T00:00:00Z",
    "body": "Yo fam this feast is lit ong",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  },
  {
    "id": "f0f87ec2-a8b5-48cc-b66a-a85ce7c7b862",
    "created_at": "2022-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "body": "What's good king?",
    "user_id": "123e4567-e89b-12d3-a456-426614174000"
  }
]
`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Get chirps by chirp id `GET /api/chirps/{chirpID}`
*URL* : `GET /api/chirps/{chirpID}`
*Path* :
`
  chirpID: id of the chirp
`
- Example: GET /api/chirps/94b7e44c-3604-42e3-bef7-ebfcc3efff8f
*Response body* :
*HTTP success status code* : `HTTP 200 OK`
`json
{
  "id": "94b7e44c-3604-42e3-bef7-ebfcc3efff8f",
  "created_at": "2021-01-01T00:00:00Z",
  "updated_at": "2021-01-01T00:00:00Z",
  "body": "fr? no clowning?",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Login `POST /api/login`
*URL* : `POST /api/login`
- Login to our social media system to post message
*Request body* :
`json
{
  "email":"email",
  "password":"password"
}`
*Response body* :
*HTTP success status code* : `HTTP 200 OK`
`json 
{
  "id": "5a47789c-a617-444a-8a80-b50359247804",
  "created_at": "2021-07-01T00:00:00Z",
  "updated_at": "2021-07-01T00:00:00Z",
  "email": "lane@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "refresh_token": "56aa826d22baab4b5ec2cea41a59ecbba03e542aedbb31d9b80326ac8ffcfa2a"
}
`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Refresh token `POST /api/refresh`
*URL* : `POST /api/refresh`
- Refresh a access token by using refresh token
- This new endpoint does not accept a request body, but does require a refresh token to be present in the headers
*Request header* :
`
Authorization: Bearer <token>
`
*Response body* :
*HTTP success status code* : `HTTP 200 OK`
`json 
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Revoke token `POST /api/revoke`
*URL* : `POST /api/refresh`
- This new endpoint does not accept a request body, but does require a refresh token to be present in the headers
*Request header* :
`
Authorization: Bearer <token>
`
*Response body* :
*HTTP success status code* : `HTTP 204 No Content`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Update user `PUT /api/users`
*URL* : `PUT /api/users`
- Endpoint so that users can update their own (but not others') email and password.
*Request header* :
`
Authorization: Bearer <token>
`
*Request body* :
`json
{
  "email":"new email",
  "password":"new password"
}`
*Response body* :
*HTTP success status code* : `HTTP 200 OK`
`json 
{
  "id": "5a47789c-a617-444a-8a80-b50359247804",
  "created_at": "2021-07-01T00:00:00Z",
  "updated_at": "2021-07-01T00:00:00Z",
  "email": "lane@example.com",
}
`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

- Delete chirp by id `DELETE /api/chirps/{chirpID}`
*URL* : `DELETE /api/chirps/{chirpID}`
- Endpoint so that deletes a chirp from the database by its id.
*Request header* :
`
Authorization: Bearer <token>
`
*Request body* :
`json
{
  "email":"new email",
  "password":"new password"
}`
*Response body* :
*HTTP success status code* : `HTTP 204 No Content`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

## Webhook `POST /api/polka/webhooks`
*URL* : `POST /api/polka/webhooks`
- They will send us webhooks whenever a user subscribes to Chirpy Red
*Request header* :
`
Authorization: Apikey <token>
`
*Request body* :
`json
{
  "data": {
    "user_id": "${userID}"
  },
  "event": "user.upgraded"
}
`
*Response body* :
*HTTP success status code* : `HTTP 204 No Content`
*HTTP error status code* : `HTTP 400, 500`
`json
{
  "error":"error message"
}`

