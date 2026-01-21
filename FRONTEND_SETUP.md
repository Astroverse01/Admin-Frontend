# Frontend Setup Guide

## Quick Start

1. **Navigate to frontend directory:**
   ```bash
   cd frontend
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Start development server:**
   ```bash
   npm run dev
   ```

4. **Access the dashboard:**
   - Open your browser and go to `http://localhost:3000`
   - Login with your admin credentials

## Features Implemented

✅ **Authentication**
- Login page with JWT token management
- Protected routes with automatic token refresh
- Logout functionality

✅ **Dashboard**
- Overview statistics for users, astrologers, and complaints
- Quick action cards
- Real-time data fetching

✅ **User Management**
- List all users with pagination
- Search and filter users
- Activate/Deactivate users

✅ **Astrologer Management**
- List all astrologers with pagination
- Search and filter astrologers
- Update status (active/inactive)
- Toggle visibility (visible/hidden)

✅ **Complaint Management**
- **User Service Complaints**: View, accept/reject with refund details
- **User General Complaints**: View and close complaints
- **Astro General Complaints**: View and close complaints
- All with pagination and filtering

✅ **Horoscope Management**
- List all horoscopes
- Create individual horoscopes
- Bulk create horoscopes
- Update horoscopes
- Delete horoscopes

✅ **Scheduler**
- Manual trigger for daily report generation
- Visual feedback and status updates
- Information about collections processed

## API Integration

All APIs from `main.go` are integrated:

- ✅ `POST /admin/login` - Authentication
- ✅ `GET /admin/users` - List users
- ✅ `PATCH /admin/users/:userId/deactivate` - Deactivate user
- ✅ `GET /admin/astros` - List astrologers
- ✅ `PATCH /admin/astros/:astroId/status` - Update astro status
- ✅ `PATCH /admin/astros/:astroId/visibility` - Toggle visibility
- ✅ `GET /admin/user-service-complaints` - List service complaints
- ✅ `GET /admin/user-service-complaints/:serviceType/:orderId` - Get complaint details
- ✅ `PATCH /admin/user-service-complaints/:orderId` - Accept/reject complaint
- ✅ `GET /admin/user-general-complaints` - List user general complaints
- ✅ `PATCH /admin/user-general-complaints/:problemId/close` - Close complaint
- ✅ `GET /admin/astro-general-complaints` - List astro general complaints
- ✅ `PATCH /admin/astro-general-complaints/:problemId/close` - Close complaint
- ✅ `GET /admin/horoscopes` - List horoscopes
- ✅ `POST /admin/horoscopes/bulk` - Bulk create horoscopes
- ✅ `PATCH /admin/horoscopes/:horoscopeId` - Update horoscope
- ✅ `DELETE /admin/horoscopes/:horoscopeId` - Delete horoscope
- ✅ `POST /admin/scheduler/trigger` - **Manual scheduler trigger** ✨

## Configuration

The API base URL can be configured via environment variable:

```env
VITE_API_URL=https://admin-panel-fe.astrosway.com
```

If not set, it defaults to `https://admin-panel-fe.astrosway.com`.

## Build for Production

```bash
npm run build
```

The production build will be in the `dist` directory.

## Technologies

- **React 18** - UI framework
- **Vite 5** - Build tool and dev server
- **React Router 6** - Routing
- **Axios** - HTTP client
- **Tailwind CSS** - Styling
- **Lucide React** - Icons
- **date-fns** - Date formatting

## Notes

- Make sure the backend server is running on port 8082 (or update `VITE_API_URL`)
- The scheduler trigger API is fully functional and can be manually triggered from the dashboard
- All API calls include JWT authentication automatically
- The UI is fully responsive and works on mobile devices

