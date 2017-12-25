package dao

//import (
//	"SilverBusinessServer/dao/sql"
//	"SilverBusinessServer/env"
//	"SilverBusinessServer/lib"
//)

////QueryFavoritesNumber :
//func QueryFavoritesNumber(deviceIDListStr string, beginTime, endTime int64) (number int, err error) {
//	session, err := sql.DB.Begin()
//	if err != nil {
//		return -1, err
//	}

//	number, err = sql.QueryFavoritesNumber(session, deviceIDListStr, beginTime, endTime)
//	if err != nil {
//		session.Rollback()
//		return -1, err
//	}

//	session.Commit()
//	return number, nil
//}

//// QueryFavoritesByPage :
//func QueryFavoritesByPage(deviceIDListStr string, beginTime, endTime int64, offset, count int) (favoriteList []lib.Favorite, err error) {
//	session, err := sql.DB.Begin()
//	if err != nil {
//		return nil, err
//	}

//	favoriteList, err = sql.QueryFavoritesByPage(session, deviceIDListStr, beginTime, endTime, offset, count)
//	if err != nil {
//		session.Rollback()
//		return nil, err
//	}

//	session.Commit()
//	return favoriteList, nil
//}

//// QueryFavoriteDevice :
//func QueryFavoriteDevice() (favoriteGroup []lib.FavoriteGroup, err error) {
//	session, err := sql.DB.Begin()
//	if err != nil {
//		return nil, err
//	}

//	favoriteGroup, err = sql.QueryFavoriteDevice(session)
//	if err != nil {
//		session.Rollback()
//		return nil, err
//	}

//	session.Commit()
//	return favoriteGroup, nil
//}

//// DeleteFavorite :
//func DeleteFavorite(fid int64) (err error) {
//	session, err := sql.DB.Begin()
//	if err != nil {
//		return err
//	}

//	err = sql.DeleteFavorite(session, fid)
//	if err != nil {
//		session.Rollback()
//		return err
//	}

//	session.Commit()
//	return nil
//}

//// AddFavorite :
//func AddFavorite(favorite lib.Favorite) (favoriteID int64, err error) {
//	session, err := sql.DB.Begin()
//	if err != nil {
//		return -1, err
//	}

//	//Get favorite id
//	next, err := sql.GetSequenceIDByName(env.SequenceNameFavoriteID)
//	if err != nil {
//		session.Rollback()
//		return -1, err
//	}
//	favoriteID = next
//	favorite.FavoriteID = favoriteID
//	err = sql.AddFavorite(session, favorite)
//	if err != nil {
//		session.Rollback()
//		return -1, err
//	}

//	session.Commit()
//	return favoriteID, nil
//}
