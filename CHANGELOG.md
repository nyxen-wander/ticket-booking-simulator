# Changelog

All notable changes to this project will be documented in this file.

## [1.3.0] - 2026-05-28
### Added
- **Administrative CRUD Module**: Introduced a secure, dedicated database management panel (`administration.go`) containing deep record handling functions (`insertMenu`, `updateMenu`, `deleteMenu`) to perform direct data operations over relational structures.
- **Polymorphic Menu Rendering Machine**: Created the structural `menuRenderer` type to centralize terminal painting behaviors, encapsulate UI messaging states, handle layout transformations, and standardize interactive guide text.
- **Constant-Memory Navigation**: Refactored the interactive interface loop away from memory-heavy recursive function calls into flat, optimized `for` loop frameworks utilizing deterministic state breaks and `continue` redirections.

### Changed
- **Unified Initialization Sequencing**: Consolidated the startup pattern in `main.go` to guarantee environmental data registration occurs cleanly in process memory before initializing the global database connection pool.
- **Strict Pointer Propagation**: Updated core utility signatures to pass explicit `*sql.DB` connection frames and `*bufio.Scanner` inputs deep down into nested routing layers, preventing terminal state conflicts and resource footprint spikes.

### Fixed
- **Call-Stack Accumulation Hazard**: Eliminated potential stack-overflow panics triggered under repetitive user typographical validation loops by introducing iterative menus.
- **Redundant Environment Discovery**: Cleaned up duplicated variable sweeps by formalizing a strict single-lifecycle load pattern at the earliest operational phase.

## [1.2.0] - 2026-05-11
### Added
- **Atomic Transactions**: Implemented `db.Begin()`, `Commit()`, and `Rollback()` to ensure that seat updates and booking records are processed as a single, unbreakable unit.
- **Project Modularization**: Refactored the single-file script into a multi-file package (`main.go`, `db_ops.go`, `utils.go`) for better maintainability and navigation.
- **Enhanced Validation**: Added a comprehensive `customerNameAndTicketAmountValidator` helper function to centralize input logic.
- **Robust Error Handling**: Added custom error returns for `clearScreen` and `printBookSummary` to prevent silent failures.

### Changed
- **Code Architecture**: Cleaned up the `main` loop by delegating database operations and utility tasks to dedicated functions.
- **Database Synchronization**: Added a post-transaction seat refresh via `getAvailableSeats` to ensure the terminal always shows the true state of the database.

### Fixed
- **Logical Inversion**: Corrected the validation check in the main loop to properly trigger on error messages.
- **Cleanup Logic**: Improved the terminal clearing function to handle cross-platform errors gracefully.

## [1.1.0] - 2026-05-07
### Added
- **Persistent Booking Summary**: Replaced the local `bookRecord` map with a dedicated `printBookSummary` function that performs a 3-table SQL JOIN to generate reports directly from the database.
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