package order

const (
	scheme      = "public"
	table       = "orders"
	tableScheme = scheme + "." + table
)
const (
	sqlGetByUserIdAndNumber = "SELECT * FROM " + tableScheme + " WHERE user_id = $1 AND number = $2"
	sqlGetByUserID          = "SELECT * FROM " + tableScheme + " WHERE user_id = $1"
	sqlGetByNumber          = "SELECT * FROM " + tableScheme + " WHERE number = $1"
	sqlInsert               = "INSERT INTO " + tableScheme + " (number, user_id, status) VALUES (:number, :user_id, :status)"
)
