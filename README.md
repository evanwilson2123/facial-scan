# **Facial Scan App (backend)**

## Overview

The Facial Scan App is a web-based application that allows users to upload their photos and receive a detailed analysis of their facial features based on predefined, objective, and scientifically-based measurements. The analysis includes criteria such as symmetry, facial definition, jawline, cheekbones, and more. The application leverages OpenAI's API for generating detailed facial feature assessments and Google Cloud's Firestore and Storage for storing user data and images.

## Features

- User registration and authentication with Firebase.
- Secure file upload and storage on Google Cloud Storage.
- Facial feature analysis using OpenAI's GPT model.
- Leaderboard for top-rated users by gender.
- Comprehensive user profile management.
- Health check endpoint for monitoring application status.

## Tech Stack

- **Backend:** Node.js, Express.js, Fiber (Go)
- **Database:** Google Firestore
- **Storage:** Google Cloud Storage
- **Authentication:** Firebase Auth
- **AI:** OpenAI API
- **Other:** Docker, CORS, Environment Variables

## Prerequisites

- Node.js and npm
- Go (for Fiber-based controllers)
- Docker (optional for containerization)
- Google Cloud account with Firestore and Storage enabled
- Firebase project setup
- OpenAI API key

## Installation

1. **Clone the repository:**
     ```sh
     git clone <repository-url>
     cd facial-scan
2. **Install Node.js dependencies:**
   ```sh
    npm install
3. **Install Go dependencies:**
    ```sh
    go get ./...
    Set up environment variables:
4. **Create a .env file in the root directory and add the following variables:**

env
Copy code
PORT=4000
GPT_KEY=your_openai_api_key
BUCKET_NAME=your_google_cloud_storage_bucket_name
APP_ID=your_firestore_project_id
IMAGE_URL=default_image_url_for_testing
Initialize Firebase:

Place your Firebase service account key file in the config directory and ensure it is named face-scan-71bdf-firebase-adminsdk-5n1l4-e49277c639.json.

Usage

5. **Start the node server:**

    ```sh
    cd node
    npm start
6. **Start the go server:**
   ```sh
    cd server
    go run main.go

***Endpoints:***

**User Registration:**
POST /api/register

Request: JSON object with user details.

**File Upload and Analysis:**
POST /api/upload

Request: FormData with the image file.

**Health Check:**
POST /api/health

Response: JSON object with health status.

**Leaderboard:**
GET /api/leaderboard-male
GET /api/leaderboard-female


**User Account Creation:**
PATCH /api/create-account

**Check if User has Username:**
GET /api/has-username

**Dashboard Load:**
GET /api/dashboard


**Get User's Score:**
GET /api/get-score






