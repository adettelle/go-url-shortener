package dbstorage

import (
	"context"
	"database/sql"
	"log"

	"github.com/adettelle/go-url-shortener/internal/helpers"
	"github.com/adettelle/go-url-shortener/internal/storage"
)

// DBStorage - это имплементация (или реализация) интерфейса Storage
type DBStorage struct {
	Ctx context.Context
	DB  *sql.DB
}

func NewURLRepo(ctx context.Context, db *sql.DB) *DBStorage {
	return &DBStorage{Ctx: ctx, DB: db}
}

func (s *DBStorage) GetOriginalURLByShortURL(shortURL string, custID string) (string, error) { // , custID string
	sqlStatement := "SELECT original_url from url_mapping  where short_url = $1 and customer_id = $2"
	row := s.DB.QueryRowContext(s.Ctx, sqlStatement, shortURL, custID)

	// переменная для чтения результата
	var originalURL string

	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s *DBStorage) AddOriginalURL(originalURL string, custID string) (string, error) { // , custID string
	log.Println("Writing to DB")

	if originalURL == "" {
		return "", &storage.EmptyOriginalURLError{}
	}

	randString, err := helpers.StringWithCharset()
	if err != nil {
		return "", err
	}

	sqlStatement := `insert into url_mapping (short_url, original_url, customer_id) 
		values ($1, $2, $3)` // TODO on conflict ?????

	_, err = s.DB.ExecContext(s.Ctx, sqlStatement, randString, originalURL, custID) // , custID
	if err != nil {
		log.Println("error in adding url:", err)
		return "", err
	}

	log.Println("Saved")
	return randString, nil
}

func (s *DBStorage) Finalize() error {
	return nil
}

// ExistsOriginalURL returns shortURL if originalURL exists
func (s *DBStorage) GetShortURLByOriginalURL(originalURL string) (string, error) { // , custID string
	sqlStatement := "SELECT short_url from url_mapping  where original_url = $1;" //  and customer_id = $2
	row := s.DB.QueryRowContext(s.Ctx, sqlStatement, originalURL)                 // , custID

	// переменная для чтения результата
	var shortURL string

	err := row.Scan(&shortURL)
	if err != nil {
		return "", err
	}

	if shortURL != "" { // errors.Is(err, &storage.OriginalURLExistsErr{})
		return shortURL, storage.NewOriginalURLExistsErr(shortURL, originalURL)
	}

	return shortURL, nil
}

type URL struct {
	ShortURL    string
	OriginalURL string
}

func (s *DBStorage) GetAllURLS(ctx context.Context, custID string) ([]URL, error) { // , userID string
	urls := make([]URL, 0)

	sqlStatement := `select short_url, original_url from url_mapping where customer_id = $1;` //  where customer_id = $1;

	rows, err := s.DB.QueryContext(ctx, sqlStatement, custID) // , userID
	if err != nil || rows.Err() != nil {
		log.Println("error in getting urls:", err)
		return nil, err
	}
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			log.Println("error in scanning url:", err)
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

/*
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
func (s *DBStorage) AddCustomer(ctx context.Context,
	userLogin string, password string) (int, error) {

	sqlCustomer := `select count(*) > 0 from customer where login = $1 limit 1;`
	row := s.DB.QueryRowContext(ctx, sqlCustomer, userLogin)

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
	row = s.DB.QueryRowContext(ctx, sqlSt, userLogin, password)
	err = row.Scan(&custID)
	if err != nil {
		return 0, err
	}

	log.Println("Customer is registered.")
	return custID, nil
}
*/
