# mini-wallet
* This project build with Golang 1.19 for managing a simple mini wallet
* Using JWT for build/construct token
* Using MySQL for the Database

# Project Structure
- main.go
- src
  - config
    - config.go
  - controllers
    - auth_controller.go
    - wallet_controller.go
  - middleware
    - jwt_middleware.go
  - models
    - customer.go
    - wallet.go
  - repositories
    - customer_repository.go
    - wallet_repository.go
  - routes
    - routes.go
  - services
    - auth_service.go
    - wallet_service.go
- test
  - controller_test
    - auth_controller_test.go
    - wallet_controller_test.go
  - mock
    - mock.go

# Explanation
- main.go:The main file to run the application.
- src: The main directory of the application.
  - config: Contains configuration-related files, such as database settings.
  - controllers: Contains controller logic to handle HTTP requests.
  - middleware: Contains middleware for validation, authentication, and authorization.
  - models: Contains the data model definitions used in the application.
  - repositories: Contains the database access logic.
  - routes: Contains the definitions of HTTP routes to be handled by the application.
  - services: Contains the business logic for each application feature.



# How to run/install the Project
- go mod tidy
- config the database (MySQL) on your local. You can look to .env file the configuration of the DB
- start the project 
  - "go run main.go"  and will run in http://localhost:8000


# Project Scope
This Project is handle the all requirements below:
- The user token implemented with JWT.
- When a user initiates, it will automatically create customer data and wallet data.
- Before enabling the wallet, the customer cannot view, add, or use its virtual money.
- If the wallet is already enabled, this endpoint would fail. This endpint should be usable again only if the wallet is disabled.
- The user can make a deposit to add balance to the wallet.
- The user can also make a withdrawal to withdraw balance from the wallet.
- Every deposit or withdrawal transaction will be recorded in the transaction table.
- The user can view all transactions.
- After adding or using virtual money, it is not expected to have the balance immediately updated. The maximum delay for updating the balance is 5 seconds. (with goroutine)

# Database Design and DDL
- You can see the DB Design and the DDL from Folder /documentation


