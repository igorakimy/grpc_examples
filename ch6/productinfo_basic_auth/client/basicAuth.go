package main

import (
	"context"
	"encoding/base64"
)

// Определяем структуру для хранения полей, которые нужно внедрять
// в удаленные вызовы (в нашем случае это такие учетные данные,
// как имя пользователя и пароль)
type basicAuth struct {
	username string
	password string
}

// GetRequestMetadata преобразует данные аутентификации в метаданные запроса.
// В нашем случае ключ равен Authorization, а значение состоит из слова Basic,
// за которым следует строка <имя_пользователя>:<пароль> в кодировке base64
func (a basicAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	auth := a.username + ":" + a.password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

// RequireTransportSecurity выступает в роли безопасного канала
// для передачи аутентификационных данных
func (a basicAuth) RequireTransportSecurity() bool {
	return true
}
