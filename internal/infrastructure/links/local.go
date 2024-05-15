package links

import "github.com/cfif1982/urlshtr.git/internal/domain/links"

// локальный репозиторий
type LocalRepository struct {
	db map[string]string
}

// Создаем локальную базу данных
func NewLocalRepository() *LocalRepository {
	return &LocalRepository{
		db: make(map[string]string),
	}
}

// узнаем - есть ли уже запись с данным ключом
func (r *LocalRepository) IsKeyExist(key string) bool {

	// проверяем - есть ли уже записm в БД с таким key
	// Если запись с таким ключом существует, то true
	_, ok := r.db[key]

	return ok
}

// Добавляем ссылку в базу данных
func (r *LocalRepository) AddLink(link *links.Link) error {

	// добавляем ссылку в БД
	r.db[link.Key()] = link.URL()

	return nil
}

// Добавляем ссылку в базу данных
func (r *LocalRepository) AddLinkBatch(links []*links.Link) error {

	for _, v := range links {
		// добавляем ссылку в БД
		r.db[v.Key()] = v.URL()
	}

	return nil
}

// находим ссылку в БД по ключу
func (r *LocalRepository) GetLinkByKey(key string) (*links.Link, error) {

	// ищем запись
	l, ok := r.db[key]

	if !ok {
		return nil, links.ErrLinkNotFound
	}

	// я так понял, что в DDD не стоит возвращать сслыки на объекты или сами объекты
	// лучше создавать новый объект, копировать в него свойства найденного объекта
	// и уже этот новый объект возвращать
	// я правильно понял?
	link, err := links.NewLink(key, l)

	if err != nil {
		return nil, err
	}

	return link, nil
}

// находим ссылку в БД по URL
func (r *LocalRepository) GetLinkByURL(URL string) (*links.Link, error) {

	// ищем запись
	for k, v := range r.db {

		if v == URL {
			link, err := links.NewLink(k, v)

			if err != nil {
				return nil, err
			}

			return link, nil
		}
	}
	return nil, links.ErrLinkNotFound
}

// узнаем доступность базы данных. Локальный репозиторий всегда доступен
func (r *LocalRepository) Ping() error {

	return nil
}
