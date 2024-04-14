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

// Добавляем ссылку в базу данных
func (r *LocalRepository) AddLink(link links.Link) error {

	// проверяем - есть ли уже записm в БД с таким key
	_, ok := r.db[link.Key()]
	if ok {
		return links.ErrKeyAlreadyExist
	}

	// если такого key в БД нет, то добавляем ссылку в БД
	r.db[link.Key()] = link.URL()

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

// возвращаем ссылку на БД для тестов
// для проведения теста нужно получить результаты работы функции AddLink.
// AddLink добавляет запись в БД. Чтобы проверить запись и получить сгенерированный key для этой запсиси,
// добавляю функцию GetDBForTest
// Так можно делать? Я имею в виду - добавлять в код функци, которые нужны только для тестирования
func (r *LocalRepository) GetDBForTest() *map[string]string {
	return &r.db
}
