# Formulas API

A powerful API for storing, sharing, and collaborating on mathematical formulas with features similar to a code repository system.

## Overview

This API allows users to store mathematical formulas, fork repositories from other users, star repositories they find useful, and manage public/private visibility settings. Built with Go and utilizing MongoDB and Redis for efficient data management.

## Tech Stack

- **Go Version**: 1.16
- **Dependencies**:
  - go-redis v8.11.4
  - golang-jwt v3.2.2+incompatible
  - uuid v1.3.0
  - gorilla/mux v1.8.0
  - mailgun/mailgun-go/v4 v4.6.0
  - cors v1.8.2
  - mongo-driver v1.8.4

## Features

### Authentication & Security
- Email verification system via Mailgun
- JWT-based authentication with Access and Refresh tokens
- Redis-backed session management for improved performance

### Core Functionality
- Complete CRUD operations for formula management
- Advanced search capabilities
- Public and private repository options
- User follow system
- Fork other users' public repositories
- Star repositories to show appreciation
- Block other spam users

### Data Management
- MongoDB for persistent storage
- Redis for caching and session management

##API Mapping

This document provides a comprehensive overview of all available routes in the application.

## Authentication Routes

### Login Routes

#### POST /login
- **Description**: Authenticates a user and provides access and refresh tokens
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response**: Returns access token, refresh token, username, and email
- **Error Codes**:
  - "010": Invalid request body
  - "030": Invalid credentials
  - "006": Token creation failed
  - "300": Authentication creation failed

#### POST /refresh
- **Description**: Refreshes the access token using a refresh token
- **Headers**: 
  - `Authorization`: Refresh token
- **Response**: Returns new access and refresh tokens
- **Error Codes**:
  - "002": Invalid token
  - "004": Invalid claims
  - "003": Token deletion failed
  - "030": User not found

#### POST /logout
- **Description**: Logs out the user by invalidating their tokens
- **Headers**:
  - `Authorization`: Access token
- **Response**: Success message
- **Error Codes**:
  - "004": Invalid token metadata
  - "400": Token deletion failed

### Signup Routes

#### POST /signup
- **Description**: Creates a new user account
- **Request Body**:
  ```json
  {
    "email": "string",
    "password": "string",
    "username": "string"
  }
  ```
- **Response**: Success message with verification code sent to email
- **Error Codes**:
  - "010": Invalid request body
  - "020": Username/Email already exists
  - "100": Database error
  - "200": Redis error

#### POST /verify-mail
- **Description**: Verifies user's email using verification code
- **Request Body**:
  ```json
  {
    "email": "string",
    "code": "string"
  }
  ```
- **Response**: Success message
- **Error Codes**:
  - "010": Invalid request body
  - "1200": Redis error
  - "005": Invalid verification code
  - "100": Database error

#### POST /resend-verification
- **Description**: Resends verification code to user's email
- **Headers**:
  - `user`: User ID
- **Request Body**:
  ```json
  {
    "email": "string"
  }
  ```
- **Response**: Success message
- **Error Codes**:
  - "010": Invalid request body
  - "030": User not found
  - "500": Redis error
  - "003": Unauthorized

## Helper Routes

#### GET /check-username/{name}
- **Description**: Checks if a username is available
- **Parameters**:
  - `name`: Username to check
- **Response**: Availability status
- **Error Codes**:
  - "100": Database error
  - "020": Username taken
  - "1020": Username available

#### GET /check-email/{mail}
- **Description**: Checks if an email is available
- **Parameters**:
  - `mail`: Email to check
- **Response**: Availability status
- **Error Codes**:
  - "100": Database error
  - "020": Email taken
  - "1020": Email available

#### GET /check-url/{url}
- **Description**: Checks if a URL is available
- **Parameters**:
  - `url`: URL to check
- **Response**: Availability status
- **Error Codes**:
  - "100": Database error
  - "020": URL taken
  - "1020": URL available

## Create Routes

#### POST /create-data
- **Description**: Creates a new formula/data entry
- **Headers**:
  - `user`: User ID (optional)
- **Request Body**:
  ```json
  {
    "name": "string",
    "Structure": "string",
    "user": "string",
    "private": boolean,
    "url": "string"
  }
  ```
- **Response**: Success message
- **Error Codes**:
  - "010": Invalid request body
  - "020": URL already exists
  - "100": Database error
  - "500": Redis error
  - "003": Unauthorized

## Read Routes

#### GET /data/{url}
- **Description**: Retrieves formula/data by URL
- **Headers**:
  - `user`: User ID (required for private data)
- **Parameters**:
  - `url`: URL of the formula/data
- **Response**: Formula data object
- **Error Codes**:
  - "040": Data not found
  - "003": Unauthorized access
  - "100": Database error

#### GET /user/{user}
- **Description**: Retrieves user profile and associated data
- **Headers**:
  - `user`: Logged-in user ID (optional)
- **Parameters**:
  - `user`: Username to fetch
- **Response**: User profile with created, starred, forked data and follow information
- **Error Codes**:
  - "010": User not found
  - "100": Database error

#### POST /search-data
- **Description**: Searches for formulas/data by name
- **Request Body**:
  ```json
  {
    "search": "string"
  }
  ```
- **Response**: Array of matching formulas
- **Error Codes**:
  - "010": Invalid request body
  - "100": Database error

#### POST /search-user
- **Description**: Searches for users by username
- **Request Body**:
  ```json
  {
    "searchUser": "string"
  }
  ```
- **Response**: Array of matching users
- **Error Codes**:
  - "010": Invalid request body
  - "100": Database error

## Update Routes

#### PUT /edit-data/{url}
- **Description**: Updates an existing formula/data
- **Headers**:
  - `user`: User ID
- **Parameters**:
  - `url`: URL of the formula to edit
- **Request Body**:
  ```json
  {
    "name": "string",
    "Structure": "string",
    "private": boolean,
    "url": "string"
  }
  ```
- **Response**: Success message
- **Error Codes**:
  - "008": Not logged in
  - "010": User not found
  - "020": URL already exists
  - "003": Unauthorized
  - "100": Database error

#### PUT /edit-user
- **Description**: Updates user email
- **Headers**:
  - `user`: User ID
- **Request Body**:
  ```json
  {
    "mail": "string"
  }
  ```
- **Response**: Success message
- **Error Codes**:
  - "008": Not logged in
  - "010": Invalid request body
  - "100": Database error

#### PUT /follow-user/{usertofollow}/{set}
- **Description**: Follows or unfollows a user
- **Headers**:
  - `user`: User ID
- **Parameters**:
  - `usertofollow`: Username to follow/unfollow
  - `set`: "1" to follow, "0" to unfollow
- **Response**: Update result
- **Error Codes**:
  - "008": Not logged in
  - "100": Database error

#### PUT /fork-data/{url}
- **Description**: Creates a fork of existing formula/data
- **Headers**:
  - `user`: User ID
- **Parameters**:
  - `url`: URL of the formula to fork
- **Response**: Success message
- **Error Codes**:
  - "008": Not logged in
  - "003": Unauthorized
  - "100": Database error

#### PUT /report-data/{url}
- **Description**: Reports a formula/data
- **Headers**:
  - `user`: User ID
- **Parameters**:
  - `url`: URL of the formula to report
- **Response**: Update result
- **Error Codes**:
  - "008": Not logged in
  - "003": Unauthorized
  - "100": Database error

#### PUT /star-data/{url}/{set}
- **Description**: Stars or unstars a formula/data
- **Headers**:
  - `user`: User ID
- **Parameters**:
  - `url`: URL of the formula
  - `set`: "1" to star, "0" to unstar
- **Response**: Update result
- **Error Codes**:
  - "008": Not logged in
  - "003": Unauthorized
  - "100": Database error

## Delete Routes

#### DELETE /delete-data/{url}
- **Description**: Deletes a formula/data
- **Headers**:
  - `user`: User ID
- **Parameters**:
  - `url`: URL of the formula to delete
- **Response**: Delete result
- **Error Codes**:
  - "008": Not logged in
  - "010": User not found
  - "003": Unauthorized
  - "100": Database error
  - "1020": Data not found





## Error Codes

### Authentication Errors
| Description | Code |
|-------------|------|
| Already Logged In | 001 |
| Logged Out | 1001 |
| Session Expired (Redirect to Login) | 002 |
| Unauthorized | 003 |
| Token Invalid/Not Passed | 004 |
| Verification Code Rejected | 005 |
| Verification Code Accepted | 1005 |
| Token Creation Error | 006 |
| Signup (Redirect to Login) | 007 |
| Redirect to Login | 008 |

### Data Errors
| Description | Code |
|-------------|------|
| Bad Data Sent | 010 |
| Data Already Present | 020 |
| Data Not Present | 1020 |
| Email or Password Incorrect | 030 |
| Data Requested Broken | 040 |
| Data Creation Error | 050 |

### Database Errors
| Description | Code |
|-------------|------|
| DB Operation Error | 100 |
| Successfully Created | 1100 |
| Cache Verification Set Fail | 200 |
| Verification Code Expired | 1200 |
| Redis Auth Set Error | 300 |
| Redis Delete Error | 400 |
| Redis Set Error | 500 |
