# JSONPlaceholder Users Management App

A modern, responsive React application built with Vite and TypeScript for managing user data from the JSONPlaceholder API.

## 🚀 Features

- **User List Display**: Clean table layout with user information
- **Responsive Design**: Mobile-first approach with breakpoints for all devices
- **User Details Modal**: Comprehensive user information with map integration
- **User Management**: Client-side user deletion functionality
- **Modern UI**: Clean design with CSS Modules and CSS custom properties
- **Accessibility**: ARIA labels, keyboard navigation, and semantic HTML
- **TypeScript**: Full type safety throughout the application

## 🛠️ Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **CSS Modules** for component-scoped styling
- **JSONPlaceholder API** for user data
- **Custom Hooks** for state management
- **Modern CSS** with custom properties and animations

## 📁 Project Structure

```
src/
├── components/           # Reusable UI components
│   ├── common/          # Generic components (Button, etc.)
│   └── ui/              # Feature-specific components
├── pages/               # Page components
├── hooks/               # Custom React hooks
├── services/            # API services
├── types/               # TypeScript type definitions
├── styles/              # Global styles and CSS variables
└── utils/               # Utility functions
```

## 🎯 Key Components

### UserTable Component

- Displays users in a responsive table format
- Includes columns for name/email, address, phone, website, and company
- Click functionality to view user details
- Delete functionality with confirmation

### UserModal Component

- Shows detailed user information in a modal
- Includes map link using geo coordinates
- Proper accessibility with focus management
- Smooth animations and responsive design

### Button Component

- Reusable button with multiple variants (primary, secondary, danger)
- Different sizes (sm, md, lg)
- Loading states and accessibility features

## 🚀 Getting Started

1. **Navigate to the project directory:**

   ```bash
   cd jsonplaceholder-api/frontend-vite
   ```

2. **Install dependencies:**

   ```bash
   npm install
   ```

3. **Start the development server:**

   ```bash
   npm run dev
   ```

4. **Open your browser and visit:**
   ```
   http://localhost:5173
   ```

## 📱 Responsive Design

The application is fully responsive with breakpoints at:

- **Mobile**: 480px and below
- **Tablet**: 768px and below
- **Desktop**: 992px and above
- **Large Desktop**: 1200px and above

## 🎨 Design System

The app uses a consistent design system with:

- **CSS Custom Properties** for colors, spacing, and typography
- **Modular CSS** with CSS Modules for component styling
- **Consistent spacing** using a spacing scale
- **Modern typography** with system fonts
- **Smooth animations** for interactions

## 🔧 Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## 📊 Features Breakdown

### Data Management

- Fetches users from JSONPlaceholder API
- Client-side filtering and deletion
- Error handling with user feedback
- Loading states for better UX

### User Experience

- Clean, modern interface
- Intuitive navigation
- Visual feedback for all interactions
- Accessible design following WCAG guidelines

### Performance

- Fast loading with Vite
- Optimized images and assets
- Efficient rendering with React 18
- Code splitting for better performance

## 🌐 API Integration

The app integrates with the JSONPlaceholder API:

- **Base URL**: `https://jsonplaceholder.typicode.com`
- **Endpoints Used**:
  - `GET /users` - Fetch all users
  - `DELETE /users/:id` - Delete user (demo only)

## 🚀 Production Build

To build for production:

```bash
npm run build
```

The built files will be in the `dist` directory and ready for deployment.

---

Built with ❤️ using React, TypeScript, and modern web standards.
