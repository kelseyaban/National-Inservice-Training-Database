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

## Sample Requests 

