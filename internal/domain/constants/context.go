package contextkey

type ContextKey string

const UserIDCtxKey ContextKey = "userID"

const ExchangerCurrencyCtxKey string = "exchangerCurrency"
const ExchangeRateToCurrencyCtxKey string = "exchangeRateToCurrency"

var Currencies = []string{"USD", "EUR", "RUB"}

type OperationType string

const (
	OperationTypeDeposit  OperationType = "DEPOSIT"
	OperationTypeWithdraw OperationType = "WITHDRAW"
	OperationTypeTransfer OperationType = "TRANSFER"
)

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
)
