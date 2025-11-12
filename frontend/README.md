# Blog Platform Frontend

A modern, responsive frontend for the Blog Platform API built with React, TypeScript, and Tailwind CSS.

## Features

### Core Functionality
- **User Authentication**: Login, register, and OAuth (Google/GitHub)
- **Blog Management**: Create, read, update, and delete blog posts
- **AI Integration**: Generate content, enhance writing, and get topic suggestions
- **Search & Filtering**: Advanced search by title, author, tags, and date
- **User Profiles**: Manage profile information and profile pictures
- **Responsive Design**: Mobile-first design with Tailwind CSS

### Technical Features
- **TypeScript**: Full type safety throughout the application
- **React Query**: Efficient data fetching and caching
- **React Router**: Client-side routing with protected routes
- **React Hook Form**: Form handling with validation
- **Axios**: HTTP client with interceptors for authentication
- **Toast Notifications**: User feedback with react-hot-toast
- **Modern UI**: Clean, professional interface with Lucide icons

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend API running on `http://localhost:8080`

### Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Open your browser and navigate to `http://localhost:5173`

### Build for Production

```bash
npm run build
```

## Project Structure

```
frontend/
├── src/
│   ├── components/          # Reusable UI components
│   │   ├── Header.tsx
│   │   ├── BlogCard.tsx
│   │   ├── LoadingSpinner.tsx
│   │   └── Pagination.tsx
│   ├── contexts/           # React contexts
│   │   └── AuthContext.tsx
│   ├── lib/                # Utility libraries
│   │   ├── api.ts          # Axios configuration
│   │   ├── auth.ts         # Authentication API
│   │   ├── blog.ts         # Blog API
│   │   ├── ai.ts           # AI API
│   │   └── utils.ts        # Utility functions
│   ├── pages/              # Page components
│   │   ├── Home.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Blogs.tsx
│   │   ├── BlogDetail.tsx
│   │   ├── CreateBlog.tsx
│   │   ├── Profile.tsx
│   │   └── AITools.tsx
│   ├── types/              # TypeScript type definitions
│   │   └── index.ts
│   ├── App.tsx             # Main application component
│   ├── main.tsx            # Application entry point
│   └── index.css           # Global styles
├── public/                 # Static assets
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── vite.config.ts
└── README.md
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## API Integration

The frontend is configured to work with the backend API at `http://localhost:8080`. The API proxy is set up in Vite configuration to handle CORS issues during development.

### Authentication

- JWT tokens are stored in localStorage
- Automatic token refresh on 401 responses
- Protected routes require authentication
- OAuth integration for Google and GitHub

### Features

#### Blog Management
- Create blogs with rich content and tags
- Edit and delete own blogs (admin can edit/delete any)
- View blogs with engagement metrics
- Like/dislike functionality
- Comment system

#### AI Tools
- Generate blog content from topics
- Enhance existing content
- Get topic suggestions based on keywords

#### Search & Filtering
- Search by blog title
- Filter by author username
- Filter by tags
- Sort by newest, oldest, or popularity

## Configuration

### Environment Variables

Create a `.env` file in the frontend root:

```env
VITE_API_BASE_URL=http://localhost:8080
```

### Customization

- **Colors**: Modify `tailwind.config.js` to customize the color scheme
- **Typography**: Font settings in `index.css`
- **API Endpoints**: Update API calls in `lib/` directory

## Contributing

1. Follow the existing code style and patterns
2. Use TypeScript for all new code
3. Ensure all components are properly typed
4. Test responsive design on different screen sizes
5. Run linting before committing: `npm run lint`

## Technologies Used

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **React Router** - Client-side routing
- **React Query** - Data fetching and state management
- **React Hook Form** - Form handling
- **Axios** - HTTP client
- **Lucide React** - Icon library
- **React Hot Toast** - Notification system

## License

This project is licensed under the MIT License.
