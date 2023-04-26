package order

const (
	scheme      = "public"
	table       = "orders"
	tableScheme = scheme + "." + table
)
const (
	sqlUpdate              = "UPDATE " + tableScheme + " SET status = :status, accrual = :accrual WHERE number = :number"
	sqlGetOrdersUnpocessed = "SELECT number FROM " + tableScheme + " WHERE status=$1 OR status=$2 ORDER BY uploaded_at"
	sqlGetByUserID         = "SELECT * FROM " + tableScheme + " WHERE user_id = $1"
	sqlGetByNumber         = "SELECT * FROM " + tableScheme + " WHERE number = $1"
	sqlInsert              = "INSERT INTO " + tableScheme + " (number, user_id, status) VALUES (:number, :user_id, :status)"
)
