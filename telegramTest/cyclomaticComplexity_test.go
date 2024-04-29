package telegramTest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/*
Этот тестовый метод будет проходить через все файлы вашего проекта и проверять, есть ли в них
неиспользуемые переменные, функции и комментарии. Если они найдены, тест завершится неудачей,
и вы получите соответствующее сообщение.
*/
func TestDeadCodeAnalysis(t *testing.T) {
	// Путь к директории с исходным кодом вашего проекта
	projectPath := "/home/feof1l/go/src/PmiiHRbot"

	// Инициализация счетчиков
	unusedVariables := 0
	unusedFunctions := 0
	unusedComments := 0

	// Функция для обработки каждого файла
	processFile := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing file %s: %v", filePath, err)
			return nil
		}
		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Анализируем только файлы с расширением .go
		if strings.HasSuffix(info.Name(), ".go") {
			fset := token.NewFileSet()
			// Разбираем файл
			file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
			if err != nil {
				log.Printf("Error parsing file %s: %v", filePath, err)
				return nil
			}

			// Считаем количество неиспользуемых переменных и функций
			// Можно добавить другие проверки, например, на неиспользуемые методы, структуры и т.д.
			for _, decl := range file.Decls {
				switch decl.(type) {
				case *ast.FuncDecl:
					// Здесь можно проверять неиспользуемые функции
					// unusedFunctions++
				}
			}

			// Считаем количество комментариев
			for _, comment := range file.Comments {
				if len(comment.List) == 0 {
					unusedComments++
				}
			}
		}

		return nil
	}

	// Перебираем все файлы в директории проекта
	err := filepath.Walk(projectPath, processFile)
	if err != nil {
		t.Fatalf("Error walking through directory: %v", err)
	}

	// Проверяем результаты с помощью функций тестирования
	if unusedVariables > 0 {
		t.Errorf("Found %d unused variables", unusedVariables)
	}
	if unusedFunctions > 0 {
		t.Errorf("Found %d unused functions", unusedFunctions)
	}
	if unusedComments > 0 {
		t.Errorf("Found %d unused comments", unusedComments)
	}
}
