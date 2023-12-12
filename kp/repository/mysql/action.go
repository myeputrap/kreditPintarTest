package mysql

import (
	"context"
	"database/sql"
	"errors"
	"goKreditPintar/domain"
	"sync"

	"strconv"

	log "github.com/sirupsen/logrus"
)

type mysqlActionRepository struct {
	Conn *sql.DB
	mu   sync.Mutex
}

// NewMySQLActionRepository is constructor of MySQL repository
func NewMySQLActionRepository(Conn *sql.DB) domain.ActionMySQLRepository {
	return &mysqlActionRepository{
		Conn: Conn,
		mu:   sync.Mutex{},
	}
}

func (db *mysqlActionRepository) GetConsumer(ctx context.Context, req domain.GetAllConsumerRequest) (res []domain.Consumer, err error) {
	log.Debug("start")
	query := `SELECT id, nik, name, place_of_birth, birth_date, salary, email, phone_number, dtm_crt, dtm_upd FROM consumers WHERE 1 = 1`
	var params []interface{}

	if req.Name != "" {
		query += " AND name = ?"
		params = append(params, req.Name)
	}
	if req.Email != "" {
		query += " AND email = ?"
		params = append(params, req.Email)
	}
	if req.Salary != "" {
		query += " AND salary > ?"
		params = append(params, req.Salary)
	}

	if req.Sort != "" {
		query += " ORDER BY " + req.Sort
	}

	if req.Order != "" {
		query += " " + req.Order
	}

	if req.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(int(req.Limit))
	}

	if req.Page > 0 {
		query += " OFFSET " + strconv.Itoa(int(req.Page))
	}

	log.Debug(query)
	db.mu.Lock()
	defer db.mu.Unlock()
	rows, err := db.Conn.QueryContext(ctx, query, params...)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i domain.Consumer
		err = rows.Scan(&i.ID, &i.Nik, &i.Name, &i.PlaceOfBirth, &i.BirthDate, &i.Salary, &i.Email, &i.PhoneNumber, &i.DtmCrt, &i.DtmUpd)
		if err != nil {
			log.Error(err)
			return
		}

		res = append(res, i)
	}
	return
}

func (db *mysqlActionRepository) GetCreditCard(ctx context.Context, req domain.GetAllCreditCardRequest) (res []domain.ConsumerCreditCard, err error) {
	log.Debug("start")
	query := `SELECT id, consumer_id, card_number, expiration_date, cvv, credit_limit, current_balance, dtm_crt, dtm_upd FROM credit_cards WHERE 1 = 1`
	var params []interface{}

	if req.CardNumber != "" {
		query += " AND card_number = ? "
		params = append(params, req.CardNumber)
	}
	if req.CreditLimit != "" {
		query += " AND credit_limit > ? "
		params = append(params, req.CreditLimit)
	}
	if req.CurrentBalance != "" {
		query += " AND current_balance > ? "
		params = append(params, req.CurrentBalance)
	}

	if req.ExpirationDate != "" {
		query += " AND expiration_date = ? "
		params = append(params, req.ExpirationDate)
	}

	if req.Sort != "" {
		query += " ORDER BY " + req.Sort
	}

	if req.Order != "" {
		query += " " + req.Order
	}

	if req.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(int(req.Limit))
	}

	if req.Page > 0 {
		query += " OFFSET " + strconv.Itoa(int(req.Page))
	}

	log.Debug(query)
	db.mu.Lock()
	defer db.mu.Unlock()
	rows, err := db.Conn.QueryContext(ctx, query, params...)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i domain.ConsumerCreditCard
		err = rows.Scan(&i.ID, &i.ConsumerID, &i.CardNumber, &i.ExpirationDate, &i.CVV, &i.CreditLimit, &i.CurrentBalance, &i.DtmCrt, &i.DtmUpd)
		if err != nil {
			log.Error(err)
			return
		}

		res = append(res, i)
	}
	return
}

func (db *mysqlActionRepository) GetBilling(ctx context.Context, req domain.GetAllBillingRequest) (res []domain.Billing, err error) {
	log.Debug("start")
	query := `SELECT id, consumer_id, transaction_id, bill_amount, due_date, status, dtm_crt, dtm_upd FROM billing WHERE 1 = 1`
	var params []interface{}

	if req.BillAmount != "" {
		query += " AND bill_amount = ?"
		params = append(params, req.BillAmount)
	}
	if req.DueDate != "" {
		query += " AND due_date = ?"
		params = append(params, req.DueDate)
	}
	if req.Status != "" {
		query += " AND status = ?"
		params = append(params, req.Status)
	}

	if req.Sort != "" {
		query += " ORDER BY " + req.Sort
	}

	if req.Order != "" {
		query += " " + req.Order
	}

	if req.Limit > 0 {
		query += " LIMIT " + strconv.Itoa(int(req.Limit))
	}

	if req.Page > 0 {
		query += " OFFSET " + strconv.Itoa(int(req.Page))
	}

	log.Debug(query)
	db.mu.Lock()
	defer db.mu.Unlock()
	rows, err := db.Conn.QueryContext(ctx, query, params...)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i domain.Billing
		err = rows.Scan(&i.ID, &i.ConsumerID, &i.TransactionID, &i.BillAmount, &i.DueDate, &i.Status, &i.DtmCrt, &i.DtmUpd)
		if err != nil {
			log.Error(err)
			return
		}

		res = append(res, i)
	}
	return
}

func (db *mysqlActionRepository) CountConsumer(ctx context.Context, req domain.GetAllConsumerRequest) (sum int64, err error) {
	query := `SELECT COUNT(id) FROM credit_cards WHERE`
	var params []interface{}

	if req.Name != "" {
		query += "name = ? AND"
		params = append(params, req.Name)
	}
	if req.Email != "" {
		query += "email = ? AND"
		params = append(params, req.Email)
	}
	if req.Salary != "" {
		query += "salary > ? "
		params = append(params, req.Salary)
	}

	if query[len(query)-5:] == "WHERE" {
		query = query[:len(query)-5]
	}

	if query[len(query)-3:] == "AND" {
		query = query[:len(query)-3]
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, params...)

	err = row.Scan(&sum)
	if err != nil {
		log.Error(err)
	}

	return
}

func (ar *mysqlActionRepository) CountCreditCard(ctx context.Context, req domain.GetAllCreditCardRequest) (sum int64, err error) {
	query := `SELECT COUNT(id) FROM credit_cards WHERE`
	var params []interface{}

	if req.CardNumber != "" {
		query += " card_number = ? AND"
		params = append(params, req.CardNumber)
	}
	if req.CreditLimit != "" {
		query += " credit_limit > ? AND"
		params = append(params, req.CreditLimit)
	}
	if req.CurrentBalance != "" {
		query += " current_balance > ? AND"
		params = append(params, req.CurrentBalance)
	}

	if req.ExpirationDate != "" {
		query += " expiration_date = ? AND"
		params = append(params, req.ExpirationDate)
	}

	if query[len(query)-5:] == "WHERE" {
		query = query[:len(query)-5]
	}

	if query[len(query)-3:] == "AND" {
		query = query[:len(query)-3]
	}
	log.Debug(query)
	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, params...)

	err = row.Scan(&sum)
	if err != nil {
		log.Error(err)
	}

	return
}

func (ar *mysqlActionRepository) CountBilling(ctx context.Context, req domain.GetAllBillingRequest) (sum int64, err error) {
	query := `SELECT COUNT(id) FROM billing WHERE`
	var params []interface{}

	if req.BillAmount != "" {
		query += " bill_amount = ? AND"
		params = append(params, req.BillAmount)
	}
	if req.DueDate != "" {
		query += " due_date = ? AND"
		params = append(params, req.DueDate)
	}
	if req.Status != "" {
		query += " status = ? "
		params = append(params, req.Status)
	}

	if query[len(query)-5:] == "WHERE" {
		query = query[:len(query)-5]
	}

	if query[len(query)-3:] == "AND" {
		query = query[:len(query)-3]
	}

	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, params...)

	err = row.Scan(&sum)
	if err != nil {
		log.Error(err)
	}

	return
}

func (ar *mysqlActionRepository) CheckValidate(ctx context.Context, coloumn string, value string) (err error) {
	query := `SELECT COUNT(id) FROM consumers WHERE ` + coloumn + ` = ?`

	log.Debug("Query : " + query)
	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}
	var count int64
	row := stmt.QueryRowContext(ctx, value)
	err = row.Scan(&count)
	if err != nil {
		log.Error("error in row scan: ", err)
		return
	}

	if count > 0 {
		err = errors.New("data exist")
		return
	}
	return
}

func (ar *mysqlActionRepository) CheckValidateCC(ctx context.Context, coloumn string, value string) (count int, err error) {
	query := `SELECT COUNT(id) FROM credit_cards WHERE ` + coloumn + ` = ?`

	log.Debug("Query : " + query)
	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}
	row := stmt.QueryRowContext(ctx, value)
	err = row.Scan(&count)
	if err != nil {
		log.Error("error in row scan: ", err)
		return
	}

	return
}

func (db *mysqlActionRepository) PostConsumer(ctx context.Context, req domain.Consumer) (err error) {
	query := `INSERT INTO consumers (nik, name, birth_date, place_of_birth, phone_number, email, salary, dtm_crt, dtm_upd)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, req.Nik, req.Name, req.BirthDate, req.PlaceOfBirth, req.PhoneNumber, req.Email, req.Salary)
	if err != nil {
		return
	}

	return
}

func (db *mysqlActionRepository) GetConsumerDetail(ctx context.Context, req int) (res domain.Consumer, err error) {
	query := `SELECT id, nik, name, place_of_birth, birth_date, salary, email, phone_number, dtm_crt, dtm_upd FROM consumers WHERE id = ?`
	log.Debug(query)

	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	row := stmt.QueryRowContext(ctx, req)
	err = row.Scan(&res.ID, &res.Nik, &res.Name, &res.PlaceOfBirth, &res.BirthDate, &res.Salary, &res.Email, &res.PhoneNumber, &res.DtmCrt, &res.DtmUpd)
	if err != nil {
		return
	}
	return
}
func (db *mysqlActionRepository) GetCreditCardDetail(ctx context.Context, req int) (res domain.ConsumerCreditCard, err error) {
	query := `SELECT id, consumer_id, card_number, expiration_date, cvv, credit_limit, current_balance, dtm_crt, dtm_upd FROM credit_cards WHERE id = ?`
	log.Debug(query)

	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	row := stmt.QueryRowContext(ctx, req)
	err = row.Scan(&res.ID, &res.ConsumerID, &res.CardNumber, &res.ExpirationDate, &res.CVV, &res.CreditLimit, &res.CurrentBalance, &res.DtmCrt, &res.DtmUpd)
	if err != nil {
		return
	}
	return
}

func (db *mysqlActionRepository) GetConsumerByParameter(ctx context.Context, coloumn string, value string) (res domain.Consumer, err error) {
	query := `SELECT id, nik, name, place_of_birth, birth_date, salary, email, phone_number, dtm_crt, dtm_upd FROM consumers WHERE ` + coloumn + ` = ?`
	log.Debug(query)

	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	row := stmt.QueryRowContext(ctx, value)
	err = row.Scan(&res.ID, &res.Nik, &res.Name, &res.PlaceOfBirth, &res.BirthDate, &res.Salary, &res.Email, &res.PhoneNumber, &res.DtmCrt, &res.DtmUpd)
	if err != nil {
		return
	}
	return
}

func (db *mysqlActionRepository) PostConsumerCreditCard(ctx context.Context, req domain.ConsumerCreditCard) (err error) {
	query := `INSERT INTO credit_cards (consumer_id, card_number, expiration_date, cvv, credit_limit, current_balance, dtm_crt, dtm_upd)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, req.ConsumerID, req.CardNumber, req.ExpirationDate, req.CVV, req.CreditLimit, req.CurrentBalance)
	if err != nil {
		return
	}

	return
}

func (db *mysqlActionRepository) PostTransactionCredit(ctx context.Context, req domain.TransactionCreditCard) (count int64, err error) {
	query := `INSERT INTO transactions (consumer_id, contract_number, otr, admin_fee, installment_count, interest_amount, purchase_amount, asset_name, status, dtm_crt, dtm_upd)
		VALUES (?, ?, ?, ?, ?, ?, ?,?, ?, NOW(), NOW())`

	db.mu.Lock()
	defer db.mu.Unlock()

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.ExecContext(ctx, req.ConsumerID, req.ContractNumber, req.OTR, req.AdminFee, req.InstallmentCount, req.InterestAmount, req.PurchaseAmount, req.AssetName, "active")
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	defer db.mu.Unlock()

	return lastID, nil

}

func (db *mysqlActionRepository) PostBillingCredit(ctx context.Context, req domain.Billing) (err error) {
	query := `INSERT INTO billing (consumer_id, transaction_id, bill_amount, due_date, status, dtm_crt, dtm_upd)
           VALUES (?, ?, ?, ?, ?,NOW(), NOW())`

	db.mu.Lock()
	defer db.mu.Unlock()
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, req.ConsumerID, req.TransactionID, req.BillAmount, req.DueDate, req.Status)
	if err != nil {
		return
	}

	return
}

func (ar *mysqlActionRepository) UpdateCreditBalance(ctx context.Context, bill float64, id string) (err error) {
	query := `UPDATE credit_cards SET current_balance = ?, dtm_upd = NOW() WHERE consumer_id = ?`
	log.Debug(query)

	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(
		ctx,
		bill,
		id,
	)

	if err != nil {
		return
	}

	return
}

func (ar *mysqlActionRepository) CountCreditBalance(ctx context.Context, id string) (res float64, err error) {
	query := `SELECT COUNT(bill_amount) FROM billing WHERE consumer_id = ? AND status = Pending`

	log.Debug("Query : " + query)

	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}
	row := stmt.QueryRowContext(ctx, id)
	err = row.Scan(&res)
	if err != nil {
		log.Error("error in row scan: ", err)
		return
	}

	return
}

func (ar *mysqlActionRepository) PatchBilling(ctx context.Context, id string) (err error) {
	query := `UPDATE status SET status = ?, dtm_upd = NOW() WHERE id = ?`
	log.Debug(query)

	ar.mu.Lock()
	defer ar.mu.Unlock()
	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(
		ctx,
		"Paid",
		id,
	)

	if err != nil {
		return
	}

	return
}

func (ar *mysqlActionRepository) CheckValidateBilling(ctx context.Context, coloumn string, value string) (count int, err error) {
	query := `SELECT COUNT(id) FROM billing WHERE ` + coloumn + ` = ?`

	log.Debug("Query : " + query)

	ar.mu.Lock()
	defer ar.mu.Unlock()

	stmt, err := ar.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}
	row := stmt.QueryRowContext(ctx, value)
	err = row.Scan(&count)
	if err != nil {
		log.Error("error in row scan: ", err)
		return
	}

	return
}
