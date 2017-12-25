package sql

import (
	"SilverBusinessServer/lib"
	"database/sql"
	"errors"
)

//QuerySearchPersonResultNumber :
func QuerySearchPersonResultNumber(session *sql.Tx) (number int, err error) {
	var rows *sql.Rows

	sqlStr := `select count(id) from imp_t_favorite`
	rows, err = session.Query(sqlStr)

	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if false == rows.Next() {
		return -1, nil
	}

	rows.Scan(&number)

	return number, nil
}

// QuerySearchPersonResultByPage :
func QuerySearchPersonResultByPage(session *sql.Tx, offset, count int64) (results []lib.SearchPersonResult, err error) {

	var rows *sql.Rows

	sqlStr := `select t1.id, t1.image_url, t1.image_time
						from imp_t_favorite t1
						where  t1.id <=
						(
						select t2.id from imp_t_favorite t2
						order by t2.id desc limit ?,1
						)
						order by t1.id desc limit ?`
	rows, err = session.Query(sqlStr, offset, count)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results = []lib.SearchPersonResult{}
	for rows.Next() {
		result := lib.SearchPersonResult{}
		if err = rows.Scan(&result.DeviceUUID, &result.ImageURL, &result.ImageTime); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// SaveSearchPersonRule :
func SaveSearchPersonRule(session *sql.Tx, SearchRuleID int64, SearchRule string) (err error) {
	sql := `insert into imp_t_searchrule values(?, ?)`

	_, err = session.Exec(sql, SearchRuleID, SearchRule)
	if err != nil {
		return err
	}

	return nil
}

// QuerySearchPersonRule :
func QuerySearchPersonRule(session *sql.Tx, SearchRuleID int64) (SearchRule string, err error) {
	var rows *sql.Rows

	sqlStr := `select rule_content from imp_t_searchrule where rule_id = ?`

	rows, err = session.Query(sqlStr, SearchRuleID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&SearchRule); err != nil {
			return "", err
		}
		return SearchRule, nil
	}
	return "", errors.New("Not found response search people rule.")
}
