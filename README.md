# PaySimple API Documentation

The Pay-Simple API provides endpoints for a simulation payment system. It includes authentication, user registration, profile management, balance inquiries, product listings, and transaction history.

## How to run an application

### Prerequisite
1. Have a [Git](https://git-scm.com) that has been installed on your computer
2. [Go](https://go.dev) Language has been installed and [PostgreSQL](https://postgresql.com)
3. Have this clone repo
4. Or please access endpoint [https://pay-simple.adaptable.app/api/v1/{resource}](https://pay-simple.adaptable.app/api/v1) using a REST Client such as Postman or Insomnia to make a request without having to do this clone repo

Note*: I recommend running the application locally, as there were some issues during the deployment process and I haven't tried all the endpoints that were created. But for the process of registering, logging in, logging out, and viewing all products it runs smoothly on the adaptable server.

### Installation
If you want to clone this repo and run it on your local computer. Follow these steps
```bash
# clone this repository using command
git clone https://github.com/Rizkyyullah/pay-simple.git

# cd to it
cd pay-simple

# you must enter key and value pair on the .env file, which is used to retrieve database configurations, server ports etc.
# you can copy the contents of .env.example file, create the .env file and fill in the value according to your configuration

# after that run this command to start application
go run main.go
```

## Authentication

### Logout [/api/v1/auth/logout] - GET

#### Description

Logout the currently authenticated user.

- Request Method: GET
- Response:

```json
{
  "status": "Success",
  "code": 200,
  "message": "Logout Successfully"
}
```

## Merchant Authentication and Registration

### Register Merchants [/api/v1/auth/register/merchants] - POST

#### Description

Register a new merchant. For balance I gave automatically 10,000,000 to facilitate the simulation process. You can also do topup and transfer to other users

Request:

- Method: `POST`
- Body:

```json
{
	"name": "Merchant 1",
	"username": "merchant 1",
	"email": "merchant1@mail.com",
	"phoneNumber": "0812345678910",
	"password": "password",
	"confirmPassword": "password"
}
```

- Response Body:
```json
{
  "meta": {
    "status": "Success",
    "code": 201,
    "message": "Register Successfully",
    "createdAt": "Thursday, 25 January 2024 18:57:13 WIB"
  },
  "data": {
    "id": "USR-0001",
    "name": "Merchant 1",
    "username": "merchant 1",
    "balance": 10000000,
    "email": "merchant1@mail.com",
    "phoneNumber": "0812345678910",
    "role": "MERCHANT"
  }
}
```

### Login Merchants [/api/v1/auth/login/merchants] - POST

#### Description

Authenticate a merchant.

Request:

- Method: `POST`
- Body:

```json
{
	"email": "merchant1@mail.com",
	"password": "password"
}
```

Response:

```json
{
  "status": "Success",
  "code": 200,
  "message": "Login Successfully"
}
```

## Customer Authentication and Registration

### Login Customers [/api/v1/auth/login] - POST

#### Description

Register a new customer.

Request:

- Method: `POST`
- Body:

```json
{
	"name": "Customer 1",
	"username": "customer 1",
	"email": "customer1@mail.com",
	"phoneNumber": "0812345678911",
	"password": "password",
	"confirmPassword": "password"
}
```

- Response Body:
```json
{
  "meta": {
    "status": "Success",
    "code": 201,
    "message": "Register Successfully",
    "createdAt": "Thursday, 25 January 2024 18:57:13 WIB"
  },
  "data": {
    "id": "USR-0002",
    "name": "Customer 2",
    "username": "customer 2",
    "balance": 10000000,
    "email": "customer2@mail.com",
    "phoneNumber": "0812345678911",
    "role": "CUSTOMER"
  }
}
```

### Login Customers [/api/v1/auth/login] - POST

#### Description

Authenticate a customer.

Request:

- Method: `POST`
- Body:

```json
{
	"email": "customer2@mail.com",
	"password": "password"
}
```

Response:

```json
{
  "status": "Success",
  "code": 200,
  "message": "Login Successfully"
}
```

## User Profile

### User Profile [/api/v1/profile] - GET

#### Description

Retrieve the profile of the authenticated user.

Request:

- Method: `GET`
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Your profile"
  },
  "data": {
    "id": "USR-0002",
    "name": "Customer 2",
    "username": "customer 2",
    "balance": 10000,
    "email": "customer2@mail.com",
    "phoneNumber": "0812345678910",
    "role": "CUSTOMER"
  }
}
```

## Merchant Balance

### Merchant Balance [/api/v1/merchants/balance] - GET

#### Description

Retrieve the balance of the authenticated merchant.

Request:

- Method: `GET`
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Your current balance"
  },
  "data": {
    "balance": 10000000
  }
}
```

## Customer Balance

### Customer Balance [/api/v1/customers/balance] - GET

#### Description

Retrieve the balance of the authenticated customers.

Request:

- Method: `GET`
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Your current balance"
  },
  "data": {
    "balance": 10000000
  }
}
```

## Product Listing

### Products [/api/v1/products] - GET

#### Description

Retrieve a list of available products.

Request:

- Method: `GET`
- Parameters:
    - page (How good page you want to access)
    - size (How many data do you want to display on one page)
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Get all products successfully"
  },
  "data": [
    {
      "id": "PRD-0001",
      "merchant": {
        "name": "Merchant 1",
        "username": "merchant 1",
        "email": "merchant1@mail.com",
        "phoneNumber": "0812345678910"
      },
      "productName": "Product A",
      "description": "Desc A",
      "stock": 90,
      "price": 40000,
      "createdAt": "2024-01-25T09:19:55Z",
      "updatedAt": "2024-01-25T09:19:55Z"
    }
  ],
  "paging": {
    "page": 1,
    "rowsPerPage": 5,
    "totalRows": 1,
    "totalPages": 1
  }
}
```

### Product By ID [/api/v1/products/:id] - GET

#### Description

Retrieve a specific product by id.

Request:

- Method: `GET`
- Parameters: id (Product ID)
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Get all products successfully"
  },
  "data": {
      "id": "PRD-0001",
      "merchant": {
        "name": "Merchant 1",
        "username": "merchant 1",
        "email": "merchant1@mail.com",
        "phoneNumber": "0812345678910"
      },
      "productName": "Product A",
      "description": "Desc A",
      "stock": 90,
      "price": 40000,
      "createdAt": "2024-01-25T09:19:55Z",
      "updatedAt": "2024-01-25T09:19:55Z"
    }
}
```

## Merchant Products

### Merchant Products Listing [/api/v1/merchants/products] - GET

#### Description

Retrieve a list of products associated with the authenticated merchant.

Request:

- Method: `GET`
- Parameters:
    - page (How good page you want to access)
    - size (How many data do you want to display on one page)
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Your products"
  },
  "data": [
    {
      "id": "PRD-0001",
      "productName": "Product A",
      "description": "Desc 1",
      "stock": 90,
      "price": 40000,
      "createdAt": "2024-01-25T09:19:55Z",
      "updatedAt": "2024-01-25T09:19:55Z"
    }
  ],
  "paging": {
    "page": 1,
    "rowsPerPage": 5,
    "totalRows": 1,
    "totalPages": 1
  }
}
```

### Add Product by Merchant [/api/v1/merchants/products] - POST

#### Description

Add a new product to the inventory of the authenticated merchant.

Request:

- Method: `POST`
- Body:

```json
{
	"productName": "Product A Mark II",
	"description": "Desc Product A Mark II",
	"stock": 10,
	"price": 10000
}
```

- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 201,
    "message": "Create product successfully",
    "createdAt": "Friday, 2 January 2024 02:36:45 WIB"
  },
  "data": {
    "id": "PRD-0002",
    "merchant": {
      "name": "Merchant 1",
      "username": "merchant 1",
      "email": "merchant1@mail.com"
    },
    "productName": "Product A Mark II",
    "description": "Desc Product A Mark II",
    "stock": 10,
    "price": 10000
  }
}
```

### Merchant Product Details [/api/v1/merchants/products/:id] - DELETE

#### Description

Delete product based on ID by authenticated merchant.

Request:

- Method: `DELETE`
- Parameters: id (Product ID)

Response:

```json
{
  "status": "Success",
  "code": 200,
  "message": "Delete product successfully"
}

```

## Customer Transactions

### Customer Transaction History [/api/v1/customers/transactions] - GET

#### Description

Retrieve the transaction history of the authenticated customer.

Request:

- Method: `GET`
- Parameters:
    - page (How good page you want to access)
    - size (How many data do you want to display on one page)

Response:
```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Your transactions history"
  },
  "data": [
    {
      "id": "TRX-0001",
      "transactionDate": "Friday, 26 January 2024",
      "transactionType": "DEBIT",
      "cashflow": "MONEY_OUT",
      "createdAt": "Friday, 26 January 2024 10:52:11 WIB"
    }
  ],
  "paging": {
    "page": 1,
    "rowsPerPage": 5,
    "totalRows": 1,
    "totalPages": 1
  }
}
```

### Customer Transaction based on ID [/api/v1/customers/transactions/:id] - GET

#### Description

Retrieve details of a specific transaction of the authenticated customer.

Request:

- Method: `GET`
- Parameters: id (Transaction ID)
- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Your transaction history"
  },
  "data": {
    "id": "TRX-0001",
    "transactionDate": "Thursday, 25 January 2024",
    "transactionType": "DEBIT",
    "paidStatus": true,
    "cashFlow": "MONEY_OUT",
    "createdAt": "Thursday, 25 January 2024 16:12:14 WIB"
  }
}
```

### Perform Customer Transaction [/api/v1/customers/transactions] - POST

#### Description

Initiate a new transaction for the authenticated customer.

Request:

- Method: `POST`
- Body:

```json
{
	"products": [
		{
			"id": "PRD-0001",
			"merchantId": "USR-0001",
			"quantity": 1
		}
	]
}
```

## Customer Balance Top-Up

### Customer Balance Top-Up [/api/v1/customers/balance/topup] - POST

#### Description

Top up the balance of the authenticated customer.

Request:

- Method: `POST`
- Body:

```json
{
	"amount": 150000
}
```

- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Topup successfully"
  },
  "data": {
    "balance": 1120000
  }
}
```

### Customer Balance Transfer [/api/v1/customers/transactions/transfer] - POST

#### Description

Transfer a specified amount to another users.

Request:

- Method: `POST`
- Body:

```json
{
	"toUserId": "USR-0001",
	"amount": 150000
}
```

- Response:

```json
{
  "meta": {
    "status": "Success",
    "code": 200,
    "message": "Transfer successfully"
  },
  "data": {
    "balance": 970000
  }
}
```
