package mysql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Feof1l/TelegramHrBot/pkg/models"
)

func TestCandidatModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	candidateName := "John Doe"
	telegramUsername := "johndoe"
	idPosition := 1

	mock.ExpectPrepare("INSERT INTO Possible_candidate").
		ExpectExec().
		WithArgs(candidateName, telegramUsername, idPosition, false, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = model.Insert(candidateName, telegramUsername, idPosition)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_InsertFeadBack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	feedback := "Good job!"
	id := 1

	mock.ExpectPrepare("UPDATE Possible_candidate").
		ExpectExec().
		WithArgs(feedback, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = model.InsertFeadBack(feedback, id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_GetFlag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	flagName := "fail_flag"
	candidateID := 1

	mock.ExpectQuery("SELECT fail_flag FROM Possible_candidate").
		WithArgs(candidateID).
		WillReturnRows(sqlmock.NewRows([]string{"fail_flag"}).AddRow(false))

	flag, err := model.GetFlag(flagName, candidateID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if flag != false {
		t.Errorf("expected flag to be false, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_GetId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	candidateName := "John Doe"
	telegramUsername := "johndoe"

	mock.ExpectQuery("SELECT id_possible_candidate FROM Possible_candidate").
		WithArgs(candidateName, telegramUsername).
		WillReturnRows(sqlmock.NewRows([]string{"id_possible_candidate"}).AddRow(1))

	id, err := model.GetId(candidateName, telegramUsername)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if id != 1 {
		t.Errorf("expected id to be 1, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_CallStoredProcedure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	procedureName := "TestProcedure"
	positionID := 1
	candidateID := 1

	mock.ExpectPrepare("CALL TestProcedure").
		ExpectExec().
		WithArgs(positionID, candidateID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = model.CallStoredProcedure(procedureName, positionID, candidateID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_UpdateStringData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	field := "Candidate_name"
	data := "John Doe"
	id := 1

	mock.ExpectPrepare("UPDATE Possible_candidate").
		ExpectExec().
		WithArgs(data, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = model.UpdateStringData(field, data, id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_UpdateBoolData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	field := "ready_flag"
	data := true
	id := 1

	mock.ExpectPrepare("UPDATE Possible_candidate").
		ExpectExec().
		WithArgs(data, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = model.UpdateBoolData(field, data, id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_UpdateIntData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	field := "id_pos"
	data := 2
	id := 1

	mock.ExpectPrepare("UPDATE Possible_candidate").
		ExpectExec().
		WithArgs(data, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = model.UpdateIntData(field, data, id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCandidatModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := &CandidatModel{DB: db}

	expectedPosition := &models.Position{Id_position: 1, Position_name: "Test Position", Profil: "Test Profile"}

	mock.ExpectQuery("SELECT id_position,Position_name,Profil FROM position").
		WillReturnRows(sqlmock.NewRows([]string{"id_position", "Position_name", "Profil"}).
			AddRow(expectedPosition.Id_position, expectedPosition.Position_name, expectedPosition.Profil))

	position, err := model.Get(1)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if position.Id_position != expectedPosition.Id_position {
		t.Errorf("expected id to be %d, got %d", expectedPosition.Id_position, position.Id_position)
	}
	if position.Position_name != expectedPosition.Position_name {
		t.Errorf("expected position name to be %s, got %s", expectedPosition.Position_name, position.Position_name)
	}
	if position.Profil != expectedPosition.Profil {
		t.Errorf("expected profile to be %s, got %s", expectedPosition.Profil, position.Profil)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
