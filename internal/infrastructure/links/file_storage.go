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

// Добавляем ссылку в базу данных
func (r *FileRepository) AddLink(link *links.Link) error {

	// создаем БД
	db := make(map[string]string)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return err
	}

	// проверяем - есть ли уже запись в БД с таким key
	_, ok := db[link.Key()]
	if ok {
		return links.ErrKeyAlreadyExist
	}

	// если такого key в БД нет, то добавляем ссылку в БД
	db[link.Key()] = link.URL()

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
	db := make(map[string]string)

	// загружаем данные из файла
	err := r.readDBFile(&db)
	if err != nil {
		return nil, err
	}

	// пробегаемся по БД и ищем нужную ссылку
	for k, v := range db {
		// если ссылка найдена, то возвращаем ее
		if k == key {
			link, err := links.NewLink(k, v)
			if err != nil {
				return nil, err
			}
			return link, nil
		}
	}

	// если ссылка не найдена, то возвращаем ошибку
	return nil, links.ErrLinkNotFound

}

// узнаем доступность базы данных. Вернем nil, т.к. эта функция нуна для БД, а здесь всавил ее для совместимости интрефейсов
func (r *FileRepository) Ping() error {

	return nil

}

// читаем файл БД
func (r *FileRepository) readDBFile(db *map[string]string) error {

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
