package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Feof1l/TelegramHrBot/pkg/models"
)

// Определяем тип который обертывает пул подключения sql.DB
type CandidatModel struct {
	DB *sql.DB
}

var ErrNoSuchRowInColumn = errors.New("Строка в таблице не найдена")

// Метод для создания записи  в базе дынных.
func (m *CandidatModel) Insert(candidateName, telegramUsername string, Id_position int) error {
	// Подготовка SQL-запроса для вставки данных в таблицу

	query := `INSERT INTO Possible_candidate (Candidate_name,Telegram_username,id_pos,fail_flag) VALUES (?,?,?,?)`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса с передачей параметров
	defaultFailFlag := false
	_, err = stmt.Exec(candidateName, telegramUsername, Id_position, defaultFailFlag)
	if err != nil {
		return err
	}
	return nil

}

// Запись фидбека после общения с ботом
func (m *CandidatModel) InsertFeadBack(feadback string, id_possible_candidate int) error {
	// Подготовка SQL-запроса для вставки данных в таблицу

	query := `UPDATE Possible_candidate SET Feadback = ? WHERE id_possible_candidate = ?`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса с передачей параметров
	_, err = stmt.Exec(feadback, id_possible_candidate)
	if err != nil {
		return err
	}
	return nil

}

// получение флага ошибки чтобы отснивать кандидатов
func (m *CandidatModel) GetFailFlag(candidate_id int) (bool, error) {
	var failFlag bool
	if err := m.DB.QueryRow("SELECT fail_flag FROM Possible_candidate WHERE id_possible_candidate = ?", candidate_id).Scan(&failFlag); err != nil {
		if err == sql.ErrNoRows {
			return false, ErrNoSuchRowInColumn
		} else {
			return false, err
		}
	}
	return failFlag, nil
}

// метод получения id по имени и нику
func (m *CandidatModel) GetId(candidateName, telegramUsername string) (int, error) {
	var id int
	err := m.DB.QueryRow("SELECT id_possible_candidate FROM Possible_candidate WHERE (Candidate_name = ? AND Telegram_username = ?) ", candidateName, telegramUsername).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrNoSuchRowInColumn
		}
		return 0, err
	}
	return id, nil
}

// обёртка для вызова хранимой процедуры из БД
func (m *CandidatModel) CallCompareEducation(position_id, possible_candidat_id int) error {
	_, err := m.DB.Exec("CALL Compare_Education(?, ?)", position_id, possible_candidat_id)
	if err != nil {
		return err
	}
	return nil
}

// Метод для добавления дынных в существующую запись  в базе дынных.
func (m *CandidatModel) Update(field, data string, id int) error {
	// Подготовка SQL-запроса для вставки данных в таблицу
	response := fmt.Sprintf(`UPDATE Possible_candidate SET %s = ? WHERE id_possible_candidate = ?`, field)
	query := response

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса с передачей параметров
	_, err = stmt.Exec(data, id)
	if err != nil {
		return err
	}
	return nil

}

/*func (m *LinkModel) Insert(education string) (int, error) {
	//stmt := `INSERT INTO Possible_candidate (Candidate_name,Telegram_username,contact_number,Citizenship,Education,Work_experience,Hours,Work_format,Expected_salary,Ready_to_relocate,Date_of_dialog,Feadback,is_blocked_flag,ready_flag,fail_flag)
	//VALUES (?,?,?,?,?,?,?,?,?,?)`
	stmt := `INSERT INTO Possible_candidate (education)
	VALUES (?)`
	result, err := m.DB.Exec(stmt, education)
	if err != nil {
		return 0, err
	}
	// Используем метод LastInsertId(), чтобы получить последний ID созданной записи из таблицу links.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}*/

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *CandidatModel) Get(id int) (*models.Position, error) {
	stmt := `SELECT id_position,Position_name,Profil FROM position`

	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(stmt, id)

	// Инициализируем указатель на новую структуру Link
	s := &models.Position{}
	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в
	// соответствующее поле в структуре Link.
	err := row.Scan(&s.Id_position, &s.Position_name, &s.Profil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
