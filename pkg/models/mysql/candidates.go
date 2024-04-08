package mysql

import (
	"database/sql"
	"errors"

	"github.com/Feof1l/TelegramHrBot/pkg/models"
)

// Определяем тип который обертывает пул подключения sql.DB
type CandidatModel struct {
	DB *sql.DB
}

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
func (m *CandidatModel) GetId(candidateName, telegramUsername string) (int, error) {
	var id int
	err := m.DB.QueryRow("SELECT id_possible_candidate FROM Possible_candidate WHERE (Candidate_name = ? AND Telegram_username = ?) ", candidateName, telegramUsername).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrNoSuchRowInColumn := errors.New("нет такой строки в столбце")
			return 0, ErrNoSuchRowInColumn
		}
		return 0, err
	}
	return id, nil
}
func (m *CandidatModel) CallCompareEducation(position_id, possible_candidat_id int) error {
	_, err := m.DB.Exec("CALL Compare_Education(?, ?)", position_id, possible_candidat_id)
	if err != nil {
		return err
	}
	return nil
}

// Метод для добавления дынных в существующую запись  в базе дынных.
func (m *CandidatModel) Update(education string, candidateName string) error {
	// Подготовка SQL-запроса для вставки данных в таблицу
	query := `UPDATE Possible_candidate SET Education = ? WHERE Candidate_name = ?`

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Выполнение запроса с передачей параметров
	_, err = stmt.Exec(education, candidateName)
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
