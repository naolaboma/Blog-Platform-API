# Blog Platform Frontend - Setup Instructions

## Quick Start

1. **Install Dependencies:**
   ```bash
   cd frontend
   npm install
   ```

2. **Start Development Server:**
   ```bash
   npm run dev
   ```

3. **Access Application:**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080 (must be running)

## Prerequisites

- Node.js 18+ and npm
- Backend API server running on port 8080
- MongoDB and Redis services for the backend

## Features Implemented

### ✅ Core Features
- **User Authentication**: Login, register, OAuth (Google/GitHub)
- **Blog Management**: Full CRUD operations for blog posts
- **AI Integration**: Content generation, enhancement, and topic suggestions
- **Search & Filtering**: Advanced search with multiple filters
- **User Profiles**: Profile management with picture uploads
- **Admin Panel**: User role management for administrators

### ✅ Technical Features
- **TypeScript**: Complete type safety
- **React Query**: Efficient data fetching and caching
- **Responsive Design**: Mobile-first with Tailwind CSS
- **Protected Routes**: Authentication-based access control
- **Form Validation**: Client-side validation with error handling
- **Toast Notifications**: User feedback system

## Project Structure

```
frontend/
├── src/
│   ├── components/          # Reusable UI components
│   ├── contexts/           # React contexts (Auth)
│   ├── lib/                # API clients and utilities
│   ├── pages/              # Page components
│   ├── types/              # TypeScript definitions
│   └── App.tsx             # Main application
├── package.json
├── vite.config.ts
└── tailwind.config.js
```

## Available Pages

- `/` - Home page with recent/popular blogs
- `/login` - User login
- `/register` - User registration
- `/blogs` - Blog listing with search/filter
- `/blogs/:id` - Individual blog post
- `/create` - Create new blog (protected)
- `/profile` - User profile management (protected)
- `/ai-tools` - AI writing tools (protected)
- `/admin` - Admin panel (admin only)

## Configuration

The frontend automatically connects to the backend at `http://localhost:8080` via Vite proxy configuration.

## Testing the Integration

1. Start your backend server (`go run cmd/server/main.go`)
2. Start the frontend (`npm run dev`)
3. Register a new user or use OAuth
4. Create blogs, test AI features, and explore admin functionality

## Production Build

```bash
npm run build
```

The build output will be in the `dist/` directory, ready for deployment.
