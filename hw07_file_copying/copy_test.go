package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	// Создаем временные файлы для тестирования
	fromFile, err := os.CreateTemp("", "from-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(fromFile.Name())

	toFile, err := os.CreateTemp("", "to-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(toFile.Name())

	// Записываем данные в исходный файл
	data := []byte("Hello, world!")
	if _, err := fromFile.Write(data); err != nil {
		t.Fatalf("failed to write to fromFile: %v", err)
	}

	tests := []struct {
		name         string
		fromPath     string
		toPath       string
		offset       int64
		limit        int64
		expectedData []byte // Ожидаемые данные в скопированном файле
		wantErr      bool
		wantError    error // Ожидаемая ошибка
	}{
		{
			name:         "Copy entire file",
			fromPath:     fromFile.Name(),
			toPath:       toFile.Name(),
			offset:       0,
			limit:        0,
			expectedData: data, // Ожидаем копию всего файла
			wantErr:      false,
			wantError:    nil,
		},
		{
			name:         "Copy part of file",
			fromPath:     fromFile.Name(),
			toPath:       toFile.Name(),
			offset:       0,
			limit:        5,
			expectedData: []byte("Hello"), // Ожидаем копию первых 5 байт файла
			wantErr:      false,
			wantError:    nil,
		},
		{
			name:         "Copy part of file with offset and over limit",
			fromPath:     fromFile.Name(),
			toPath:       toFile.Name(),
			offset:       7,
			limit:        55,
			expectedData: []byte("world!"), // Ожидаем копию с 7-ого байта и до конца
			wantErr:      false,
			wantError:    nil,
		},
		{
			name:         "Offset exceeds file size",
			fromPath:     fromFile.Name(),
			toPath:       toFile.Name(),
			offset:       100,
			limit:        0,
			expectedData: nil,
			wantErr:      true,
			wantError:    ErrOffsetExceedsFileSize,
		},
		{
			name:         "Non regular file",
			fromPath:     "/dev/random",
			toPath:       toFile.Name(),
			offset:       0,
			limit:        0,
			expectedData: nil,
			wantErr:      true,
			wantError:    ErrUnsupportedFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.fromPath, tt.toPath, tt.offset, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !errors.Is(err, tt.wantError) {
				t.Errorf("Copy() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// Если не ожидается ошибка, проверяем содержимое скопированного файла
			if !tt.wantErr {
				// Открываем скопированный файл для чтения
				copiedFile, err := os.Open(tt.toPath)
				if err != nil {
					t.Fatalf("failed to open copied file: %v", err)
				}
				defer copiedFile.Close()

				// Считываем данные из скопированного файла
				copiedData, err := io.ReadAll(copiedFile)
				if err != nil {
					t.Fatalf("failed to read copied file: %v", err)
				}

				// Сравниваем данные с ожидаемыми данными
				if !bytes.Equal(copiedData, tt.expectedData) {
					t.Errorf("copied data does not match expected data")
				}
			}
		})
	}
}
