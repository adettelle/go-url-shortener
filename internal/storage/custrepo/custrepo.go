package custrepo

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

type CustomerRepo struct {
	DB *sql.DB
}

func NewCustomerRepo(db *sql.DB) *CustomerRepo {
	return &CustomerRepo{
		DB: db,
	}
}

// CustomerExistsErr represents an error indicating that a customer with the given login already exists.
type CustomerExistsErr struct {
	login string // The login of the existing customer
}

func (ce *CustomerExistsErr) Error() string {
	return fmt.Sprintf("customer %s already exists", ce.login)
}

func NewCustomerExistsErr(login string) *CustomerExistsErr {
	return &CustomerExistsErr{
		login: login,
	}
}

// IsCustomerExistsErr checks if an error is of type CustomerExistsErr.
func IsCustomerExistsErr(err error) bool {
	var customErr *CustomerExistsErr
	return errors.As(err, &customErr)
}

// AddCustomer adds a new customer to the database.
// Returns the ID of the newly added customer and an error if the operation fails.
func (cr *CustomerRepo) AddCustomer(ctx context.Context,
	userLogin string, password string) (int, error) {

	sqlCustomer := `select count(*) > 0 from customer where login = $1 limit 1;`
	row := cr.DB.QueryRowContext(ctx, sqlCustomer, userLogin)

	// переменная для чтения результата
	var customerExists bool

	err := row.Scan(&customerExists)
	if err != nil {
		return 0, err
	}
	if customerExists {
		return 0, NewCustomerExistsErr(userLogin) // пользователь уже существует
	}

	sqlSt := `insert into customer (login, pwd) values ($1, $2) returning id;`

	var custID int
	row = cr.DB.QueryRowContext(ctx, sqlSt, userLogin, password)
	err = row.Scan(&custID)
	if err != nil {
		return 0, err
	}

	log.Println("Customer is registered.")
	return custID, nil
}

func (cr *CustomerRepo) VerifyUser(ctx context.Context, login string, pass string) (bool, int) {
	log.Println("---------------------")
	if login == "" || pass == "" {
		return false, 0
	}
	// Generate a hash of the provided password
	hashedPassword := sha256.Sum256([]byte(pass))
	hashStringPassword := hex.EncodeToString(hashedPassword[:]) // дополнительно кодируем пароль

	// Retrieve the customer by login
	cust, err := cr.GetCustomerByLogin(ctx, login)
	if err != nil {
		log.Printf("Error in authorization %s", cust.Login)
		return false, 0
	}
	if cust == nil {
		log.Printf("Error in authorization %s, user not found", login)
		return false, 0
	}

	return cust.Password == hashStringPassword, cust.ID
}

type CustomerGetByLogin struct {
	ID       int
	Login    string
	Password string
}

func (cr *CustomerRepo) GetCustomerByLogin(ctx context.Context, login string) (*CustomerGetByLogin, error) {
	sqlSt := `select id, login, pwd from customer where login = $1;`

	row := cr.DB.QueryRowContext(ctx, sqlSt, login)

	var customer CustomerGetByLogin

	err := row.Scan(&customer.ID, &customer.Login, &customer.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // считаем, что это не ошибка, просто не нашли пользователя
		}
		return nil, err
	}
	return &customer, nil
}
