package links

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cfif1982/urlshtr.git/pkg/log"
	"github.com/pressly/goose/v3"

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
func NewPostgresRepository(ctx context.Context, databaseDSN string, logger *log.Logger) (*PostgresRepository, error) {

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	// создаю контекст для пинга
	ctx2, cancel2 := context.WithTimeout(ctx, 1*time.Second)
	defer cancel2()

	// пингую БД. Если не отвечает, то возвращаю ошибку
	if err = db.PingContext(ctx2); err != nil {
		return nil, err
	}

	// начинаю миграцию
	logger.Info("Start migrating database")

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal(err.Error())
	}

	// такое указание пути к папке с миграциями на локальном компе работает, но в тестах не проходит
	// пришлось перенести файлы в дпапку cmd/shortener/migrations и указать путь по другому
	// err = goose.Up(db, "../../internal/infrastructure/migrations")
	err = goose.Up(db, "migrations")
	if err != nil {
		logger.Info(err.Error())
	}

	logger.Info("migrating database finished")

	return &PostgresRepository{
		db: db,
	}, nil
}

// узнаем - есть ли уже запись с данным ключом
func (r *PostgresRepository) IsKeyExist(key string) bool {

	// проверяем - есть ли уже запись в БД с таким key
	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT count(*) FROM links WHERE link_key='%v'", key)
	row := r.db.QueryRowContext(ctx, query)

	var urlKey string

	_ = row.Scan(&urlKey)

	// Если запись с таким ключом существует, то true
	return urlKey != "0"
}

// Добавляем ссылку в базу данных
func (r *PostgresRepository) AddLink(link *links.Link) error {

	// создаю контекст для запроса
	ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel1()

	query := fmt.Sprintf("INSERT INTO links(link_key, link_url) VALUES ('%v', '%v')", link.Key(), link.URL())
	_, err := r.db.ExecContext(ctx1, query)
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

	// в эту переменную будет сканиться результат запроса
	var urlLink string

	err := row.Scan(&urlLink)

	if err != nil {
		return nil, err
	}

	// создаем объект ссылку и возвращаем ее
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

	// в эту переменную будет сканиться результат запроса

	var urlKey string

	err := row.Scan(&urlKey)

	if err != nil {
		return nil, err
	}

	// создаем объект ссылку и возвращаем ее
	link, err := links.NewLink(urlKey, URL)

	if err != nil {
		return nil, err
	}

	return link, nil
}

// Добавляем ссылку в базу данных
func (r *PostgresRepository) AddLinkBatch(links []*links.Link) error {

	// создаю транзакцию для вставки всех ссылок одним запросом
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// подготавливаю запрос для транзакции
	stmt, err := r.db.PrepareContext(ctx,
		"INSERT INTO links(link_key, link_url) "+
			"VALUES ($1, $2) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// добавляю запросы в транзакцю
	for _, v := range links {
		_, err := stmt.ExecContext(ctx, v.Key(), v.URL())
		if err != nil {
			return err
		}
	}

	// зпускаю транзакцию
	return tx.Commit()
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
