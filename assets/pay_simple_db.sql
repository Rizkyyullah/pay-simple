CREATE DATABASE pay_simple_db;

CREATE TYPE role_type AS ENUM('CUSTOMER', 'MERCHANT');
CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');
CREATE TYPE cashflow_type AS ENUM ('MONEY_IN', 'MONEY_OUT');

CREATE SEQUENCE users_id_seq;
CREATE SEQUENCE products_id_seq;
CREATE SEQUENCE transactions_id_seq;
CREATE SEQUENCE transaction_details_id_seq;

CREATE TABLE users (
  id VARCHAR(10) PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(100) NOT NULL,
  balance INT DEFAULT 10000000,
  email VARCHAR(255) UNIQUE NOT NULL,
  phone_number CHAR(13),
  password VARCHAR(255) NOT NULL,
  role role_type NOT NULL,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
  id VARCHAR(10) PRIMARY KEY NOT NULL,
  user_id VARCHAR(10) NOT NULL,
  product_name VARCHAR(100) NOT NULL,
  description TEXT,
  stock INT DEFAULT 0,
  price INT DEFAULT 0,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE transactions (
  id VARCHAR(10) PRIMARY KEY NOT NULL,
  user_id VARCHAR(10) NOT NULL,
  transaction_date TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  transaction_type transaction_type NOT NULL,
  paid_status BOOLEAN DEFAULT false,
  cashflow cashflow_type NOT NULL,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE transaction_details (
  id VARCHAR(10) PRIMARY KEY NOT NULL,
  transaction_id VARCHAR(10) NOT NULL,
  product_id VARCHAR(10) NOT NULL,
  quantity INT DEFAULT 1,
  total_price INT NOT NULL,
  created_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP
);
