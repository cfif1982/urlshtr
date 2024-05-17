package links

import "github.com/cfif1982/urlshtr.git/internal/domain/links"

// локальный репозиторий
type LocalRepository struct {
	db map[int]LRLink
}

// структура для хранения ссылк в локальном репозитории
type LRLink struct {
	Key string
	URL string
}

// Создаем локальную базу данных
func NewLocalRepository() *LocalRepository {
	return &LocalRepository{
		db: make(map[int]LRLink),
	}
}

// узнаем - есть ли уже запись с данным ключом
func (r *LocalRepository) IsKeyExist(key string) bool {

	// проверяем - есть ли уже записm в БД с таким key
	// Если запись с таким ключом существует, то true
	for _, v := range r.db {
		if v.Key == key {
			return true
		}
	}

	return false
}

// Добавляем ссылку в базу данных
func (r *LocalRepository) AddLink(link *links.Link) error {

	l := LRLink{
		Key: link.Key(),
		URL: link.URL(),
	}
	// добавляем ссылку в БД
	r.db[link.UserID()] = l

	return nil
}

// Добавляем ссылку в базу данных
func (r *LocalRepository) AddLinkBatch(links []*links.Link) error {

	for _, v := range links {
		l := LRLink{
			Key: v.Key(),
			URL: v.URL(),
		}
		// добавляем ссылку в БД
		r.db[v.UserID()] = l
	}

	return nil
}

// находим ссылку в БД по ключу
func (r *LocalRepository) GetLinkByKey(key string) (*links.Link, error) {

	// ищем запись
	for k, v := range r.db {
		if v.Key == key {
			// я так понял, что в DDD не стоит возвращать сслыки на объекты или сами объекты
			// лучше создавать новый объект, копировать в него свойства найденного объекта
			// и уже этот новый объект возвращать
			// я правильно понял?
			link, err := links.NewLink(key, v.URL, k)

			if err != nil {
				return nil, err
			}

			return link, nil
		}
	}

	return nil, links.ErrLinkNotFound
}

// находим ссылку в БД по URL
func (r *LocalRepository) GetLinkByURL(URL string) (*links.Link, error) {

	// ищем запись
	for k, v := range r.db {
		if v.URL == URL {
			link, err := links.NewLink(v.Key, URL, k)

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
