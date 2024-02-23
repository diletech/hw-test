package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	// получаем список файлов
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envData := make(Environment)
	// экспортируем всё окружение в карту данных
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		value := pair[1]
		envData[key] = EnvValue{Value: value}
	}

	for _, file := range dirEntry {
		if file.IsDir() {
			continue // пропустить директории
		}

		// проверяем, что это файл
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		// Пропускаем файл если он не регулярный.
		if !fileInfo.Mode().IsRegular() {
			continue
		}

		// Начинаем обработку регулярного файла.
		fileName := file.Name() // вытаскиваем имя строковое файла
		// проверяем, что имя файла содержит символ "="
		if strings.Contains(fileName, "=") {
			return nil, fmt.Errorf("invalid character '=' in file name: %s", fileName)
		}
		// проверяем на размер в 0 байт
		if fileInfo.Size() == 0 {
			// проверяем присутствует ли этот env в карте со всем окружением
			if _, ok := envData[fileName]; ok {
				// если обнаружен, то выставляем соответсвующий флаг
				envData[fileName] = EnvValue{NeedRemove: true}
			}
			continue // и переходим к следующему
		}
		// отправляем на извлечение параметра через разбор первой строчки файла
		data, err := FileRead(filepath.Join(dir, fileName))
		if err != nil {
			return nil, err
		}
		envData[fileName] = EnvValue{Value: data}
	}

	return envData, nil
}
