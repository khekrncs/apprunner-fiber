# API Documentation - CURL Examples

Base URL: `http://localhost:8080`

## Table of Contents

- [Health Check Endpoints](#health-check-endpoints)
- [User Management APIs](#user-management-apis)
- [File Management APIs](#file-management-apis)
- [Complete Workflow Example](#complete-workflow-example)

## Health Check Endpoints

### 1. Root Endpoint

Check if the service is running.

```bash
curl -X GET http://localhost:8080/
```

**Response:**

```json
{
  "status": "I'm running!"
}
```

### 2. Health Check

Get service health status.

```bash
curl -X GET http://localhost:8080/health
```

**Response:**

```json
{
  "status": "ok"
}
```

## User Management APIs

### 1. Create User

Create a new user in the system.

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "name": "John Doe",
    "metadata": {
      "role": "admin",
      "department": "engineering"
    }
  }'
```

**Response:**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "metadata": {
    "role": "admin",
    "department": "engineering"
  },
  "created_at": "2025-05-22T10:30:00Z",
  "updated_at": "2025-05-22T10:30:00Z"
}
```

### 2. Get User by ID

Retrieve a specific user by their ID.

```bash
curl -X GET http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000
```

**Response:**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "metadata": {
    "role": "admin",
    "department": "engineering"
  },
  "created_at": "2025-05-22T10:30:00Z",
  "updated_at": "2025-05-22T10:30:00Z"
}
```

### 3. Update User

Update an existing user's information.

```bash
curl -X PUT http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "metadata": {
      "role": "senior_admin",
      "department": "engineering",
      "team": "backend"
    }
  }'
```

**Response:**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "john.doe@example.com",
  "name": "John Smith",
  "metadata": {
    "role": "senior_admin",
    "department": "engineering",
    "team": "backend"
  },
  "created_at": "2025-05-22T10:30:00Z",
  "updated_at": "2025-05-22T11:45:00Z"
}
```

### 4. List All Users

Get a list of all users in the system.

```bash
curl -X GET http://localhost:8080/api/v1/users
```

**Response:**

```json
{
  "users": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "john.doe@example.com",
      "name": "John Smith",
      "metadata": {
        "role": "senior_admin",
        "department": "engineering"
      },
      "created_at": "2025-05-22T10:30:00Z",
      "updated_at": "2025-05-22T11:45:00Z"
    },
    {
      "id": "987fcdeb-51a2-43d1-9f4e-123456789abc",
      "email": "jane.doe@example.com",
      "name": "Jane Doe",
      "metadata": {
        "role": "user",
        "department": "marketing"
      },
      "created_at": "2025-05-22T09:15:00Z",
      "updated_at": "2025-05-22T09:15:00Z"
    }
  ],
  "count": 2
}
```

### 5. Delete User

Delete a user and all their associated files.

```bash
curl -X DELETE http://localhost:8080/api/v1/users/123e4567-e89b-12d3-a456-426614174000
```

**Response:**

```
HTTP 204 No Content
```

## File Management APIs

### 1. Upload File

Upload a file for a specific user.

```bash
curl -X POST http://localhost:8080/api/v1/files/upload/123e4567-e89b-12d3-a456-426614174000 \
  -F "file=@document.pdf"
```

**Alternative with explicit content type:**

```bash
curl -X POST http://localhost:8080/api/v1/files/upload/123e4567-e89b-12d3-a456-426614174000 \
  -F "file=@document.pdf" \
  -H "Content-Type: multipart/form-data"
```

**Response:**

```json
{
  "key": "users/123e4567-e89b-12d3-a456-426614174000/files/document.pdf",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "filename": "document.pdf",
  "content_type": "application/pdf",
  "size": 1048576,
  "etag": "\"d41d8cd98f00b204e9800998ecf8427e\"",
  "metadata": {
    "user-id": "123e4567-e89b-12d3-a456-426614174000",
    "original-name": "document.pdf"
  },
  "uploaded_at": "2025-05-22T12:00:00Z",
  "last_modified": "2025-05-22T12:00:00Z"
}
```

### 2. Download File

Download a specific file for a user.

```bash
curl -X GET http://localhost:8080/api/v1/files/123e4567-e89b-12d3-a456-426614174000/document.pdf \
  --output downloaded_document.pdf
```

**Download with headers visible:**

```bash
curl -X GET http://localhost:8080/api/v1/files/123e4567-e89b-12d3-a456-426614174000/document.pdf \
  -v \
  --output downloaded_document.pdf
```

**Response Headers:**

```
Content-Type: application/pdf
Content-Length: 1048576
ETag: "d41d8cd98f00b204e9800998ecf8427e"
```

### 3. List User Files

Get all files for a specific user.

```bash
curl -X GET http://localhost:8080/api/v1/files/123e4567-e89b-12d3-a456-426614174000
```

**Response:**

```json
{
  "files": [
    {
      "key": "document.pdf",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "filename": "document.pdf",
      "content_type": "",
      "size": 1048576,
      "etag": "\"d41d8cd98f00b204e9800998ecf8427e\"",
      "last_modified": "2025-05-22T12:00:00Z"
    },
    {
      "key": "image.png",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "filename": "image.png",
      "content_type": "",
      "size": 512000,
      "etag": "\"e4d909c290d0fb1ca068ffaddf22cbd0\"",
      "last_modified": "2025-05-22T11:30:00Z"
    }
  ],
  "count": 2
}
```

### 4. Delete File

Delete a specific file for a user.

```bash
curl -X DELETE http://localhost:8080/api/v1/files/123e4567-e89b-12d3-a456-426614174000/document.pdf
```

**Response:**

```
HTTP 204 No Content
```

## Complete Workflow Example

Here's a complete workflow demonstrating the API usage:

### Step 1: Check Service Health

```bash
curl -X GET http://localhost:8080/health
```

### Step 2: Create a User

```bash
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "name": "Demo User",
    "metadata": {"role": "demo"}
  }')

USER_ID=$(echo $USER_RESPONSE | jq -r '.id')
echo "Created user with ID: $USER_ID"
```

### Step 3: Upload Files

```bash
# Create a test file
echo "Hello, World!" > test.txt

# Upload the file
curl -X POST http://localhost:8080/api/v1/files/upload/$USER_ID \
  -F "file=@test.txt"

# Upload another file
echo '{"message": "Hello JSON"}' > data.json
curl -X POST http://localhost:8080/api/v1/files/upload/$USER_ID \
  -F "file=@data.json"
```

### Step 4: List User Files

```bash
curl -X GET http://localhost:8080/api/v1/files/$USER_ID
```

### Step 5: Download a File

```bash
curl -X GET http://localhost:8080/api/v1/files/$USER_ID/test.txt \
  --output downloaded_test.txt
```

### Step 6: Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Demo User",
    "metadata": {"role": "demo", "status": "active"}
  }'
```

### Step 7: Get Updated User

```bash
curl -X GET http://localhost:8080/api/v1/users/$USER_ID
```

### Step 8: Clean Up (Delete File)

```bash
curl -X DELETE http://localhost:8080/api/v1/files/$USER_ID/test.txt
```

### Step 9: List All Users

```bash
curl -X GET http://localhost:8080/api/v1/users
```

### Step 10: Delete User (Optional)

```bash
curl -X DELETE http://localhost:8080/api/v1/users/$USER_ID
```

## Error Responses

### 400 Bad Request

```json
{
  "error": "Invalid request body"
}
```

### 404 Not Found

```json
{
  "error": "User not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "Failed to save user: failed to put object: access denied"
}
```

## Testing Script

Save this as `test-api.sh` and make it executable:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🚀 Testing API endpoints..."

# 1. Health check
echo "1️⃣ Health check..."
curl -s $BASE_URL/health | jq
echo

# 2. Create user
echo "2️⃣ Creating user..."
USER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "name": "Test User"}')
USER_ID=$(echo $USER_RESPONSE | jq -r '.id')
echo "User ID: $USER_ID"
echo

# 3. Upload file
echo "3️⃣ Uploading file..."
echo "Test content" > test-file.txt
curl -s -X POST $BASE_URL/api/v1/files/upload/$USER_ID \
  -F "file=@test-file.txt" | jq
echo

# 4. List files
echo "4️⃣ Listing files..."
curl -s -X GET $BASE_URL/api/v1/files/$USER_ID | jq
echo

# 5. Download file
echo "5️⃣ Downloading file..."
curl -s -X GET $BASE_URL/api/v1/files/$USER_ID/test-file.txt
echo

# 6. List users
echo "6️⃣ Listing users..."
curl -s -X GET $BASE_URL/api/v1/users | jq '.count'
echo

# Clean up
rm test-file.txt

echo "✅ API testing complete!"
```

Run with:

```bash
chmod +x test-api.sh
./test-api.sh
```
