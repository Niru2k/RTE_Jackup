# OnlinePurchase Project
This repository contains a Go language web application for managing online orders and products. The application provides APIs for user signup, login, managing products, placing orders, and checking order status. It is built using the Echo framework and GORM for database interactions.

## Technologies used
The project is built using the following technologies:
- **Golang**  : The backend is written in Go (Golang), a statically typed, compiled language.
- **Echo**   : The Echo web framework is used to create RESTful APIs and handle HTTP requests.
- **JWT**     : JSON Web Tokens are used for secure user authentication and authorization.
- **bcrypt**  : Passwords are stored securely in hashed form using the bcrypt hashing algorithm.
- **Postgres**: Here, users data and Post articles data are handled in Postgres SQL.

## Project Structure
The project is organized into several packages, each responsible for specific functionalities:
- `handlers`  : Contains the HTTP request handlers for different API endpoints.
- `logs`      : Custom package for logging.
- `middleware`: Custom middleware for handling authentication and authorization.
- `models`    : Defines the data models used in the application.
- `repository`: Contains functions for interacting with the database.
- `drivers`   : Contains functions for establish a connection to database.
- `helper`    : Custom package that contains all the constants.
- `Lookup`    : Contains functions for checking the database updations.

## Endpoints
The following endpoints are available in the application:

### User Signup
- `POST /signup`: Sign up a new user with the required details such as username , email, password, and role (admin or user).

### User Login
- `POST /login`: Authenticate a user with email and password and return a JWT token for further authentication.

### Product Management
- `POST /product`: Add a new product with details such as brand name, product price, RAM capacity, etc. (Admin access required)
- `GET /products`: Get a list of all products.
- `PUT /product/:product_id`: Update product details by product ID. (Admin access required)
- `DELETE /product/:product_id`: Delete a product by product ID. (Admin access required)

### Order Management
- `POST /order`: Place a new order with details such as brand name, product price, RAM capacity, etc. (User access required)
- `DELETE /order/:order_id`: Cancel an order by order ID. (User access required)
- `GET /orders`: Get a list of all orders for the current user. (User access required)
- `POST /payment/:order_id`: Make a payment for an order with the specified order ID. (User access required)
- `PUT /orderstatus/:order_id`: Update the status of an order by order ID. (Admin access required)
- `GET /orderstatus/:order_id`: Get the status of an order by order ID.
- `GET /orderstatuses`: Get a list of all order statuses. (Admin access required)

## Authentication
The application uses JWT (JSON Web Token) for authentication. To access the protected endpoints, users need to include the JWT token in the `Authorization` header of the request.

### User Authentication
- For user authentication, a valid JWT token obtained during the login process should be included in the header as follows:

### Admin Authentication
- For admin authentication, the same process applies with a valid JWT token obtained during the login process for an admin user.

## Error Handling
The application handles various error scenarios and provides appropriate error responses with corresponding status codes and messages.
