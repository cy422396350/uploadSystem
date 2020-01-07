package db

import (
	"log"
	myDb "uploadSystem/db/mysql"
	"uploadSystem/meta"
)

func OnfileFinish(fileMeta meta.FileMeta) bool {
	db := myDb.GetDb()
	prepare, err := db.Prepare("INSERT INTO `uploadfile` (`file_sha1`,`filename`,`filesize`,`fileaddr`,`status`) values (?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	exec, err := prepare.Exec(fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize, fileMeta.Location, 1)
	if err != nil {
		log.Println(err)
	}
	if eff, err := exec.RowsAffected(); nil == err && eff > 0 {
		return true
	}
	return false
}
