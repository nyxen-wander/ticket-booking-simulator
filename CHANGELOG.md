# Changelog

All notable changes to this project will be documented in this file.

## [1.4.0] - 2026-06-03
### Added
- **Coordinate-Based Seat Matrix**: Implemented `physical_seats` as an immutable structural blueprint and `session_seats` for live inventory, establishing a 1-to-1 physical grid map (e.g., A1, B4) across theaters.
- **Dynamic Terminal Painter**: Added matrix rendering logic to visually display theater layouts and live booking statuses (e.g., `[ A1 ]`, `[ XX ]`) directly within the CLI.
- **Automated Inventory Lifecycle**: Introduced `sessionSeatsRegeneration` to auto-clone blueprint seats upon session activation and `bulkDeleteSessionSeats` to safely purge unbooked database bloat when a session is deactivated.

### Changed
- **Atomic Booking Transactions**: Refactored the `ticketBooking` engine to perform row-level locks utilizing strict `WHERE is_booked = false` boundaries and `RETURNING id`, guaranteeing absolute concurrency safety during simultaneous purchases.
- **Ledger-Based Receipts**: Updated the `bookings` schema to rely on a strict foreign key relationship with specific physical coordinates (`session_seat_id`) instead of integer subtraction, permanently preserving historical sales data.
- **Session State Management**: Transitioned the `is_active` flag into a strict draft/publish trigger, maximizing database storage efficiency by delaying live seat generation until a session officially opens.

## [1.3.0] - 2026-05-28
### Added
- **Administrative CRUD Module**: Introduced a secure, dedicated database management panel (`administration.go`) containing deep record handling functions (`insertMenu`, `updateMenu`, `deleteMenu`) to perform direct data operations over relational structures.
- **Polymorphic Menu Rendering Machine**: Created the structural `menuRenderer` type to centralize terminal painting behaviors, encapsulate UI messaging states, handle layout transformations, and standardize interactive guide text.
- **Constant-Memory Navigation**: Refactored the interactive interface loop away from memory-heavy recursive function calls into flat, optimized `for` loop frameworks utilizing deterministic state breaks and `continue` redirections.

### Changed
- **Unified Initialization Sequencing**: Consolidated the startup pattern in `main.go` to guarantee environmental data registration occurs cleanly in process memory before initializing the global database connection pool.
- **Strict Pointer Propagation**: Updated core utility signatures to pass explicit `*sql.DB` connection frames and `*bufio.Scanner` inputs deep down into nested routing layers, preventing terminal state conflicts and resource footprint spikes.

### Fixed
- **Call-Stack Accumulation Hazard**: Eliminated potential stack-overflow panics triggered under repetitive user typographical validation loops by introducing the iterative routing loop.

## [1.2.0] - 2026-05-15
### Added
- **Dynamic Data Formatting**: Implemented `CAST` and `TO_CHAR` functions inside SQL queries to format data output directly at the database layer before reaching Go.
- **Summary Export**: Added an administrative function to generate reports directly from the database.
- **Dynamic Session Switching**: Added a `list` command within the main booking loop, allowing employees to switch active sessions without restarting the application.
- **Enhanced Time Handling**: Integrated `time.Format()` for standardized UTC/Local time strings, replacing manual string splitting for more reliable database filtering.
- **Modular Data Fetching**: Extracted seat availability logic into a standalone `getAvailableSeats` function to improve code readability and reuse.

### Changed
- **Database Query Optimization**: Updated the summary query to use PostgreSQL's `to_char` for consistent date formatting in the terminal output.
- **Input Normalization**: Switched command checks (like `exit` and `list`) to be case-insensitive for a smoother user experience.
- **Workflow Improvements**: Added a post-loop prompt to selectively print the booking summary rather than printing it automatically.

### Fixed
- Improved error handling in the `main` scanner to prevent silent failures during user input.

## [1.0.0] - 2026-04-28
### Added
- **Relational Schema**: Initial implementation of `movies`, `theaters`, `sessions`, and `bookings` tables.
- **Core Booking Logic**: Basic loop for seat subtraction and record insertion using `pgx/v5`.
- **Environment Management**: Integration of `.env` files via `godotenv` for secure credential handling.
- **Validation**: Regex-based name sanitization and basic error messaging.
