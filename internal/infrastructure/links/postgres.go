package links

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

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
	query := "CREATE TABLE IF NOT EXISTS links(link_key TEXT, link_url TEXT)"
	_, err = db.ExecContext(ctx2, query)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

// Добавляем ссылку в базу данных
func (r *PostgresRepository) AddLink(link *links.Link) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO links(link_key, link_url) VALUES ('%v', '%v')", link.Key(), link.URL())
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
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

	var url_link string

	err := row.Scan(&url_link)

	if err != nil {
		return nil, err
	}

	link, err := links.NewLink(key, url_link)

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
