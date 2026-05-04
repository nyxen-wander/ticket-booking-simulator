# Ticket Booking Simulator

Backend-driven cinema management system built with **Go** and **PostgreSQL**. This project serves as the core engine for a multi-theater ticket simulator, featuring strict data integrity and relational mapping.

## Overview
The system manages the "Session-First" workflow, allowing employees to select from currently active movie showings while automatically filtering out theaters under maintenance or inactive sessions.

### Key Features
- **Relational Data Model:** Decoupled architecture using `movies`, `theaters`, `sessions`, and `bookings` tables.
- **Maintenance Guard:** Automated filtering of sessions where the physical theater room is marked as inactive (`t.is_active = 't'`).
- **Dynamic Session Selection:** Interactive CLI allows users to view all survivors of the "active" filters before proceeding to booking.
- **Secure Configuration:** Zero-footprint credential management using `.env` files and `godotenv`.
- **Input Sanitization:** Uses `regexp` for customer name validation and `strconv` for ticket parsing to prevent runtime crashes.

## System Flow
1. **Bootstrap**: Loads `.env` and establishes a PostgreSQL connection via `pgx`.
2. **Session Discovery**: Queries a 3-table join to find available showtimes.
3. **Selection**: Validates the chosen `session_id` against the active pool.
4. **Transaction Logic**:
    - Updates `sessions.available_seats` to reflect the booking.
    - Inserts a record into the `bookings` table for historical tracking.
5. **Summary**: Generates a local booking record for the current work session.

## Tech Stack
- **Language:** Go (Golang)
- **Database:** PostgreSQL (Relational)
- **Library:** `github.com/jackc/pgx/v5`, `github.com/joho/godotenv`
- **Environment:** WSL2 / Windows

## Getting Started

### 1. Database Setup
Ensure your PostgreSQL instance has the following tables with appropriate foreign key constraints:

- `movies`
- `theaters`
- `sessions`
- `bookings`

Below is the schema reference for each:

#### Table: movies

Stores the catalog of available films.

|Column|Type|Description|
|--------|---|--------------|
|id|SERIAL|	Primary Key.|
|movie_title|VARCHAR(50)|The title of the film.|
|rating|NUMERIC|Critical rating (e.g., 6.9). |
|duration|INTEGER|Run time in minutes. |

#### Table: theaters

Manages the physical rooms and their operational status.

Column	Type	Description
|id	|SERIAL	|Primary Key.|
|---|-------|------------|
|total_capacity|INTEGER|Total seats in the room.|
|is_active|BOOLEAN|Maintenance toggle (filters sessions if false). |

#### Table: sessions (The Core Engine)

The bridge between movies and theaters. This is the primary target for bookings.  

Column	Type	Description
|id|SERIAL|Primary Key.  |
|---|---|---|
|movie_id|INTEGER|Foreign Key referencing movies(id).  |
|theater_id|INTEGER|Foreign Key referencing theaters(id).  |
|available_seats|INTEGER|Real-time count of remaining tickets. |
|is_active|BOOLEAN|Scheduling toggle for specific showtimes.  |

#### Table: bookings

Audit trail for every transaction processed by the system.  

|Column|Type|Description|
|---|---|---|
|id|SERIAL|Primary Key.  |
|session_id|INTEGER|Foreign Key referencing sessions(id).  |
|customer_name|VARCHAR(50)|Sanitized name of the booker.  |
|seat_count|INTEGER|Number of tickets purchased.  |
|created_at|TIMESTAMP|Auto-generated timestamp of the transaction.|

### 2. Environment Configuration
Create a `.env` file in the project root that is containing:

    DB_URL=postgres://user:password@localhost:5432/ticket_booking_db

### 3. Execution
Run in your terminal:

    go run main.go

## Future Roadmap

    [ ] ACID Compliance: Implement db.Begin() for atomic transactions to prevent data desync.

    [ ] Web API: Transition from CLI input to a RESTful JSON interface.

    [ ] GUI Evolution: Build an on-site dashboard for employees using Go-based UI toolkits.
