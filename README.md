# LoanGuard

LoanGuard is a loan management system designed to streamline the process of applying for loans, managing loan applications, and handling administrative tasks such as loan approval, rejection, and system logging. Built using Go with the Gin framework, LoanGuard implements clean architecture principles, ensuring a scalable and maintainable codebase.

## Features

### User Functionalities
- **Apply for Loan**: Users can submit loan applications with details like amount, interest rate, and loan purpose.
- **View Loan Status**: Users can check the status of their specific loan applications.

### Admin Functionalities
- **View All Loans**: Admins can view all loan applications with filtering options based on status (`pending`, `approved`, `rejected`) and ordering (`asc`, `desc`).
- **Approve/Reject Loan**: Admins can approve or reject loan applications.
- **Delete Loan**: Admins can delete specific loan applications.
- **View System Logs**: Admins can retrieve system logs to track actions like login attempts, loan submissions, loan status updates, and password reset activities.

## Project Structure

LoanGuard follows a clean architecture approach, dividing the project into separate layers:

- **Models**: Defines the core data structures (e.g., `Loan`, `SystemLog`).
- **Repositories**: Handles database operations for various entities.
- **Use Cases**: Contains the business logic for handling requests and processing data.
- **Controllers**: Manages the HTTP request handling, invoking the appropriate use cases.
- **Routes**: Defines the API endpoints and associates them with controllers.

## LoanGuard Setup Guide

This guide provides instructions on setting up and configuring the LoanGuard project. Follow the steps below to prepare your environment.

## 1. Environment Variables

LoanGuard requires several environment variables for configuration. You should define these in a `.env` file located at the root of the project.

### Sample `.env` File

#### Database Configuration
DB_NAME=YourDatabaseName
MONGO_URI=mongodb://localhost:27017

#### Security Keys
ACCESS_SECRET_KEY=your_access_secret_key
REFRESH_SECRET_KEY=your_refresh_secret_key
VERIFICATION_SECRET_KEY=your_verification_secret_key

#### SMTP Configuration (for email services)
USERNAME=your_email_address@gmail.com
SMTP_HOST=smtp.your_email_provider.com
SMTP_PORT=587
PASSWORD=your_email_password

#### Server Configuration
PORT=8080

#### Cache Configuration (e.g., Redis)
CACHE_PORT=6379
CACHE_HOST=localhost

#### Cloud Storage (e.g., Cloudinary)
CLOUDINARY_NAME=your_cloudinary_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret
CLOUDINARY_UPLOAD_FOLDER=your_upload_folder


## 2. Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/Naoldaba/LoanGuard.git
   cd LoanGuard
2. **Install Dependencies**
    go mod tidy
3. **Run the App**
    go run cmd/main.go

## 3. Postman Documentation
    - https://documenter.getpostman.com/view/31532211/2sAXjM4C46