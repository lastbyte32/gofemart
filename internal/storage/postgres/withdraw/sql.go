package withdraw

const (
	scheme      = "public"
	table       = "withdrawals"
	tableScheme = scheme + "." + table
)
const (
	sqlSumAccrualByUserID = "SELECT * FROM " + tableScheme + " WHERE user_id = $1"

	sqlGetByUserID = "SELECT * FROM " + tableScheme + " WHERE user_id = $1"
	sqlInsert      = "INSERT INTO " + tableScheme + " (user_id, order_number, sum) VALUES (:user_id, :order_number, :sum)"
)
