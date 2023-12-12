package mysql

import (
	"context"
	"database/sql"
	"goKreditPintar/domain"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type mysqlAuthRepository struct {
	Conn *sql.DB
}

// NewMySQLAuthRepository is constructor of MySQL repository
func NewMySQLAuthRepository(Conn *sql.DB) domain.AuthMySQLRepository {
	return &mysqlAuthRepository{Conn}
}

func (db *mysqlAuthRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Client, err error) {
	log.Debug(query)

	rows, err := db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Error(errRow)
		}
	}()

	result = make([]domain.Client, 0)
	for rows.Next() {
		t := domain.Client{}
		err = rows.Scan(
			&t.ID,
			&t.NIK,
			&t.Name,
			&t.BirthDate,
			&t.PhoneNumber,
			&t.Email,
			&t.DtmCrt,
			&t.DtmUpd,
		)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (db *mysqlAuthRepository) GetAllClient(ctx context.Context, page, limit int64, sort, order string, calegId ...string) (client []domain.Client, err error) {
	query := `SELECT id,  nik, name, birth_date, prov_id, kab_id, kec_id, kel_id, address, phone_number, email, active, log, dtm_crt, dtm_upd FROM client`

	if sort != "" {
		query += " ORDER BY " + sort
	}
	if order != "" {
		query += " " + order
	}
	if limit > 0 {
		query += " LIMIT " + strconv.Itoa(int(limit))
	}
	if page > 0 {
		query += " OFFSET " + strconv.Itoa(int(page))
	}

	log.Debug(query)
	client, err = db.fetch(ctx, query)
	if err != nil {
		return
	}

	if len(client) == 0 {
		err = sql.ErrNoRows
	}

	return
}

func (db *mysqlAuthRepository) GetClientByID(ctx context.Context, id int64) (client domain.Client, err error) {
	query := `SELECT id,  nik, name, birth_date, prov_id, kab_id, kec_id, kel_id, address, phone_number, email, active, log, dtm_crt, dtm_upd FROM client WHERE id = ?`

	log.Debug(query + " | " + strconv.FormatInt(id, 10))
	list, err := db.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		client = list[0]
	} else {
		err = sql.ErrNoRows
	}

	return
}

func (db *mysqlAuthRepository) InsertClient(ctx context.Context, client domain.Client) (err error) {
	query := `INSERT INTO client( nik, name, birth_date, prov_id, kab_id, kec_id, kel_id, address, phone_number, email, active, log, dtm_crt, dtm_upd) 
			 VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(
		ctx,
		client.NIK,
		client.Name,
		client.BirthDate,
		client.PhoneNumber,
		client.Email,
	)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	client.ID = lastID
	log.Debug(lastID)

	return
}

func (db *mysqlAuthRepository) LoginAdmin(ctx context.Context, req domain.LoginRequest) (res domain.IsAdminResponse, err error) {
	query := `SELECT user_name, password  FROM admin WHERE user_name = ?`
	log.Debug(query)

	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return
	}

	row := stmt.QueryRowContext(ctx, req.Username)
	err = row.Scan(&res.Username, &res.Password)
	if err != nil {
		log.Error("error in row scan: ", err)
		return
	}

	return

}

func (db *mysqlAuthRepository) LoginConsumer(ctx context.Context, req domain.LoginRequest) (res domain.Client, err error) {
	query := `SELECT id, nik, name, birth_date, phone_number, email, dtm_crt, dtm_upd FROM consumers WHERE phone_number = ? and email = ?`

	list, err := db.fetch(ctx, query, req.PhoneNumber, req.Email)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		err = sql.ErrNoRows
	}

	return
}
