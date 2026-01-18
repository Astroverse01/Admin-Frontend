# Astro Admin Frontend

A modern React + Vite dashboard for managing the Astro Admin backend.

## Features

- 🚀 Fast development with Vite
- 🎨 Modern UI with Tailwind CSS
- 🔐 JWT Authentication
- 📊 Dashboard with statistics
- 👥 User Management
- 🔮 Astrologer Management
- 📝 Complaint Management (User Service, User General, Astro General)
- 🌟 Horoscope Management
- ⏰ Scheduler Trigger

## Getting Started

### Prerequisites

- Node.js 16+ and npm/yarn

### Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file (optional, defaults to `http://localhost:8082`):
```env
VITE_API_URL=http://localhost:8082
```

4. Start the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

### Preview Production Build

```bash
npm run preview
```

## Project Structure

```
frontend/
├── src/
│   ├── components/          # React components
│   │   ├── Dashboard.jsx
│   │   ├── Users.jsx
│   │   ├── Astros.jsx
│   │   ├── UserServiceComplaints.jsx
│   │   ├── UserGeneralComplaints.jsx
│   │   ├── AstroGeneralComplaints.jsx
│   │   ├── Horoscopes.jsx
│   │   ├── Scheduler.jsx
│   │   ├── Login.jsx
│   │   └── Layout.jsx
│   ├── context/             # React contexts
│   │   └── AuthContext.jsx
│   ├── services/            # API services
│   │   └── api.js
│   ├── App.jsx             # Main app component
│   ├── main.jsx            # Entry point
│   └── index.css           # Global styles
├── index.html
├── package.json
├── vite.config.js
└── tailwind.config.js
```

## API Integration

All API calls are handled through the `services/api.js` file. The API base URL can be configured via the `VITE_API_URL` environment variable.

### Authentication

The app uses JWT tokens stored in localStorage. The token is automatically included in all API requests via axios interceptors.

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Technologies Used

- React 18
- Vite 5
- React Router 6
- Axios
- Tailwind CSS
- Lucide React (Icons)
- date-fns

## Notes

- Make sure the backend server is running on port 8082 (or update the API URL)
- The scheduler trigger API is available at `/admin/scheduler/trigger`
- All protected routes require authentication

