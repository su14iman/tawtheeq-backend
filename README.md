# tawtheeq-backend

A digital file signing and verification system built with Go and RSA.

## Features

- **Digital file signing** (images, PDF)
- **Signature verification**
- **User management** (login, password change, roles)
- **Team management** (create, add/remove members, delete)
- **File storage** (local or S3/MinIO)
- **RSA key generation**
- **API documentation** via Swagger

---

## Requirements

- Go 1.24+
- MySQL
- Redis
- exiftool
- MinIO (optional for cloud storage)
- Docker (optional)

---

## Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/su14iman/tawtheeq-backend.git
   cd tawtheeq-backend
   ```

2. **Generate RSA keys:**
   ```bash
   ./generate_keys.sh
   ```

3. **Configure environment variables:**
   - Copy `.env.example` to `.env` and update values as needed.

4. **Run the server:**
   ```bash
   go run main.go
   ```
   Or with Docker:
   ```bash
   ./start.sh
   ```

---

## API Documentation

- Available via Swagger at:  
  `http://localhost:8080/swagger/index.html`

---

## API Endpoints

### Authentication

| Method | Endpoint                  | Description                        | Roles Required      |
|--------|---------------------------|------------------------------------|---------------------|
| POST   | `/api/auth/login`         | User login                         | Public              |
| POST   | `/api/auth/forgot-password`| Request password reset             | Public              |
| POST   | `/api/auth/reset-password` | Reset password                     | Public              |

---

### File Signing & Verification

| Method | Endpoint                | Description                        | Roles Required      |
|--------|-------------------------|------------------------------------|---------------------|
| POST   | `/api/upload`           | Upload and digitally sign a file   | Any authenticated   |
| GET    | `/api/verify/:id`       | Verify file signature by ID        | Public              |

---

### User Management

| Method | Endpoint                        | Description                        | Roles Required      |
|--------|---------------------------------|------------------------------------|---------------------|
| GET    | `/api/users/`                   | List all users                     | SuperAdmin          |
| POST   | `/api/users/`                   | Create a new user                  | SuperAdmin, TeamLeader |
| DELETE | `/api/users/:id/remove`          | Remove a user                      | SuperAdmin          |
| PUT    | `/api/users/:id/role`            | Change user role                   | SuperAdmin          |

---

### Profile Management

| Method | Endpoint                        | Description                        | Roles Required      |
|--------|---------------------------------|------------------------------------|---------------------|
| PUT    | `/api/myself/name`              | Update your name                   | Any authenticated   |
| PUT    | `/api/myself/password`          | Change your password               | Any authenticated   |

---

### Team Management

| Method | Endpoint                                 | Description                        | Roles Required      |
|--------|------------------------------------------|------------------------------------|---------------------|
| GET    | `/api/teams/`                            | List all teams                     | SuperAdmin          |
| POST   | `/api/teams/`                            | Create a new team                  | SuperAdmin          |
| DELETE | `/api/teams/:id/remove`                  | Remove a team                      | SuperAdmin          |
| PUT    | `/api/teams/:id/name`                    | Update team name                   | SuperAdmin          |
| PUT    | `/api/teams/:id/leader`                  | Change team leader                 | SuperAdmin          |
| GET    | `/api/teams/:team_id/members`            | List all members in a team         | SuperAdmin          |
| POST   | `/api/teams/members`                     | Add user to a team                 | SuperAdmin          |
| DELETE | `/api/teams/members/:team_id/:user_id`   | Remove user from a team            | SuperAdmin          |

---

### My Team

| Method | Endpoint                                 | Description                        | Roles Required      |
|--------|------------------------------------------|------------------------------------|---------------------|
| GET    | `/api/my/team`                           | Get your team info                 | Any authenticated   |
| GET    | `/api/my/team/members`                   | List members in your team          | TeamLeader          |
| POST   | `/api/my/team/members`                   | Add user to your team              | TeamLeader          |
| DELETE | `/api/my/team/members/:user_id`          | Remove user from your team         | TeamLeader          |

---

### Document Management

| Method | Endpoint                                         | Description                                 | Roles Required      |
|--------|--------------------------------------------------|---------------------------------------------|---------------------|
| GET    | `/api/documents/visible`                         | List all visible documents                  | SuperAdmin          |
| GET    | `/api/documents/hidden`                          | List all hidden documents                   | SuperAdmin          |
| GET    | `/api/documents/team/:team_id/visible`           | List visible documents for a team           | SuperAdmin          |
| GET    | `/api/documents/team/:team_id/hidden`            | List hidden documents for a team            | SuperAdmin          |
| GET    | `/api/documents/user/:user_id/visible`           | List visible documents for a user           | SuperAdmin          |
| GET    | `/api/documents/user/:user_id/hidden`            | List hidden documents for a user            | SuperAdmin          |
| GET    | `/api/documents/user/me`                         | List your visible documents                 | Any authenticated   |
| GET    | `/api/documents/user/me/hidden`                  | List your hidden documents                  | SuperAdmin          |
| GET    | `/api/documents/:id/hide`                        | Hide a document (SuperAdmin)                | SuperAdmin          |
| GET    | `/api/documents/:id/show`                        | Unhide a document (SuperAdmin)              | SuperAdmin          |
| GET    | `/api/documents/my`                              | List your visible documents                 | Any authenticated   |
| GET    | `/api/documents/myteam`                          | List all documents for your team            | TeamLeader          |
| GET    | `/api/documents/myteam/:id/hide`                 | Hide a document from your team              | TeamLeader          |
| GET    | `/api/documents/my/:id/hide`                     | Hide a document from your own list          | Any authenticated   |

---

**Note:**  
- `SuperAdmin` and `TeamLeader` are user roles with different permissions.
- `Any authenticated` means any logged-in user.

---

## Contribution

Contributions are welcome! Please open an Issue or Pull Request for suggestions or improvements.

---

## License

[MIT](LICENSE)