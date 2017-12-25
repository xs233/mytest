package dao

import (
	"SilverBusinessServer/dao/sql"
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
)

//QuerySearchPersonResultNumber :
func QuerySearchPersonResultNumber() (number int, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return -1, err
	}

	number, err = sql.QuerySearchPersonResultNumber(session)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	session.Commit()
	return number, nil
}

// QuerySearchPersonResultByPage :
func QuerySearchPersonResultByPage(offset, count int64) (resultList []lib.SearchPersonResult, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return nil, err
	}

	resultList, err = sql.QuerySearchPersonResultByPage(session, offset, count)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return resultList, nil
}

// SaveSearchPersonRule :
func SaveSearchPersonRule(SearchRule string) (SearchRuleID int64, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return -1, err
	}

	SearchRuleID, err = sql.GetSequenceIDByName(env.SequenceNameSearchRuleID)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	err = sql.SaveSearchPersonRule(session, SearchRuleID, SearchRule)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	session.Commit()
	return SearchRuleID, nil
}

// QuerySearchPersonRule :
func QuerySearchPersonRule(SearchRuleID int64) (SearchRule string, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return "", err
	}

	SearchRule, err = sql.QuerySearchPersonRule(session, SearchRuleID)
	if err != nil {
		session.Rollback()
		return "", err
	}

	session.Commit()
	return SearchRule, nil
}
