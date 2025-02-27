# SCForm Notes

A web application built with Go (Fiber) and modern frontend technologies for managing and processing form data. The application uses HTMX for dynamic interactions and Tailwind CSS for styling.

## 🚀 Features

- Modern web interface with Tailwind CSS and DaisyUI
- Server-side rendering with Go templates
- Dynamic UI updates using HTMX
- Form data processing and management
- Asset management system
- Environment configuration support

## 🛠 Tech Stack

### Backend
- Go (Fiber web framework)
- HTML Templates
- Environment configuration with godotenv

### Frontend
- HTMX for dynamic interactions
- Hyperscript for enhanced interactivity
- Tailwind CSS with DaisyUI components
- Webpack for asset bundling

## 📦 Prerequisites

- Go 1.x
- Node.js and pnpm
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
   go run main.go
   ```

The application will be available at `http://localhost:3000`

## 🔧 Development

### Frontend Development
- Watch for Tailwind CSS changes:
  ```bash
  pnpm run watch
  ```

### Backend Development
- The application uses Go modules for dependency management
- Main application entry point is in `main.go`
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
