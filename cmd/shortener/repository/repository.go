package repository

// Локальная БД
type LocalDatabase struct {
	ReceivedURL map[string]string
}

// // инициализируем БД
// func (database LocalDatabase) InitDatabase() {
// 	database.ReceivedURL = make(map[string]string)
// }

// сохраняем переданные польователем данные
func (database LocalDatabase) SaveURL(key string, value string) {
	database.ReceivedURL[key] = value
}

// находим данные в БД по ключу
func (database LocalDatabase) GetURL(key string) string {
	return database.ReceivedURL[key]
}
