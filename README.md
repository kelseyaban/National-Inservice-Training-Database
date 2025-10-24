# National In-Service Training Database API

The **National In-Service Training Database** is a RESTful API designed for the **Belize Police Department** to manage officer training data, monitor participation, generate reports, and analyze training trends.  
This system supports role-based access, comprehensive reporting, and analytics to ensure officers meet national training requirements(to be impemented).

---

## Features

### Account Privilege
**Administrator**
Full access to add, edit, and manage users, training data, and reports.
**Content Contributor**
Create and manage training courses, record participant data, and print or export training records.
**System User**
Access and download personal or unit training reports.

---

## Endpoints

### Healthcheck
- **GET** `/v1/healthcheck` – Check API status

### Users
- **POST** `/v1/users` – Register user  
- **PUT** `/v1/users/activated` – Activate user  
- **POST** `/v1/tokens/authentication` – User login/authentication  
- **PATCH** `/v1/users/update/:id` – Update user info  
- **GET** `/v1/users/details` – List users  
- **DELETE** `/v1/users/delete/:id` – Delete user  
- **PATCH** `/v1/users/update-password/:id` – Update password  

### Roles
- **POST** `/v1/roles` – Create role  
- **GET** `/v1/roles/:id` – View role  
- **PATCH** `/v1/roles/:id` – Update role  
- **DELETE** `/v1/roles/:id` – Delete role  
- **GET** `/v1/roles` – List roles  

### User Roles
- **POST** `/v1/users/assign-role` – Assign role to user  
- **GET** `/v1/users/user_roles/:id` – View user roles  
- **PATCH** `/v1/users/update-role/:id` – Update user role  
- **DELETE** `/v1/users/delete-role/:id` – Remove user role  
- **GET** `/v1/users/user_roles` – List all user-role mappings  

### Facilitator Ratings
- **POST** `/v1/facilitator-rating` – Add rating  
- **GET** `/v1/facilitator-rating/:id` – View rating  
- **GET** `/v1/facilitator-rating` – List ratings  

### Courses
- **POST** `/v1/courses` – Create course  
- **GET** `/v1/courses/:id` – View course  
- **PATCH** `/v1/courses/:id` – Update course  
- **DELETE** `/v1/courses/:id` – Delete course  
- **GET** `/v1/courses` – List courses  

### Course Postings
- **POST** `/v1/course/posting` – Create posting  
- **GET** `/v1/course/posting/:id` – View posting  
- **PATCH** `/v1/course/posting/:id` – Update posting  
- **DELETE** `/v1/course/posting/:id` – Delete posting  
- **GET** `/v1/course/posting` – List postings  

### Sessions
- **POST** `/v1/session` – Create session  
- **GET** `/v1/session/:id` – View session  
- **PATCH** `/v1/session/:id` – Update session  
- **DELETE** `/v1/session/:id` – Delete session  
- **GET** `/v1/session` – List sessions  

### User Sessions
- **POST** `/v1/user_session` – Create user session  
- **GET** `/v1/user_session/:id` – View user session  
- **PATCH** `/v1/user_session/:id` – Update user session  
- **DELETE** `/v1/user_session/:id` – Delete user session  
- **GET** `/v1/user_session` – List user sessions  

### Attendance
- **POST** `/v1/attendance` – Create attendance record  
- **GET** `/v1/attendance/:id` – View individual attendance  
- **PATCH** `/v1/attendance/:id` – Update attendance


## Future Updates

### Data Analysis Functions

- Calculate **total officers trained** and **percentage trained** by:
  - Day, week, month, quarter, mid-year, or annually  
  - Region, Formation, or Unit/Branch
- Group training by **mandatory or elective topics**
- Generate statistics for **national**, **regional**, and **unit-level** participation
- Identify or flag officers who **do not meet annual training hour requirements**
- Visual outputs include **graphs**, **bar charts**, and **pie charts**

---

### Report Capabilities

- Daily, Weekly, Monthly, Quarterly, and Yearly Reports  
- Individual **Report Cards/Transcripts**  
- **Regional / Formation / Unit** Training Reports  

---

## Tech Stack

- **Language:** Golang 
- **Database:** PostgreSQL  
- **Migrations:** `golang-migrate`  
- **API Framework:** `net/http` or `httprouter`  
---

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/kelseyaban/National-Inservice-Training-Database.git
cd National-Inservice-Training-Database

```
### 2. Creation of Database
```bash
# Create the application database
CREATE DATABASE national_inservice_training;

# Create a dedicated user for the application
CREATE USER facilitator WITH PASSWORD 'your_password_here';

# Grant privileges to the new user
ALTER DATABASE national_inservice_training OWNER TO facilitator;

# Enable the citext extension to make email comparisons case-insensitive:
CREATE EXTENSION IF NOT EXISTS citext;

# Run migrations to set up the database tables
db/migrations/up
```

## Sample Requests 

This section provides sample `CURL` commands for testing each endpoint in the **National Inservice Training Database API**.

## Users
### Create User
```bash
BODY='{
  "regulation_number": "R14345",
  "username": "kelseyaban",
  "fname": "Kelsey",
  "lname": "Aban",
  "email": "kelsey@example.com",
  "gender": "F",
  "formation": 4,
  "rank": 2,
  "postings": 5,
  "password": "kelsey123"
}'
curl -d "$BODY" localhost:4000/v1/users
```
### Activate User
```bash
curl -X PUT -d '{"token": "3EX3UPDHJMX5SUBRWOZTPLMDGM"}' localhost:4000/v1/users/activated
```
### Read Users
```bash
curl -i "localhost:4000/v1/users?page=1&page_size=2"
```
### Update User
```bash
curl -X PATCH http://localhost:4000/v1/users/update/1 \
-H "Content-Type: application/json" \
-d '{
  "regulation_number": "REG12345",
  "username": "newusername",
  "fname": "John",
  "lname": "Doe",
  "email": "john.doe@example.com",
  "formation": 2,
  "rank": 5,
  "postings": 3
}'
```
### Update Password
```bash
curl -X PATCH http://localhost:4000/v1/users/update-password/4 \
-H "Content-Type: application/json" \
-d '{"new_password": "myNewSecret123"}'
```
### Delete User
```bash
curl -X DELETE localhost:4000/v1/users/delete/1
```
## Roles
### Create Role
```bash
BODY='{"role": "Manager"}'
curl -d "$BODY" localhost:4000/v1/roles
```
### Read Roles
```bash
# One role
curl -i localhost:4000/v1/roles/1

# All roles
curl -i localhost:4000/v1/roles
```

### Update Role
```bash
curl -X PATCH -d '{"role": "Supervisor"}' localhost:4000/v1/roles/5
```
### Delete Role
```bash
curl -X DELETE localhost:4000/v1/roles/5
```
## User Roles
### Assign Roles to User
```bash
BODY='{"user_id": 1, "role_ids": [1,2]}'
curl -d "$BODY" localhost:4000/v1/users/assign-role
```
### Read User Roles
```bash
# One user’s roles
curl -i localhost:4000/v1/users/user_roles/1

# All user-role mappings
curl -i localhost:4000/v1/users/user_roles
```
### Update User Role
```bash
BODY='{"old_role_id": 2, "new_role_id": 3}'
curl -X PATCH -d "$BODY" localhost:4000/v1/users/update-role/1
```
### Delete User Role
```bash
BODY='{"role_id": 3}'
curl -X DELETE -d "$BODY" localhost:4000/v1/users/delete-role/1
```
## Facilitator Rating
### Create Rating
```bash
BODY='{"user_id": 2, "rating": 5}'
curl -d "$BODY" localhost:4000/v1/facilitator-rating
```
### Read Ratings
```bash
# One record
curl -i localhost:4000/v1/facilitator-rating/1

# All records
curl -i localhost:4000/v1/facilitator-rating
```
## Courses
### Create Course
```bash
BODY='{
  "course": "Narcotic Detection",
  "description": "Trains the dog to identify and locate explosive materials in various environments."
}'
curl -d "$BODY" localhost:4000/v1/courses
```
### Read Courses
```bash
curl -i "localhost:4000/v1/courses?page=1&page_size=2"
curl -i "localhost:4000/v1/courses?sort=-id"
```

### Update Course
```bash
curl -X PATCH -H "Content-Type: application/json" \
-d '{"description": "Searching for substances."}' \
http://localhost:4000/v1/courses/2
``` 
###Delete Course
```bash
curl -X DELETE localhost:4000/v1/courses/1
```
## Course Posting
### Create Course Posting
```bash
BODY='{
  "course_id": 1,
  "posting_id": 2,
  "mandatory": true,
  "credithours": 40,
  "rank_id": 3
}'
curl -d "$BODY" localhost:4000/v1/course/posting
```
### Read Course Postings
```bash
curl -i "localhost:4000/v1/course/postings?page=1&page_size=2"
```
### Update Course Posting
```bash
BODY='{
  "course_id": 4,
  "posting_id": 3,
  "mandatory": true,
  "credithours": 4,
  "rank_id": 1
}'
curl -X PATCH -d "$BODY" localhost:4000/v1/course/posting/1
```
### Delete Course Posting
```bash
curl -X DELETE localhost:4000/v1/course/postings/1
```
## Session
### Create Session
```bash
BODY='{"course_id": 1, "formation_id": 2, "facilitator_id": 3}'
curl -d "$BODY" localhost:4000/v1/session
```
### Read Sessions
```bash
curl -i localhost:4000/v1/session/1
curl -i localhost:4000/v1/session
```
### Update Session
```bash
BODY='{"facilitator_id": 5}'
curl -X PATCH -d "$BODY" localhost:4000/v1/session/1
```
### Delete Session
```bash
curl -X DELETE localhost:4000/v1/session/1
```
## User Session
### Create User Session
```bash
BODY='{
  "session_id": 2,
  "credithours_completed": 4,
  "grade": "B",
  "feedback": "Good performance",
  "trainee_id": 4
}'
curl -d "$BODY" localhost:4000/v1/user_session
``` 
### Read User Sessions
```bash
curl -i localhost:4000/v1/user_session/5
curl -i localhost:4000/v1/user_session
``` 
### Update User Session
```bash
BODY='{
  "credithours_completed": 4,
  "grade": "A+",
  "feedback": "Excellent improvement"
}'
curl -X PATCH -d "$BODY" localhost:4000/v1/user_session/5
```
### Delete User Session
```bash
curl -X DELETE localhost:4000/v1/user_session/4
```
## Attendance
### Create Attendance Record
```bash
BODY='{"user_session_id": 1, "attendance": true, "date": "2025-10-19"}'
curl -d "$BODY" localhost:4000/v1/attendance
```
### Read Attendance
```bash
curl -i "localhost:4000/v1/attendance?page=1&page_size=2"
```
### Update Attendance
```bash
BODY='{"attendance": false, "date": "2025-10-20"}'
curl -X PATCH -d "$BODY" localhost:4000/v1/attendance/2
```
## Authentication Example

Use authentication to generate a token for protected routes.
### Create Authentication Token
```bash
BODY='{"email": "kelsey@example.com", "password": "kelsey123"}'
curl -i -d "$BODY" localhost:4000/v1/tokens/authentication
``` 
Then use the token in subsequent requests:
```bash
curl -i -H "Authorization: Bearer YOUR_TOKEN_HERE" localhost:4000/v1/session
```
