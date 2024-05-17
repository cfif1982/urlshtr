package links

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// файловый репозиторий
type FileRepository struct {
	fileName string
}

// структура для хранения ссылк в файловом репозитории
type FRLink struct {
	Key    string
	URL    string
	UserID int
}

// Создаем файл для хранения БД
func NewFileRepository(fileName string) (*FileRepository, error) {

	// узнаем путь к файлу и назвние самого файла
	absPathToFile, _ := filepath.Abs(fileName)
	absPathToFolder := filepath.Dir(absPathToFile)

	// проверяем существует ли файл
	_, err := os.Stat(absPathToFile)

	// если файл не существует, то создаем его
	if err != nil {
		if os.IsNotExist(err) {

			// проверяем существует ли папка для файла
			_, err = os.Stat(absPathToFolder)

			// если папка не существует, то создаем папку
			if err != nil {
				if os.IsNotExist(err) {
					// создаем папку для файла
					err = os.MkdirAll(absPathToFolder, 0755)
					if err != nil {
						return nil, err
					}
				}
			}

			// создаем файл
			file, err := os.Create(absPathToFile)
			if err != nil {
				return nil, err
			}

			// после использования файла закрываем его
			file.Close()
		}
	}

	return &FileRepository{
		fileName: fileName,
	}, nil
}

// узнаем - есть ли уже запись с данным ключом
func (r *FileRepository) IsKeyExist(key string) bool {

	// создаем БД
	db := make([]FRLink, 0)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return false
	}

	// проверяем - есть ли уже записm в БД с таким key
	// Если запись с таким ключом существует, то true
	for _, v := range db {
		if v.Key == key {
			return true
		}
	}

	return false
}

// Добавляем ссылку в базу данных
func (r *FileRepository) AddLink(link *links.Link) error {

	// создаем БД
	db := make([]FRLink, 0)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return err
	}

	l := FRLink{
		Key:    link.Key(),
		URL:    link.URL(),
		UserID: link.UserID(),
	}

	// добавляем ссылку в БД
	db = append(db, l)

	// маршалим полученный объект в строку для сохранения в файле
	data, err := json.Marshal(&db)
	if err != nil {
		return err
	}

	// записываем данные в файл
	err = os.WriteFile(r.fileName, data, 0666)

	return err
}

// Добавляем ссылку в базу данных
func (r *FileRepository) AddLinkBatch(links []*links.Link) error {

	// создаем БД
	db := make([]FRLink, 0)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return err
	}

	for _, v := range links {
		l := FRLink{
			Key:    v.Key(),
			URL:    v.URL(),
			UserID: v.UserID(),
		}
		// добавляем ссылку в БД
		db = append(db, l)
	}

	// маршалим полученный объект в строку для сохранения в файле
	data, err := json.Marshal(&db)
	if err != nil {
		return err
	}

	// записываем данные в файл
	err = os.WriteFile(r.fileName, data, 0666)

	return err

}

// находим ссылку в БД по ключу
func (r *FileRepository) GetLinkByKey(key string) (*links.Link, error) {

	// создаем БД
	db := make([]FRLink, 0)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return nil, err
	}

	// ищем запись
	for _, v := range db {
		if v.Key == key {
			// я так понял, что в DDD не стоит возвращать сслыки на объекты или сами объекты
			// лучше создавать новый объект, копировать в него свойства найденного объекта
			// и уже этот новый объект возвращать
			// я правильно понял?
			link, err := links.NewLink(key, v.URL, v.UserID)

			if err != nil {
				return nil, err
			}

			return link, nil
		}
	}

	return nil, links.ErrLinkNotFound
}

// находим ссылку в БД по URL
func (r *FileRepository) GetLinkByURL(URL string) (*links.Link, error) {

	// создаем БД
	db := make([]FRLink, 0)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return nil, err
	}

	// ищем запись
	for _, v := range db {
		if v.URL == URL {
			link, err := links.NewLink(v.Key, URL, v.UserID)

			if err != nil {
				return nil, err
			}

			return link, nil
		}
	}

	return nil, links.ErrLinkNotFound
}

// узнаем доступность базы данных. Вернем nil, т.к. эта функция нуна для БД, а здесь всавил ее для совместимости интрефейсов
func (r *FileRepository) Ping() error {

	return nil
}

// читаем файл БД
func (r *FileRepository) readDBFile(db *[]FRLink) error {

	// читаем данные из файла
	data, err := os.ReadFile(r.fileName)
	if err != nil {
		return err
	}

	// если файл пустой, то ничего не делаем
	if len(data) != 0 {
		// анмаршалим данные в БД
		err = json.Unmarshal(data, db)
		if err != nil {
			return err
		}
	}

	return nil
}
