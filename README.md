# GNotifier - UQAM Grade Notification System

GNotifier is a lightweight tool that automatically checks for grade updates in your UQAM courses and sends email notifications when new grades are posted.

## Features

- **Automatic monitoring** of your UQAM course grades
- **Email notifications** when grades change
- **Simple setup** with a one-time configuration
- **Secure storage** of previously fetched grades
- **Scheduled checks** via cron job (hourly by default)

## Requirements

- Go (for building from source)
- Gmail account (for sending notifications)
- UQAM student credentials

## Installation

### Option 1: Manual Setup

1. Clone the repository:
   ```
   git clone https://github.com/felixlheureux/uqam-grade-notifier.git
   cd uqam-grade-notifier
   ```

2. Build the application:
   ```
   make build
   ```

3. Edit the configuration file:
   ```
   cp main/config/prod.config.json config.json
   ```
   Update the file with your personal information (see Configuration section).

4. Run the application:
   ```
   ./dist/gnotifier -config ./config.json
   ```

### Option 2: Automatic Installation (Linux)

Run the setup script as root:
```
sudo ./scripts/setup.sh
```

This will:
- Build the application
- Install it to `/usr/local/bin/gnotifier`
- Create necessary configuration directories
- Set up a cron job to run the app every hour

## Configuration

Edit the `config.json` file with your information:

```json
{
  "username": "YOUR_UQAM_USERNAME",
  "password": "YOUR_UQAM_PASSWORD",
  "gmail_app_password": "YOUR_GMAIL_APP_PASSWORD",
  "mail_to": "your.email@example.com",
  "mail_from": "your.gmail@gmail.com",
  "semester": "20251",
  "courses": ["COURSE1/GROUP", "COURSE2/GROUP"],
  "store_path": "/path/to/store/grades.json"
}
```

### Configuration Values

| Field | Description | Example |
|-------|-------------|---------|
| `username` | Your UQAM username | `"SERF12345678"` |
| `password` | Your UQAM password | `"mypassword"` |
| `gmail_app_password` | App password for Gmail (see note below) | `"abcd efgh ijkl mnop"` |
| `mail_to` | Email address to receive notifications | `"your.email@example.com"` |
| `mail_from` | Gmail address used to send notifications | `"your.gmail@gmail.com"` |
| `semester` | Current semester code | `"20251"` (Winter 2025) |
| `courses` | Array of course codes with group numbers | `["INF1132/10", "INF2050/30"]` |
| `store_path` | Path to store grade data | `"/etc/gnotifier/grades.json"` |

**Note on Gmail App Password:**  
You need to generate an app password for your Gmail account:
1. Enable 2-Step Verification on your Google Account
2. Go to [App passwords](https://myaccount.google.com/apppasswords)
3. Select "Mail" and "Other (Custom name)"
4. Generate and copy the 16-character password

## Usage

After configuration, the application will:
1. Authenticate with UQAM
2. Check your grades for each course
3. Compare with previously stored grades
4. Send an email notification if a grade has changed
5. Update the stored grades

When set up with cron, this process runs automatically.

## Structure

- `main/` - Main application code
- `pkg/` - Supporting packages:
    - `alert/` - Email notification functionality
    - `auth/` - UQAM authentication
    - `grade/` - Grade fetching logic
    - `helper/` - Utility functions
    - `store/` - Grade storage management
- `scripts/` - Installation and maintenance scripts

## Security Considerations

- Your UQAM credentials and Gmail app password are stored in plain text in the configuration file
- Make sure to restrict access to this file (e.g., `chmod 600 config.json`)
- Consider using environment variables or a secure password manager for production use

## License

[MIT License]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

Felix L'Heureux (felixslheureux@gmail.com)