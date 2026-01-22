# Admin Backend

A comprehensive admin panel for managing the Astroway platform, built with Go and Gin framework.

## Features

- **User Management**: List users and deactivate accounts
- **Astrologer Management**: Manage astrologer status and visibility
- **Complaint Management**: Handle user service complaints, user general complaints, and astrologer complaints
- **Horoscope Management**: Create, update, delete, and manage horoscopes
- **Daily Report Scheduler**: Automated daily reports via email with CSV attachments
- **JWT Authentication**: Secure admin authentication
- **MongoDB Integration**: MongoDB database for data storage

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT
- **Scheduler**: Cron (robfig/cron/v3)
- **Email**: SMTP (Gmail)

## Prerequisites

- Go 1.21 or higher
- MongoDB instance
- Gmail account with App Password (for daily reports)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd ADMIN-BE
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the environment example file:
```bash
cp env.example .env
```

4. Configure your `.env` file with the required values (see Environment Variables section below)

5. Run the application:
```bash
go run main.go
```

The server will start on the port specified in the `PORT` environment variable (default: 8080).

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
MONGO_URI=your_mongodb_connection_string

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Server Configuration
PORT=8082
ADMIN_ID=your_admin_id

# Email Configuration for Daily Reports (Gmail)
EMAIL_USER=support@astrosway.com
EMAIL_PASS=your_gmail_app_password
```

### Gmail App Password Setup

For the daily report scheduler to work, you need to:

1. Enable 2-Step Verification on your Google Account
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a new app password for "Mail"
   - Use this 16-character password as `EMAIL_PASS`

## API Endpoints

### Authentication

- `POST /admin/login` - Admin login (returns JWT token)

### User Management

- `GET /admin/users` - List all users (with pagination and filters)
- `PATCH /admin/users/:userId/deactivate` - Deactivate a user

### Astrologer Management

- `GET /admin/astros` - List all astrologers
- `PATCH /admin/astros/:astroId/status` - Update astrologer status
- `PATCH /admin/astros/:astroId/visibility` - Toggle astrologer visibility

### User Service Complaints

- `GET /admin/user-service-complaints` - List user service complaints
- `GET /admin/user-service-complaints/:serviceType/:orderId` - Get complaint details
- `PATCH /admin/user-service-complaints/:reportId` - Accept/Reject complaint

### User General Complaints

- `GET /admin/user-general-complaints` - List user general complaints
- `PATCH /admin/user-general-complaints/:problemId/close` - Close a complaint

### Astrologer General Complaints

- `GET /admin/astro-general-complaints` - List astrologer general complaints
- `PATCH /admin/astro-general-complaints/:problemId/close` - Close a complaint

### Horoscope Management

- `GET /admin/horoscopes` - List horoscopes (with pagination and filters)
- `POST /admin/horoscopes/bulk` - Bulk create horoscopes
- `PATCH /admin/horoscopes/:horoscopeId` - Update a horoscope
- `DELETE /admin/horoscopes/:horoscopeId` - Delete a horoscope

**Note**: All admin endpoints (except login) require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Daily Report Scheduler

The application includes an automated daily report scheduler that:

- Runs every day at **11:55 PM** (23:55:00)
- Fetches data from 12 MongoDB collections based on `createdOn` field for today's records:
  - appointments
  - chat
  - conversionHistory
  - feedback
  - ivrCall
  - orderScore
  - rewards
  - serviceReports
  - user
  - userPayment
  - userProblem
  - videoCall
- Generates CSV files for each collection
- Sends all CSV files as email attachments to `EMAIL_USER` (support@astrosway.com)

The scheduler starts automatically when the application starts. Ensure `EMAIL_USER` and `EMAIL_PASS` are configured in your `.env` file.

### Date Format

Records are filtered based on the `createdOn` field using UTC timezone. The format used is: `2025-05-25T05:25:22.193+00:00`

## Project Structure

```
ADMIN-BE/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── Dockerfile                 # Docker configuration
├── docker-compose.yml         # Docker Compose configuration
├── env.example                # Environment variables example
├── internal/
│   ├── config/                # Configuration management
│   ├── database/              # Database connections (MongoDB)
│   ├── dto/                   # Data Transfer Objects
│   ├── handlers/              # HTTP handlers
│   ├── middleware/            # HTTP middleware (auth, etc.)
│   ├── models/                # Data models
│   ├── repository/            # Data access layer
│   ├── services/              # Business logic layer
│   └── utils/                 # Utility functions
```

## Database

The application uses MongoDB database named **astroway**. Make sure your MongoDB connection string in `MONGO_URI` points to the correct database.

## Docker

The project includes Docker support. You can use Docker Compose to run the application:

```bash
docker-compose up -d
```

Or build and run with Docker:

```bash
docker build -t admin-be .
docker run -p 8082:8082 --env-file .env admin-be
```

## Development

### Running in Development Mode

```bash
go run main.go
```

### Building for Production

```bash
go build -o admin-be main.go
./admin-be
```

## Error Handling

The application includes comprehensive error handling and logging. Check the console output for detailed error messages and logs.

## Security

- JWT-based authentication for all admin endpoints
- Environment variables for sensitive configuration
- Input validation using go-playground/validator
- Secure password storage (use App Passwords for Gmail)

## License

[Add your license information here]

## Support

For issues and questions, please contact the development team.

