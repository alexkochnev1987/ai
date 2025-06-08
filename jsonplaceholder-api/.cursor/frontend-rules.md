# JSONPlaceholder Frontend App Development Rules

## Project Structure & Architecture

### Directory Structure

```
src/
├── components/           # Reusable UI components
│   ├── common/          # Generic components (Button, Modal, etc.)
│   ├── layout/          # Layout components (Header, Footer, etc.)
│   └── ui/              # UI-specific components
├── pages/               # Page components
├── hooks/               # Custom React hooks
├── services/            # API services and external integrations
├── types/               # TypeScript type definitions
├── utils/               # Utility functions
├── contexts/            # React contexts
├── styles/              # Global styles and themes
│   ├── globals.css      # Global CSS
│   ├── variables.css    # CSS custom properties
│   └── mixins.css       # CSS mixins
├── assets/              # Static assets (images, icons)
├── tests/               # Test utilities and setup
│   ├── __mocks__/       # Mock files
│   ├── fixtures/        # Test data fixtures
│   └── utils/           # Test utilities
└── docs/                # Component documentation
```

### Component Organization

- Each component should have its own directory with:
  - `ComponentName.tsx` - Main component file
  - `ComponentName.module.css` - CSS Module styles
  - `ComponentName.test.tsx` - Component tests
  - `ComponentName.stories.tsx` - Storybook stories (optional)
  - `index.ts` - Export barrel file

## TypeScript Configuration

### Type Safety Rules

- Always use strict TypeScript configuration
- No `any` types allowed - use proper typing or `unknown`
- Use interface for object types, type aliases for unions/primitives
- Define prop types as interfaces with descriptive names
- Use generic types for reusable components
- Implement proper error types for API responses

### API Types

```typescript
// Define types matching the backend API
interface User {
  id: number;
  name: string;
  username: string;
  email: string;
  phone: string;
  website: string;
  address: Address;
  company: Company;
}

interface Address {
  street: string;
  suite: string;
  city: string;
  zipcode: string;
  geo: Geo;
}

interface Geo {
  lat: string;
  lng: string;
}

interface Company {
  name: string;
  catchPhrase: string;
  bs: string;
}

// API Response types
interface APIResponse<T> {
  success: boolean;
  message: string;
  data: T;
  timestamp: string;
}

interface PaginatedResponse<T> extends APIResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}
```

## React Component Rules

### Functional Components

- Use functional components with hooks exclusively
- Implement proper PropTypes with TypeScript interfaces
- Use React.memo for performance optimization when appropriate
- Follow the Single Responsibility Principle

### Component Structure Template

```typescript
import React, { useState, useEffect } from "react";
import styles from "./ComponentName.module.css";

interface ComponentNameProps {
  // Define all props with proper types
  prop1: string;
  prop2?: number; // Optional props
  onAction?: (param: string) => void; // Event handlers
}

export const ComponentName: React.FC<ComponentNameProps> = ({
  prop1,
  prop2 = defaultValue,
  onAction,
}) => {
  // State management
  const [state, setState] = useState<StateType>(initialState);

  // Effects
  useEffect(() => {
    // Effect logic
  }, [dependencies]);

  // Event handlers
  const handleEvent = (param: string) => {
    // Handler logic
    onAction?.(param);
  };

  return <div className={styles.container}>{/* JSX content */}</div>;
};

export default ComponentName;
```

### Custom Hooks Rules

- Extract complex stateful logic into custom hooks
- Use descriptive names starting with "use"
- Return objects for multiple values, not arrays
- Include proper TypeScript return types

## API Integration Rules

### Service Layer

```typescript
// services/api.ts
const API_BASE_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";

class ApiService {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  setAuthToken(token: string) {
    this.token = token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<APIResponse<T>> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers = {
      "Content-Type": "application/json",
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers,
    };

    try {
      const response = await fetch(url, { ...options, headers });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      throw new Error(`API request failed: ${error}`);
    }
  }

  // User-related methods
  async getUsers(page = 1, limit = 10): Promise<PaginatedResponse<User>> {
    return this.request<User[]>(`/api/v1/users?page=${page}&limit=${limit}`);
  }

  async getUser(id: number): Promise<APIResponse<User>> {
    return this.request<User>(`/api/v1/users/${id}`);
  }

  async deleteUser(id: number): Promise<APIResponse<null>> {
    return this.request<null>(`/api/v1/users/${id}`, { method: "DELETE" });
  }
}

export const apiService = new ApiService(API_BASE_URL);
```

### Error Handling

- Implement global error boundary for unhandled errors
- Use try-catch blocks for async operations
- Provide user-friendly error messages
- Log errors for debugging

## User Interface Requirements

### User List Table Component

```typescript
interface UserTableProps {
  users: User[];
  onUserClick: (user: User) => void;
  onUserDelete: (userId: number) => void;
  isLoading?: boolean;
}

export const UserTable: React.FC<UserTableProps> = ({
  users,
  onUserClick,
  onUserDelete,
  isLoading = false,
}) => {
  return (
    <div className={styles.tableContainer}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Name & Email</th>
            <th>Address</th>
            <th>Phone</th>
            <th>Website</th>
            <th>Company</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <UserTableRow
              key={user.id}
              user={user}
              onClick={() => onUserClick(user)}
              onDelete={() => onUserDelete(user.id)}
            />
          ))}
        </tbody>
      </table>
      {isLoading && <LoadingSpinner />}
    </div>
  );
};
```

### User Detail Modal Component

```typescript
interface UserModalProps {
  user: User | null;
  isOpen: boolean;
  onClose: () => void;
}

export const UserModal: React.FC<UserModalProps> = ({
  user,
  isOpen,
  onClose,
}) => {
  if (!isOpen || !user) return null;

  const mapUrl = `https://maps.google.com/?q=${user.address.geo.lat},${user.address.geo.lng}`;

  return (
    <div className={styles.modalOverlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.modalHeader}>
          <h2>{user.name}</h2>
          <button className={styles.closeButton} onClick={onClose}>
            ×
          </button>
        </div>
        <div className={styles.modalContent}>
          {/* User details */}
          <a href={mapUrl} target="_blank" rel="noopener noreferrer">
            View on Map
          </a>
        </div>
      </div>
    </div>
  );
};
```

## CSS Modules & Styling Rules

### CSS Architecture

- Use CSS Modules for component-specific styles
- Follow BEM methodology within CSS Modules
- Use CSS custom properties for theming
- Implement mobile-first responsive design

### CSS Variables Structure

```css
/* styles/variables.css */
:root {
  /* Colors */
  --color-primary: #007bff;
  --color-secondary: #6c757d;
  --color-success: #28a745;
  --color-danger: #dc3545;
  --color-warning: #ffc107;
  --color-info: #17a2b8;

  --color-text: #212529;
  --color-text-secondary: #6c757d;
  --color-background: #ffffff;
  --color-surface: #f8f9fa;
  --color-border: #dee2e6;

  /* Typography */
  --font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
    sans-serif;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;

  /* Spacing */
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 2rem;
  --spacing-2xl: 3rem;

  /* Breakpoints */
  --breakpoint-sm: 576px;
  --breakpoint-md: 768px;
  --breakpoint-lg: 992px;
  --breakpoint-xl: 1200px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);

  /* Border radius */
  --radius-sm: 0.125rem;
  --radius-md: 0.25rem;
  --radius-lg: 0.5rem;

  /* Transitions */
  --transition-fast: 0.15s ease-in-out;
  --transition-normal: 0.3s ease-in-out;
  --transition-slow: 0.5s ease-in-out;
}
```

### Component Styling Template

```css
/* ComponentName.module.css */
.container {
  /* Container styles */
}

.title {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-md);
}

/* Responsive design */
@media (max-width: 768px) {
  .container {
    padding: var(--spacing-sm);
  }
}

/* State modifiers */
.container.loading {
  opacity: 0.6;
  pointer-events: none;
}

.button {
  background: var(--color-primary);
  color: white;
  border: none;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.button:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
```

### Responsive Design Rules

- Mobile-first approach (min-width media queries)
- Flexible grid system using CSS Grid/Flexbox
- Responsive typography with clamp() function
- Touch-friendly interactive elements (min 44px)
- Proper spacing for different screen sizes

## State Management Rules

### Local State

- Use useState for component-level state
- Use useReducer for complex state logic
- Lift state up when multiple components need access

### Global State (if needed)

- Use React Context for global state
- Implement custom hooks for context consumption
- Consider state management libraries for complex apps

### Example Context Implementation

```typescript
// contexts/UserContext.tsx
interface UserContextType {
  users: User[];
  selectedUser: User | null;
  isLoading: boolean;
  error: string | null;
  fetchUsers: () => Promise<void>;
  selectUser: (user: User | null) => void;
  deleteUser: (userId: number) => Promise<void>;
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const useUsers = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error("useUsers must be used within UserProvider");
  }
  return context;
};
```

## Testing Rules

### Testing Framework Setup

- Use Jest + React Testing Library for unit/integration tests
- Use MSW (Mock Service Worker) for API mocking
- Implement Cypress or Playwright for E2E tests
- Use @testing-library/jest-dom for custom matchers

### Test Structure

```typescript
// Component.test.tsx
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { UserTable } from "./UserTable";
import { mockUsers } from "../tests/fixtures/users";

describe("UserTable", () => {
  const defaultProps = {
    users: mockUsers,
    onUserClick: jest.fn(),
    onUserDelete: jest.fn(),
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders user table with correct data", () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByRole("table")).toBeInTheDocument();
    expect(screen.getByText(mockUsers[0].name)).toBeInTheDocument();
  });

  it("calls onUserClick when user row is clicked", () => {
    render(<UserTable {...defaultProps} />);

    fireEvent.click(screen.getByText(mockUsers[0].name));

    expect(defaultProps.onUserClick).toHaveBeenCalledWith(mockUsers[0]);
  });

  it("calls onUserDelete when delete button is clicked", async () => {
    render(<UserTable {...defaultProps} />);

    fireEvent.click(screen.getByLabelText("Delete user"));

    await waitFor(() => {
      expect(defaultProps.onUserDelete).toHaveBeenCalledWith(mockUsers[0].id);
    });
  });
});
```

### Testing Rules

- Test behavior, not implementation
- Use data-testid sparingly, prefer semantic queries
- Mock external dependencies
- Test error states and loading states
- Aim for 80%+ code coverage
- Write integration tests for critical user flows

### Test Utilities

```typescript
// tests/utils/renderWithProviders.tsx
export const renderWithProviders = (
  ui: ReactElement,
  { preloadedState = {}, ...renderOptions } = {}
) => {
  const Wrapper = ({ children }: { children: ReactNode }) => (
    <UserProvider initialState={preloadedState}>{children}</UserProvider>
  );

  return render(ui, { wrapper: Wrapper, ...renderOptions });
};
```

### Mock Data

```typescript
// tests/fixtures/users.ts
export const mockUsers: User[] = [
  {
    id: 1,
    name: "John Doe",
    username: "johndoe",
    email: "john@example.com",
    phone: "123-456-7890",
    website: "john.com",
    address: {
      street: "123 Main St",
      suite: "Apt 1",
      city: "Anytown",
      zipcode: "12345",
      geo: { lat: "40.7128", lng: "-74.0060" },
    },
    company: {
      name: "Acme Corp",
      catchPhrase: "Innovation at its best",
      bs: "synergistic solutions",
    },
  },
];
```

## Performance Optimization Rules

### React Performance

- Use React.memo for expensive components
- Implement useCallback for stable function references
- Use useMemo for expensive calculations
- Lazy load components with React.lazy
- Implement virtual scrolling for large lists

### Bundle Optimization

- Use code splitting for route-based chunks
- Optimize images with appropriate formats
- Implement service worker for caching
- Use tree shaking for unused code elimination

### User Experience

- Implement loading states for all async operations
- Add skeleton screens for better perceived performance
- Use optimistic updates where appropriate
- Implement proper error boundaries

## Accessibility Rules

### ARIA Implementation

- Use semantic HTML elements
- Implement proper ARIA labels and roles
- Ensure keyboard navigation works
- Provide focus management for modals
- Include skip links for navigation

### Testing Accessibility

- Use @testing-library/jest-dom accessibility matchers
- Test with screen readers
- Ensure color contrast ratios meet WCAG guidelines
- Test keyboard-only navigation

## Documentation Rules

### Component Documentation

````typescript
/**
 * UserTable component displays a list of users in a table format
 *
 * @example
 * ```tsx
 * <UserTable
 *   users={users}
 *   onUserClick={(user) => setSelectedUser(user)}
 *   onUserDelete={(id) => deleteUser(id)}
 *   isLoading={loading}
 * />
 * ```
 */
export interface UserTableProps {
  /** Array of user objects to display */
  users: User[];
  /** Callback fired when a user row is clicked */
  onUserClick: (user: User) => void;
  /** Callback fired when delete button is clicked */
  onUserDelete: (userId: number) => void;
  /** Whether the table is in loading state */
  isLoading?: boolean;
}
````

### README Structure

```markdown
# JSONPlaceholder Frontend App

## Overview

Brief description of the application and its purpose.

## Features

- User list display with table layout
- User detail modal with comprehensive information
- Responsive design for all screen sizes
- Modern UI with animations and interactions

## Technology Stack

- React 18 with TypeScript
- CSS Modules for styling
- React Testing Library for testing
- MSW for API mocking

## Getting Started

Installation and running instructions

## API Integration

Documentation of API endpoints and data flow

## Component Documentation

Links to component documentation

## Testing

How to run tests and testing strategy

## Deployment

Deployment instructions
```

## Build & Development Rules

### Development Environment

- Use Create React App or Vite for build tooling
- Configure ESLint and Prettier for code quality
- Use Husky for git hooks
- Implement absolute imports for cleaner import statements

### Environment Configuration

```typescript
// config/env.ts
export const config = {
  apiUrl: process.env.REACT_APP_API_URL || "http://localhost:8080",
  environment: process.env.NODE_ENV || "development",
  enableMocking: process.env.REACT_APP_ENABLE_MOCKING === "true",
} as const;
```

### Error Monitoring

- Implement error boundary for production error handling
- Use error logging service (Sentry, LogRocket)
- Provide fallback UI for error states

## Testing Configuration & Rules

### Jest Configuration

```javascript
// jest.config.js
module.exports = {
  testEnvironment: "jsdom",
  setupFilesAfterEnv: ["<rootDir>/src/tests/setup.ts"],
  moduleNameMapping: {
    "\\.(css|less|scss|sass)$": "identity-obj-proxy",
    "^@/(.*)$": "<rootDir>/src/$1",
  },
  collectCoverageFrom: [
    "src/**/*.{ts,tsx}",
    "!src/**/*.d.ts",
    "!src/tests/**/*",
    "!src/stories/**/*",
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
};
```

### Test Categories

1. **Unit Tests**: Individual component and function testing
2. **Integration Tests**: Component interaction testing
3. **API Tests**: Service layer and API integration testing
4. **E2E Tests**: Full user journey testing
5. **Accessibility Tests**: Screen reader and keyboard navigation testing

### Test Naming Conventions

- Test files: `*.test.tsx` for components, `*.test.ts` for utilities
- Test descriptions: Use descriptive names explaining the behavior
- Test organization: Group related tests using `describe` blocks
- Mock files: Place in `__mocks__` directory with matching structure

## Documentation Configuration & Rules

### Documentation Structure

```
docs/
├── components/          # Component documentation
│   ├── UserTable.md    # Individual component docs
│   └── UserModal.md    # Modal component docs
├── api/                # API documentation
│   ├── endpoints.md    # API endpoint documentation
│   └── types.md        # Type definitions documentation
├── guides/             # Development guides
│   ├── setup.md        # Project setup guide
│   ├── testing.md      # Testing guide
│   └── deployment.md   # Deployment guide
└── README.md           # Main project documentation
```

### Component Documentation Template

````markdown
# ComponentName

## Overview

Brief description of the component's purpose and functionality.

## Props

| Prop  | Type   | Required | Default | Description          |
| ----- | ------ | -------- | ------- | -------------------- |
| prop1 | string | Yes      | -       | Description of prop1 |
| prop2 | number | No       | 0       | Description of prop2 |

## Usage Examples

```tsx
// Basic usage
<ComponentName prop1="value" />

// Advanced usage
<ComponentName
  prop1="value"
  prop2={42}
  onAction={(data) => handleAction(data)}
/>
```
````

```

## Styling

Description of available CSS classes and customization options.

## Accessibility

Description of accessibility features and ARIA attributes.

## Testing

Examples of how to test this component.
```

### Auto-Documentation Rules

- Use TypeScript interfaces for automatic prop documentation
- Implement JSDoc comments for all public APIs
- Generate documentation from code using tools like Typedoc
- Keep examples up-to-date with actual implementation
- Include visual examples using Storybook or similar tools

## Modern React Patterns (2024+)

### React 18+ Features

```typescript
// Concurrent Features
import { startTransition, useDeferredValue } from "react";

const SearchComponent = () => {
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query);

  const handleSearch = (value: string) => {
    startTransition(() => {
      setQuery(value);
    });
  };

  return <SearchResults query={deferredQuery} />;
};

// Suspense for Data Fetching
const UserProfile = ({ userId }: { userId: number }) => {
  return (
    <Suspense fallback={<ProfileSkeleton />}>
      <UserData userId={userId} />
    </Suspense>
  );
};
```

### Server Components Integration

- Use Next.js 13+ App Router for server components when applicable
- Implement streaming for better perceived performance
- Use server actions for form submissions

## Enhanced Security Rules

### Content Security Policy

```typescript
// Next.js security headers
const securityHeaders = [
  {
    key: "Content-Security-Policy",
    value: `
      default-src 'self';
      script-src 'self' 'unsafe-eval' 'unsafe-inline';
      style-src 'self' 'unsafe-inline';
      img-src 'self' data: https:;
      connect-src 'self' ${process.env.NEXT_PUBLIC_API_URL};
    `
      .replace(/\s{2,}/g, " ")
      .trim(),
  },
];
```

### XSS Prevention

```typescript
// Sanitize user input
import DOMPurify from "dompurify";

const SafeHTML = ({ content }: { content: string }) => {
  const sanitizedContent = DOMPurify.sanitize(content);
  return <div dangerouslySetInnerHTML={{ __html: sanitizedContent }} />;
};
```

## Performance Monitoring

### Web Vitals Tracking

```typescript
// utils/webVitals.ts
export const reportWebVitals = (metric: any) => {
  // Send to analytics service
  if (process.env.NODE_ENV === "production") {
    // Google Analytics, DataDog, etc.
    console.log(metric);
  }
};
```

These comprehensive rules ensure a maintainable, scalable, and professional frontend application that integrates seamlessly with the JSONPlaceholder API backend while maintaining high code quality, comprehensive testing, and thorough documentation.

```

```
