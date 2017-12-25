package orm

import (
	"errors"
)

// GetSequenceIDByName : Get sequence id by name
func GetSequenceIDByName(name string) (id int64, err error) {
	session, err := DB.Begin()
	if err != nil {
		return -1, err
	}

	sql := "select current_value, increment_value from lic_t_sequence where sequence_name = ?"
	rows, err := session.Query(sql, name)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	if false == rows.Next() {
		rows.Close()
		session.Rollback()
		return -1, errors.New("sequence name is not existing")
	}

	var currentValue, incrementValue int64
	if err = rows.Scan(&currentValue, &incrementValue); err != nil {
		rows.Close()
		session.Rollback()
		return -1, err
	}

	rows.Close()

	sql = "update lic_t_sequence set current_value = current_value + increment_value where sequence_name = ?"
	_, err = session.Exec(sql, name)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	session.Commit()
	return currentValue, nil
}
