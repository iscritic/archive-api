# Archive & Email API

A Go-based API that provides functionality to handle file archiving and email distribution. Built with the Gin framework, this API allows users to retrieve archive details, create zip archives, and send files via email to specified recipients.

## Features

- **Archive Information Retrieval**: Upload files to retrieve detailed metadata.
- **Zip Archive Creation**: Create a ZIP archive from multiple files.
- **Send Files via Email**: Distribute uploaded files as email attachments to multiple recipients.
## Prerequisites

- **Go**: Version 1.18 or higher.
- **SMTP Server**: Required for sending emails.

## Installation

1. **Clone the Repository**:

   ```sh
   git clone https://github.com/yourusername/archive-email-api.git
   cd archive-email-api
   ```

2. **Install Dependencies**:

   ```sh
   go mod tidy
   ```

3. **Build the Project**:

   ```sh
   go build -o archive-email-api
   ```

## Configuration

Set up environment variables or use a configuration file to define SMTP settings.

### Environment Variables

Create a `.env` file:

```sh
CONFIG_PATH=./config/local.yaml

SMTP_USER=YOUR_GMAIL
SMTP_PASSWORD=YOUR_APP_PASSWORD

```

## Running the Server

Run the server with the following command:

```sh
./archive-email-api
```

The server will run by default on `http://localhost:8080`.

## API Endpoints

### 1. Ping

**Endpoint**: `GET /api/ping`

**Description**: Health check to ensure the server is running.

**Response**:

```json
{
  "message": "pong"
}
```

### 2. Get Archive Information

**Endpoint**: `POST /api/archive/information`

**Description**: Upload a file to retrieve archive metadata.

**Request**: Form data with `file` field.

**Example**:

```sh
curl -X POST http://localhost:8080/api/archive/information \
  -F "file=@/path/to/your/file.zip"
```

### 3. Create Archive

**Endpoint**: `POST /api/archive/files`

**Description**: Upload multiple files to create a ZIP archive.

**Request**: Form data with `files[]` field.

**Example**:

```sh
curl -X POST http://localhost:8080/api/archive/files \
  -F "files[]=@/path/to/file1.txt" \
  -F "files[]=@/path/to/file2.jpg" \
  -o archive.zip
```

### 4. Send File to Emails

**Endpoint**: `POST /api/mail/file`

**Description**: Send an uploaded file as an attachment to multiple recipients.

**Request**: Form data with fields `file`, `emails`, `subject` (optional), and `body` (optional).

**Example**:

```sh
curl -X POST http://localhost:8080/api/mail/file \
  -F "file=@/path/to/your/file.pdf" \
  -F "emails=email1@example.com,email2@example.com" \
  -F "subject=Here is your file" \
  -F "body=Please find the attached file."
```

