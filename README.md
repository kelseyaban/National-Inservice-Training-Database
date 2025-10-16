# 🛡️ National In-Service Training Database API

The **National In-Service Training Database** is a RESTful API designed for the **Belize Police Department** to manage officer training data, monitor participation, generate reports, and analyze training trends.  
This system supports role-based access, comprehensive reporting, and analytics to ensure officers meet national training requirements(to be impemented).

---

## 🛠️ Features

### 👤 Account Privilege
**Administrator**
Full access to add, edit, and manage users, training data, and reports.
**Content Contributor**
Create and manage training courses, record participant data, and print or export training records.
**System User**
Access and download personal or unit training reports.

---

## 📊 Data Analysis Functions

- Calculate **total officers trained** and **percentage trained** by:
  - Day, week, month, quarter, mid-year, or annually  
  - Region, Formation, or Unit/Branch
- Group training by **mandatory or elective topics**
- Generate statistics for **national**, **regional**, and **unit-level** participation
- Identify or flag officers who **do not meet annual training hour requirements**
- Visual outputs include **graphs**, **bar charts**, and **pie charts**

---

## 🛠️ Endpoints

Method	Endpoint
Get     /v1/healthcheck             Confirm that your application and its dependencies are running properly
POST	/v1/users	                Register new system users
POST	/v1/tokens/authentication	Authenticate user login
POST	/v1/training-courses	    Create a new training course
GET     /v1/training-courses	    View all courses and participants
GET     /v1/reports/summary	        Generate analytical reports (daily,monthly etc.)
GET     /v1/reports/officers	    View officers’ training status and compliance
GET     /v1/reports/unit	        Regional/unit training summaries
GET     /v1/reports/cards/:id	    Individual officer transcripts

## 🧾 Report Capabilities

- Daily, Weekly, Monthly, Quarterly, and Yearly Reports  
- Individual **Report Cards/Transcripts**  
- **Regional / Formation / Unit** Training Reports  

---

## 🛠️ Tech Stack

- **Language:** Golang 
- **Database:** PostgreSQL  
- **Migrations:** `golang-migrate`  
- **API Framework:** `net/http` or `httprouter`  
---

## ⚙️ Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/kelseyaban/National-Inservice-Training-Database.git
cd National-Inservice-Training-Database
