# Mathematical Formulas API

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
