package main

import (
	"os"
	"testing"
)

// Вспомогательная функция для создания временного файла с заданным содержимым.
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}

	// Записываем данные в файл
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		t.Fatal(err)
	}
	tmpfile.Close()

	return tmpfile.Name()
}

func TestFileRead(t *testing.T) {
	t.Run("FileWithSimple", func(t *testing.T) {
		// Создаем временный файл для тестирования
		// простой случай, выделяем как есть без обработки
		content := "FOOBAR"
		tmpfile := createTempFile(t, content)
		// Удаляем временный файл после тестирования
		defer os.Remove(tmpfile)

		// Проверяем функцию FileRead на этом временном файле
		line, err := FileRead(tmpfile)
		if err != nil {
			t.Errorf("FileRead(%q) returned error: %v", tmpfile, err)
		}

		// ожидаемый возврат: пустая строка
		expected := content
		if line != expected {
			t.Errorf("FileRead(%q) = %q, want %q", tmpfile, line, expected)
		}
	})
	t.Run("FileWithSpace", func(t *testing.T) {
		// Создаем временный файл для тестирования
		// так как будет очистка от пробелов справа то будет пустая строка
		content := " \n"
		tmpfile := createTempFile(t, content)
		// Удаляем временный файл после тестирования
		defer os.Remove(tmpfile)

		// Проверяем функцию FileRead на этом временном файле
		line, err := FileRead(tmpfile)
		if err != nil {
			t.Errorf("FileRead(%q) returned error: %v", tmpfile, err)
		}

		// ожидаемый возврат: пустая строка
		expected := ""
		if line != expected {
			t.Errorf("FileRead(%q) = %q, want %q", tmpfile, line, expected)
		}
	})
	t.Run("FileWithSoupContent", func(t *testing.T) {
		// Создаем временный файл для тестирования
		// Проверка на обработку только одной строки
		// замену всех терминальных нулей
		// и очистку в конце от пробелов и табов,
		// и ожидается многострочник (с сохранением пробелов у первого)
		content := "  AAA  " + "\x00" + "A" + "\x00" + "BBB     \t\t\nCCC\n"
		tmpfile := createTempFile(t, content)
		// Удаляем временный файл после тестирования
		defer os.Remove(tmpfile)

		// Проверяем функцию FileRead на этом временном файле
		line, err := FileRead(tmpfile)
		if err != nil {
			t.Errorf("FileRead(%q) returned error: %v", tmpfile, err)
		}

		// ожидаемый возврат: одна строка с удаленными в конце пробелов и табуляций и замененный %00 на \n
		expected := "  AAA  \nA\nBBB"
		if line != expected {
			t.Errorf("FileRead(%q) = %q, want %q", tmpfile, line, expected)
		}
	})
}
