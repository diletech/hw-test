package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

func FileRead(file string) (line string, err error) {
	// Открываем файл для чтения
	r, err := os.Open(file)
	if err != nil {
		return line, err
	}
	defer r.Close()

	// Создаем новый сканнер для файла
	scanner := bufio.NewScanner(r)

	// Считываем только первую строку из файла
	if scanner.Scan() {
		// получаем текст первой строки
		line = scanner.Text()
		// обрезаем пробелы и табуляции в конце строки
		line = strings.TrimRight(line, " \t")
		// заменяем терминальные нули на перевод строки
		line = string(bytes.ReplaceAll([]byte(line), []byte{0x00}, []byte("\n")))
	}

	// Проверяем наличие ошибок при сканировании файла
	if err := scanner.Err(); err != nil {
		return line, err
	}
	return line, nil
}
