# SCForm Notes

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/99XCDj?referralCode=-nFAKR)

## Project Purpose

SCForm Notes is a helper application designed to extract and manage grade information from GALIA, a learning management system/ERP used by educational institutions. This tool addresses specific limitations of GALIA, which does not provide:

- A way to download or export grades
- Calculation of a global GPA (Grade Point Average)

This application allows students to efficiently access, track, and analyze their academic performance data outside of GALIA's native interface.

[Version française](README.fr.md)

A web application built with Go (Fiber) and modern frontend technologies for managing and processing form data. The application uses HTMX for dynamic interactions and Tailwind CSS for styling.

## 🚀 Features

- Modern web interface with Tailwind CSS and DaisyUI
- Server-side rendering with Go templates
- Dynamic UI updates using HTMX
- Form data processing and management
- Asset management system
- Environment configuration support
- Live reload development with Air

## 🛠 Tech Stack

### Backend
- Go (Fiber web framework)
- HTML Templates
- Environment configuration with godotenv
- Air (Live reload)

### Frontend
- HTMX for dynamic interactions
- Hyperscript for enhanced interactivity
- Tailwind CSS with DaisyUI components
- Webpack for asset bundling

## 📦 Prerequisites

- Go 1.x
- Node.js and pnpm
- Air (Go live reload tool)
- Environment variables (copy from .env.example)

## 🚀 Getting Started

1. Clone the repository
2. Copy the environment configuration:
   ```bash
   cp .env.example .env
   ```

3. Install frontend dependencies:
   ```bash
   pnpm install
   ```

4. Build the frontend assets:
   ```bash
   pnpm run build
   ```

5. Run the application:
   ```bash
   # Using standard Go
   go run main.go

   # Using Air for live reload during development
   air
   ```

The application will be available at `http://localhost:3000`

## 🐳 Docker Deployment

You can also run the application using Docker:

1. Pull the Docker image:
   ```bash
   docker pull ghcr.io/romaingallez/scform-notes:latest
   ```

2. Create a `.env.docker` file from the example:
   ```bash
   cp .env.docker.example .env.docker
   ```
   
   Make sure to update the environment variables in `.env.docker` as needed.

3. Run browserless/chrome container (required for form processing):
   ```bash
   docker run -d -p 1337:3000 --rm --name chrome browserless/chrome
   ```

4. Run the application container:
   ```bash
   docker run -d -p 3000:3000 --rm --env-file .env.docker --name scform-notes ghcr.io/romaingallez/scform-notes:latest
   ```

5. Access the application at `http://localhost:3000`

### Docker Compose (Alternative)

You can also use Docker Compose to run both containers:

1. Create a `docker-compose.yml` file:
   ```yaml
   services:
     app:
       image: ghcr.io/romaingallez/scform-notes:latest
       ports:
         - "3000:3000"
       env_file:
         - .env.docker
       depends_on:
         - chrome
     chrome:
       image: browserless/chrome
       ports:
         - "1337:3000"
   ```

2. Run with Docker Compose:
   ```bash
   docker-compose up -d
   ```

## 🔧 Development

### Frontend Development
- Watch for Tailwind CSS changes:
  ```bash
  pnpm run watch
  ```

### Backend Development
- The application uses Go modules for dependency management
- Main application entry point is in `main.go`
- For live reload during development, use Air:
  ```bash
  # Install Air if you haven't already
  go install github.com/cosmtrek/air@latest

  # Run the application with Air
  air
  ```
- Core logic is organized in the `internals` directory:
  - `scform/`: Form-related functionality
  - `utils/`: Utility functions and helpers
  - `web/`: Web server and routing logic

## 📁 Project Structure

```
.
├── assets/          # Frontend assets
├── internals/       # Core application logic
│   ├── scform/      # Form processing
│   ├── utils/       # Utility functions
│   └── web/         # Web server and routing
├── views/           # HTML templates
├── main.go         # Application entry point
├── go.mod          # Go dependencies
└── package.json    # Frontend dependencies
```

## 📄 License

ISC License

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

# SCForm Notes Viewer

A web application to view and manage grades from SCForm (Student Course Form) system.

## Features

### Excel-like Table Interface
- **Sortable Columns**: Click on any column header to sort by Matière, Évaluation, Note, Coefficient, Date, or Type
- **Real-time Search**: Filter results by course name, evaluation title, or type
- **Pagination**: Navigate through large datasets with configurable items per page (10, 25, 50, 100)
- **Alternating Row Colors**: Excel-like zebra striping for better readability
- **Responsive Design**: Works on desktop and mobile devices
- **Live Updates**: Real-time progress updates during grade retrieval

### Data Management
- **Import/Export**: Import grades from JSON files or export to JSON/Excel formats
- **Print Support**: Generate print-friendly versions of grade reports
- **WebSocket Integration**: Real-time progress updates during data retrieval

### User Interface
- **Modern Design**: Clean, professional interface using Tailwind CSS
- **Interactive Elements**: Hover effects, loading states, and smooth transitions
- **Accessibility**: Keyboard navigation and screen reader support

## Technology Stack

- **Backend**: Go with Fiber framework
- **Frontend**: Alpine.js for reactivity, HTMX for dynamic updates
- **Styling**: Tailwind CSS with DaisyUI components
- **Real-time**: WebSocket for progress updates
- **Data Processing**: Chrome DevTools Protocol for web scraping

## Installation

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set environment variables (see Configuration section)
4. Run the application: `go run main.go`

## Configuration

Set the following environment variables:
- `SCFORM_URL`: Your SCForm instance URL
- `SCFORM_USERNAME`: Default username (optional)
- `SCFORM_PASSWORD`: Default password (optional)

## Usage

1. Navigate to the application in your browser
2. Enter your SCForm credentials or use the default ones
3. Click "Obtenir les Notes" to retrieve your grades
4. Use the Excel-like table to sort, filter, and navigate your grades
5. Export data or generate print reports as needed

## API Endpoints

- `GET /api/grades`: Returns grades data as JSON for the table interface
- `POST /grades`: Initiates grade retrieval process
- `GET /export`: Download grades as JSON
- `GET /export/excel`: Download grades as Excel file
- `POST /import`: Import grades from JSON file
- `GET /print`: Generate print-friendly version
