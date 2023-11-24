package user

const (
	scheme      = "public"
	table       = "users"
	tableScheme = scheme + "." + table
)

const (
	sqlAccrualAmountByUserID     = "SELECT COALESCE(SUM(o.accrual), 0) as sum FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE u.id = $1 GROUP BY u.id"
	sqlWithdrawalsAmountByUserID = "SELECT COALESCE(SUM(o.sum), 0) as sum FROM users u LEFT JOIN withdrawals o ON u.id = o.user_id WHERE u.id = $1 GROUP BY u.id"
	sqlGetByLogin                = "SELECT * FROM " + tableScheme + " WHERE login = $1"
	sqlInsertUser                = "INSERT INTO " + tableScheme + " (id, login, password) VALUES (:id, :login, :password)"
)
