package links

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// postgres репозиторий
type PostgresRepository struct {
	db *sql.DB
}

// Создаем репозиторий БД
func NewPostgresRepository(databaseDSN string) (*PostgresRepository, error) {

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	// TODO: Нужно разобраться. Если оставляю этот код, то соединене закрывается и нет доступа к БД. Где это нужно делать?
	//defer db.Close()

	// создаю контекст для пинга
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// пингую БД. Если не отвечает, то возвращаю ошибку
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	// создаю контекст для запроса на создание таблицы
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	// если в БД нет таблицы, то создаю ее
	query := "CREATE TABLE IF NOT EXISTS links(" +
		"link_key TEXT," +
		"link_url TEXT UNIQUE NOT NULL)"
	_, err = db.ExecContext(ctx2, query)
	if err != nil {
		return nil, err
	}

	// // создаю контекст для запроса на создание индекса
	// ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel3()

	// // если в БД нет таблицы, то создаю ее
	// query = "CREATE UNIQUE INDEX idx_url " +
	// 	"ON links (link_url)"
	// _, err = db.ExecContext(ctx3, query)
	// if err != nil {
	// 	return nil, err
	// }

	return &PostgresRepository{
		db: db,
	}, nil
}

// Добавляем ссылку в базу данных
func (r *PostgresRepository) AddLink(link *links.Link) error {

	// проверяем - есть ли уже запись в БД с таким key
	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT link_key FROM links WHERE link_key='%v'", link.Key())
	row := r.db.QueryRowContext(ctx, query)

	var urlKey string

	err := row.Scan(&urlKey)

	// Если запись с таким ключом существует, то возвращаем ошибку
	// Если нашли запись, т.е. err == nil, то возвращаем ошибку
	if err == nil {
		return links.ErrKeyAlreadyExist
	}

	// создаю контекст для запроса
	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()

	query = fmt.Sprintf("INSERT INTO links(link_key, link_url) VALUES ('%v', '%v')", link.Key(), link.URL())
	_, err = r.db.ExecContext(ctx1, query)
	if err != nil {
		// проверяем ошибку на предмет вставки URL который уже есть в БД
		// создаем объект *pgconn.PgError - в нем будет храниться код ошибки из БД
		var pgErr *pgconn.PgError

		// преобразуем ошибку к типу pgconn.PgError
		if errors.As(err, &pgErr) {
			// если ошибка- запис существует, то возвращаем эту ошибку
			if pgErr.Code == pgerrcode.UniqueViolation {
				return links.ErrURLAlreadyExist
			}
		} else {
			return err
		}
	}

	return nil
}

// находим ссылку в БД по ключу
func (r *PostgresRepository) GetLinkByKey(key string) (*links.Link, error) {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT link_url FROM links WHERE link_key='%v'", key)
	row := r.db.QueryRowContext(ctx, query)

	var urlLink string

	err := row.Scan(&urlLink)

	if err != nil {
		return nil, err
	}

	link, err := links.NewLink(key, urlLink)

	if err != nil {
		return nil, err
	}

	return link, nil

}

// находим ссылку в БД по URL
func (r *PostgresRepository) GetLinkByURL(URL string) (*links.Link, error) {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT link_key FROM links WHERE link_url='%v'", URL)
	row := r.db.QueryRowContext(ctx, query)

	var urlKey string

	err := row.Scan(&urlKey)

	if err != nil {
		return nil, err
	}

	link, err := links.NewLink(urlKey, URL)

	if err != nil {
		return nil, err
	}

	return link, nil

}

// узнаем доступность базы данных
func (r *PostgresRepository) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err := r.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
