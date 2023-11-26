# BasicTrade Application - V1

BasicTrade is an application for managing products and variants, equipped with authentication, CRUD operations, and Cloudinary integration.

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Key Features](#key-features)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [API Endpoints](#api-endpoints)
- [Deployment](#deployment)
- [Contributions](#contributions)
- [License](#license)

## Introduction

BasicTrade is a platform designed to streamline product and variant management for admins. Admins can register and log in to access these features. With a secure authentication process in place, the application ensures the security of user accounts. The CRUD operations enable admins to easily create, update, and delete product and variant information.

## Prerequisites

List any prerequisites or dependencies that need to be installed before setting up the application.

- Go programming language
- MySQL database
- Cloudinary account


## Key Features

- **Authentication:** Secure login and register processes for admins.
- **CRUD Operations:** Create, Read, Update, and Delete operations for products and variants.
- **Photo Management:** Cloudinary integration for efficient storage and retrieval of product photos.
- **Modular Structure:** The application is organized into distinct modules for easy development and maintenance.

## Technology Stack

- **Gin Gonic:** Web framework for building APIs in Go.
- **Gorm:** ORM (Object Relational Mapper) for database interactions.
- **JWT (JSON Web Tokens):** Used for secure authentication.
- **Cloudinary:** Cloud-based image and video management service.

## Project Structure
```bash
BasicTrade
|-- controllers
|-- middleware
|-- models
|-- routes
|-- utils
|-- helpers
|-- main.go
```

## Getting Started
Provide steps to install and set up the BasicTrade application on a local machine.
1. Clone this repository.
```bash
# Clone the repository
git clone https://github.com/diazamaliana/basictrade.git

# Change into the project directory
cd basictrade

# Install dependencies
go mod download
```
2. Set up the necessary environment variables, such as Cloudinary credentials and database connection details.
```bash
HOST="your-host"
DB_USER="your-db-username"
DB_PASSWORD="your-db-password"
DB_NAME="your-db-name"
DB_PORT="your-db-port"
CLOUDINARY_CLOUD_NAME="your-cloudinary-name"
CLOUDINARY_API_KEY="your-cloudinary-api-key"
CLOUDINARY_API_SECRET="your-cloudinary-api-secret"
CLOUDINARY_UPLOAD_FOLDER="your-cloudinary-folder-name"
JWT_SECRET_KEY="your-jwt-secret-key"
PORT="5050"
```

3. Run the application using `go run main.go`.
4. Access the application at `http://localhost:5050` (or the specified port).

## API Endpoints

1. **POST /auth/register:** Register an admin.
2. **POST /auth/login:** Log in as an admin.
3. **GET /products:** Get all products.
4. **POST /products:** Create a product.
5. **PUT /products/:productUUID:** Update product details.
6. **DELETE /products/:productUUID:** Delete a product.
7. **GET /products/:productUUID:** Get product details.
8. **GET /products/variants:** Get all variants.
9. **POST /products/variants/:variantUUID:** Create a variant.
10. **PUT /products/variants/:variantUUID:** Update variant details.
11. **DELETE /products/variants/:variantUUID:** Delete a variant.
12. **GET /products/variants/:variantUUID:** Get variant details.

## Deployment

The BasicTrade application can be deployed on the [Railway](https://railway.app/) platform. Ensure you configure the necessary environment variables for successful deployment. 

## Contributions

Contributions are welcome! If you find any issues or have suggestions for improvement, feel free to create a pull request or raise an issue.

## License

This project is licensed under the [MIT License](LICENSE).
