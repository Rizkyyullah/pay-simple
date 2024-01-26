package configs

const (
  // API Group
  APIGroup = "/api/v1"
  MerchantsGroup = "/merchants"
  CustomersGroup = "/customers"
  
  // Merchants Group
  MerchantRegister = "/auth/register/merchants"
  MerchantLogin = "/auth/login/merchants"
  MerchantLogout = "/auth/logout/merchants"
  
  // Customers Group
  CustomerRegister = "/auth/register"
  CustomerLogin = "/auth/login"
  CustomerTransaction = "/transactions"
  CustomerTransactionWithIDParam = "/transactions/:id"
  CustomerTopup = "/topup"
  
  // Common
  Balance = "/balance"
  Products = "/products"
  ProductsWithIDParam = "/products/:id"
  Transfer = "/transactions/transfer"
)