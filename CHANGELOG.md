# Changelog

All notable changes to this project will be documented in this file.

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