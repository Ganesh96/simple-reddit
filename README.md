# Simple Reddit Clone

This project is a full-stack web application that mimics some of the basic functionalities of Reddit. It features a modern backend API and responsive frontend interface.

## Key Features

- User authentication system with secure JWT implementation
- Complete post management (create, read, update, delete)
- Community/subreddit creation and moderation
- User profile system
- RESTful API design following best practices

## Technical Stack

**Backend:**
- Built with Go using the Gin framework
- MongoDB for data storage
- JWT for authentication

**Frontend:**
- Angular framework with TypeScript
- Angular Material for UI components
- Responsive design

## Getting Started

### Prerequisites

- Go (version 1.15+)
- Node.js and npm (version 14+)
- MongoDB database

### Installation

1. Configure environment variables for both backend and frontend
2. Install required dependencies
3. Start both backend server and frontend development server

The application runs with backend serving APIs on port 8080 and frontend running on port 4200.

## API Overview

The application provides RESTful endpoints for:

### Authentication
- User registration
- Login with JWT token generation

### Posts
- Full CRUD operations for posts
- Pagination and filtering support

### Communities
- Create and manage communities
- Community-specific post feeds

### Users
- Profile management
- Account deletion

## Project Highlights

- Implemented secure authentication flow
- Designed scalable API architecture
- Created responsive UI components
- Optimized database queries for performance
- Followed REST API best practices

Note: Implementation details and source code structure are intentionally omitted as this is a portfolio showcase.