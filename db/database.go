package db

import (
	"database/sql"
	"fmt"

	"github.com/Narutchai01/ProjectC-BE/data"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*data.Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(*data.Account) error
	GetAccountByID(int) (*data.Account, error)
	GetAccounts() ([]*data.Account, error)
	GetAccountByEmail(string) (*data.Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostGresStore() (*PostgresStore, error) {
	connStr := "user=admin password=admin123 dbname=kmutt sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS Account (
		user_id SERIAL PRIMARY KEY,
		name varchar(100),
		email varchar(100),
		encrypted_password bpchar,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *data.Account) (int, error) {
	query := `INSERT INTO Account 
		(name, email, encrypted_password, created_at)
		VALUES ($1, $2, $3, $4) RETURNING user_id`

	id := 0
	err := s.db.QueryRow(
		query,
		acc.Name, acc.Email, acc.EncryptedPassword, acc.CreatedAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *PostgresStore) UpdateAccount(*data.Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM Account WHERE user_id = $1", id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*data.Account, error) {
	rows, err := s.db.Query("SELECT * FROM Account WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByEmail(email string) (*data.Account, error) {
	rows, err := s.db.Query("SELECT * FROM Account WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with email [%s] not found", email)
}

func (s *PostgresStore) GetAccounts() ([]*data.Account, error) {
	rows, err := s.db.Query("SELECT * FROM Account")
	if err != nil {
		return nil, err
	}

	accounts := []*data.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*data.Account, error) {
	account := new(data.Account)
	err := rows.Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.EncryptedPassword,
		&account.CreatedAt)

	return account, err
}