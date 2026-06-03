# Ticket Booking Simulator

Backend-driven cinema management system built with **Go** and **PostgreSQL**. This project serves as the core engine for a multi-theater ticket simulator, featuring strict data integrity, physical coordinate-based seat mapping, and an interactive administrative panel.

## Overview
The system manages a "Session-First" workflow, allowing operators to orchestrate movie screenings, manage visual seating matrices, run secure database lookups, and execute live customer ticket bookings while dynamically filtering out theaters under maintenance or inactive sessions. 

The architecture acts as a permanent ledger, guaranteeing that booked seats are permanently archived for accounting purposes while dynamically freeing up database storage by pruning unbooked inventory.

### Key Features
- **Coordinate-Based Seat Matrix:** Replaces generic integer counters with a 1-to-1 physical grid map (e.g., A1, B4), dynamically generated and rendered in the terminal based on theater capacity.
- **Automated Inventory Lifecycle:** Features automated blueprint cloning (`sessionSeatsRegeneration`) when sessions are published, and a Safe Bulk Delete trigger that aggressively clears unused database bloat while permanently protecting sold ticket records.
- **Relational Data Model:** Decoupled multi-table schema using `movies`, `theaters`, `physical_seats`, `sessions`, `session_seats`, and `bookings` bound by strict foreign keys and composite unique constraints.
- **Full Administrative CRUD Module:** Isolated management workspace enabling secure create, read, update, and delete operations across all system tables without disrupting active booking workflows.
- **Polymorphic Menu Rendering Engine (`menuRenderer`):** Dedicated structural state machine handling centralized terminal clearing, header rendering, custom collection layouts, active debugging states, and live database string transformations.
- **Flat-Memory Iterative Routing:** Eliminates deep call-stack accumulation by transitioning from recursive UI redraw mechanics into deterministic, constant-memory `for` loop menus utilizing controlled `return` and `continue` keywords.
- **ACID Transaction Safeguard:** Guarantees absolute data atomicity across high-concurrency ticket sales by anchoring specific seat coordinates via safe `db.Begin()`, row-level `WHERE is_booked = false` locks, and `Commit()` boundaries.

## System Flow
1. **Bootstrap Phase**: The system loads environmental variables into memory space via a single unified loader before spinning up a highly optimized global connection pool via `pgx`.
2. **Interactive Selection & Interface Layer**: Operators are routed through an optimized constant-memory main loop managed by the centralized `menuRenderer`. 
3. **Admin Controls (CRUD Workspace)**: Authorized administrators can cleanly execute transactional database alterations across tables, utilizing primary key constraints as immutable records anchors.
4. **Ticket Booking Engine**: Real-time ticket transactions render a live matrix grid of the theater. The transaction block securely isolates the selected grid coordinate, updates the physical state, logs the receipt history, and updates the terminal view dynamically.

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

### Table: physical_seats (The Blueprint)
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key |
| theater_id | INTEGER | Foreign Key referencing `theaters(id)` |
| row_letter | VARCHAR(2) | The physical row designation (e.g., 'A', 'B') |
| seat_num | INTEGER | The physical seat number within the row |

### Table: sessions
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key (Immutable Identity Anchor) |
| movie_id | INTEGER | Foreign Key referencing `movies(id)` |
| theater_id | INTEGER | Foreign Key referencing `theaters(id)` |
| show_time | TIMESTAMP | The specific date/time making this session a unique event |
| is_active | BOOLEAN | Draft/Publish toggle (Triggers Seat Auto-Generation/Deletion) |

### Table: session_seats (Live Inventory)
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key |
| session_id | INTEGER | Foreign Key referencing `sessions(id)` |
| physical_seat_id | INTEGER | Foreign Key referencing `physical_seats(id)` |
| is_booked | BOOLEAN | Live status flag with strict `false` state locking |
*Note: Enforces a `UNIQUE (session_id, physical_seat_id)` constraint to prevent duplicate chairs.*

### Table: bookings (The Ledger)
| Column | Type | Description |
|---|---|---|
| id | SERIAL | Primary Key (Immutable Identity Anchor) |
| session_id | INTEGER | Foreign Key referencing `sessions(id)` |
| session_seat_id | INTEGER | Foreign Key referencing `session_seats(id)` |
| customer_name | VARCHAR(50) | Normalized and sanitized booker identity |
| created_at | TIMESTAMP | Database-generated audit timestamp |

## Getting Started

### 1. Prerequisites
Ensure you have the following installed on your environment (WSL2 preferred):
- Go (version 1.21 or higher)
- PostgreSQL server instance running locally or remotely

### 2. Database Initialization
Log into your PostgreSQL shell and create the target database schema. Ensure foreign key constraints and composite unique keys are properly linked to maintain data integrity as documented in the architecture segment.

### 3. Environment Configuration
Create a `.env` file in the project root directory to store your secure connection string and administrative credentials:
```env
DB_URL=postgres://your_user:your_password@localhost:5432/ticket_booking_db
ADMIN_PASSWD=your_secure_admin_password
```

### 4. Installation & Execution
# Fetch dependencies
```bash
go mod tidy
```

# Run the simulator
```bash
go run .
```
