# Ticket Booking Simulator

Backend-driven cinema management system built with **Go** and **PostgreSQL**. This project serves as the core engine for a multi-theater ticket simulator, featuring strict data integrity, relational mapping, and an interactive administrative panel.

## Overview
The system manages a "Session-First" workflow, allowing operators to orchestrate movie screenings, manage seating capacities, run secure database lookups, and execute live customer ticket bookings while dynamically filtering out theaters under maintenance or inactive sessions.

### Key Features
- **Relational Data Model:** Decoupled multi-table schema using `movies`, `theaters`, `sessions`, and `bookings` tables bound by foreign keys.
- **Full Administrative CRUD Module:** Isolated management workspace enabling secure create, read, update, and delete operations across all system tables without disrupting active booking workflows.
- **Polymorphic Menu Rendering Engine (`menuRenderer`):** Dedicated structural state machine handling centralized terminal clearing, header rendering, custom collection layouts, active debugging states, and live database string transformations.
- **Flat-Memory Iterative Routing:** Eliminates deep call-stack accumulation by transitioning from recursive UI redraw mechanics into deterministic, constant-memory `for` loop menus utilizing controlled `return` and `continue` keywords.
- **ACID Transaction Safeguard:** Guarantees absolute data atomicity across high-concurrency ticket sales by anchoring seat updates and booking logs via safe `db.Begin()`, `Commit()`, and `Rollback()` boundaries.
- **Automated Input Normalization & Security:** Live regular expression parsing and input normalization prevent buffer collisions or data mismatches during CLI operations.

## System Flow
1. **Bootstrap Phase**: The system loads environmental variables into memory space via a single unified loader before spinning up a highly optimized global connection pool via `pgx`.
2. **Interactive Selection & Interface Layer**: Operators are routed through an optimized constant-memory main loop managed by the centralized `menuRenderer`. 
3. **Admin Controls (CRUD Workspace)**: Authorized administrators can cleanly execute transactional database alterations across tables, utilizing primary key constraints as immutable records anchors.
4. **Ticket Booking Engine**: Real-time ticket transactions are managed through isolation blocks that securely evaluate inventory states, alter available slots, log receipt histories, and update the view dynamically.

## Tech Stack
- **Language:** Go (Golang)
- **Database:** PostgreSQL (Relational Engine)
- **Driver / Libraries:** `github.com/jackc/pgx/v5` (Standard library compatibility layer), `github.com/joho/godotenv`
- **Security / Terminal Modules:** `golang.org/x/term` (Secure password masking)
- **Environment:** Linux (Ubuntu/WSL2) / Windows

## Database Architecture & Structural Integrity

### Table: movies
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key (Immutable Identity Anchor) |
| movie_title | VARCHAR(100) | Full title of the film |
| rating | VARCHAR(5) | Content classification rating (e.g., SU, R) |
| duration | INTERVAL | Total running length of the movie |

### Table: theaters
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key (Immutable Identity Anchor) |
| total_capacity | INTEGER | Total structural seats within the physical room |
| is_active | BOOLEAN | Operational maintenance toggle |

### Table: sessions
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key (Immutable Identity Anchor) |
| movie_id | INTEGER | Foreign Key referencing `movies(id)` |
| theater_id | INTEGER | Foreign Key referencing `theaters(id)` |
| available_seats | INTEGER | Live tracking capacity of remaining ticket count |
| is_active | BOOLEAN | Scheduling toggle for custom availability |

### Table: bookings
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key (Immutable Identity Anchor) |
| session_id | INTEGER | Foreign Key referencing `sessions(id)` |
| customer_name | VARCHAR(50) | Normalized and sanitized booker identity |
| seat_count | INTEGER | Total volume of tickets claimed |
| created_at | TIMESTAMP | Database-generated audit timestamp |

## Getting Started

### 1. Prerequisites
Ensure you have the following installed on your environment (WSL2 preferred):
- Go (version 1.21 or higher)
- PostgreSQL server instance running locally or remotely

### 2. Database Initialization
Log into your PostgreSQL shell and create the target database schema. Ensure foreign key constraints are properly linked to maintain data integrity as documented in the architecture segment.

### 3. Environment Configuration
Create a `.env` file in the project root directory to store your secure connection string and administrative credentials:
```env
DB_URL=postgres://your_user:your_password@localhost:5432/ticket_booking_db
ADMIN_PASSWD=your_secure_admin_password
```

### 4. Installation & Execution
Clone the repository, fetch the necessary module dependencies, and launch the application:
```Bash
# Fetch dependencies
go mod tidy

# Run the simulator
go run .
```

## Future Roadmap

    [ ] Advanced Seat Mapping: Introduce physical coordinate arrays (e.g., A1, B5) mapped via a localized theater inventory layout table.

    [ ] RESTful Web API Interface: Refactor the core operational logic away from the terminal layout into a decoupled JSON API utilizing a framework like Gin or Fiber.

    [ ] GUI Evolution: Build an on-site dashboard for employees using Go-based UI toolkits.
