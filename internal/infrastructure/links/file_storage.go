package links

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// локальный репозиторий
type FileRepository struct {
	// filename *os.File // файл для записи
	fileName string
}

// Создаем локальную базу данных
func NewFileRepository(fileName string) (*FileRepository, error) {

	// currentWorkingDirectory, error := os.Getwd()
	// if error != nil {
	// 	log.Fatal(error)
	// }
	// fmt.Println(currentWorkingDirectory)

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

	db := make(map[string]string)

	err := r.readDBFile(&db)

	if err != nil {
		return err
	}

	// проверяем - есть ли уже записm в БД с таким key
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

	db := make(map[string]string)

	err := r.readDBFile(&db)

	if err != nil {
		return nil, err
	}

	for k, v := range db {

		if k == key {
			link, err := links.NewLink(k, v)
			if err != nil {
				return nil, err
			}
			return link, nil
		}

	}

	return nil, links.ErrLinkNotFound

}

func (r *FileRepository) readDBFile(db *map[string]string) error {

	data, err := os.ReadFile(r.fileName)
	if err != nil {
		return err
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, db)
		if err != nil {
			return err
		}
	}

	return nil

}
