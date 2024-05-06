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

	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", databaseDSN, `postgres`, `123`, `videos`)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, err
	}

	// TODO: Нужно разобраться. Если оставляю этот код, то соединене закрывается и нет доступа к БД. Где это нужно делать?
	//defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

// Добавляем ссылку в базу данных
func (r *PostgresRepository) AddLink(link *links.Link) error {

	return nil
}

// находим ссылку в БД по ключу
func (r *PostgresRepository) GetLinkByKey(key string) (*links.Link, error) {

	// если ссылка не найдена, то возвращаем ошибку
	return nil, nil

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
