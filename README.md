# GNotifier - Grade Notification System for UQAM Students

GNotifier is a lightweight Go application that automatically monitors your UQAM course grades and sends you email notifications when new grades are posted.

## 📋 Features

- **Automated Grade Checking:** Regularly polls the UQAM portal for your course grades
- **Instant Notifications:** Sends email alerts when grades change
- **Minimal Configuration:** Simple one-time setup process
- **Secure Data Storage:** Locally stores grade history
- **Background Operation:** Runs silently via scheduled tasks

## 🔧 Requirements

- Go programming language (for building from source)
- Gmail account (for sending notifications)
- UQAM student credentials

## 🚀 Installation

### Option 1: Quick Setup (Linux)

Run our automated installation script:

```bash
sudo ./scripts/setup.sh
```

This script will:
- Install Go and required dependencies
- Build the application
- Create the necessary directories
- Configure basic firewall settings
- Set up instructions for scheduling hourly checks

### Option 2: Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/felixlheureux/uqam-grade-notifier.git
   cd uqam-grade-notifier
   ```

2. Build the application:
   ```bash
   cd main
   make build
   ```

3. Create your configuration file:
   ```bash
   cp main/config/prod.config.json config.json
   ```

4. Edit the configuration with your personal information

5. Run the application:
   ```bash
   ./dist/gnotifier -config ./config.json
   ```

## ⚙️ Configuration

Edit the `config.json` file with your personal information:

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

### Configuration Options Explained

| Option | Description | Example |
|--------|-------------|---------|
| `username` | Your UQAM username | `"SERF12345678"` |
| `password` | Your UQAM password | `"mypassword"` |
| `gmail_app_password` | App-specific password for your Gmail account | `"abcd efgh ijkl mnop"` |
| `mail_to` | Email where you want to receive notifications | `"your.email@example.com"` |
| `mail_from` | Gmail address used to send notifications | `"your.gmail@gmail.com"` |
| `semester` | Current semester code | `"20251"` (Winter 2025) |
| `courses` | List of courses and group numbers to monitor | `["INF1132/10", "INF2050/30"]` |
| `store_path` | Location to save grade history | `"/home/user/app/grades.json"` |

### Gmail App Password Setup

To send email notifications, you'll need to generate an app password:

1. Enable 2-Step Verification on your Google Account
2. Visit [Google App Passwords](https://myaccount.google.com/apppasswords)
3. Select "Mail" and "Other (Custom name)"
4. Use the generated 16-character password in your configuration

## 📐 How It Works

When run, GNotifier performs the following operations:

1. Authenticates with your UQAM credentials
2. Fetches the current grades for each of your courses
3. Compares these grades with previously stored values
4. If changes are detected, sends an email notification
5. Updates the stored grade data for future comparisons

For best results, set up a recurring scheduled task (cron job) to run the application regularly.

## 📁 Project Structure

- `main/` - Application entry point and configuration
- `pkg/` - Core functionality packages:
    - `alert/` - Email notification handling
    - `auth/` - UQAM authentication module
    - `grade/` - Grade fetching logic
    - `helper/` - Utility functions
    - `store/` - Grade storage management
- `scripts/` - Installation and maintenance scripts

## 🔒 Security Notes

- Your UQAM credentials and Gmail app password are stored as plain text in the configuration file
- Restrict access to this file with appropriate permissions (`chmod 600 config.json`)
- Consider using environment variables for sensitive information in production environments

## 🤝 Contributing

Contributions are welcome! Feel free to submit pull requests or open issues for bugs and feature requests.

## 📄 License

This project is licensed under the MIT License.

## 👤 Author

Felix L'Heureux (felixslheureux@gmail.com)