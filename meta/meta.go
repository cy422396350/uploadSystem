package meta

import (
	"log"
	myDb "uploadSystem/db/mysql"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
	CreateAt string
}

type Operation interface {
	UpdateFilemetas() bool
	GetFileMeta(string) *FileMeta
	DeleteFile(hash string) bool
}

var FileMetas map[string]FileMeta

const table = "uploadfile"

func init() {
	FileMetas = make(map[string]FileMeta)
}

func (f *FileMeta) UpdateFilemetas() bool {
	if f.CreateAt != "" {
		// TODO 更新
		return true
	}
	db := myDb.GetDb()
	prepare, err := db.Prepare("INSERT INTO " + table + " (`file_sha1`,`filename`,`filesize`,`fileaddr`,`status`) values (?,?,?,?,?)")
	if err != nil {
		log.Println(err)
	}
	defer prepare.Close()
	exec, err := prepare.Exec(f.FileSha1, f.FileName, f.FileSize, f.Location, 1)
	if err != nil {
		log.Println(err)
	}
	if eff, err := exec.RowsAffected(); nil == err && eff > 0 {
		return true
	}
	return false
}

func (f *FileMeta) GetFileMeta(hash string) *FileMeta {
	filemeta := FileMeta{
		FileSha1: hash,
		FileName: "",
		FileSize: 0,
		Location: "",
		UploadAt: "",
	}
	db := myDb.GetDb()
	prepare, err := db.Prepare("select filename,filesize,fileaddr,created_at from " + table + " where file_sha1 = ? and  status = 1 limit 1")
	if err != nil {
		log.Println(err)
	}
	defer prepare.Close()
	res := prepare.QueryRow(filemeta.FileSha1)
	err = res.Scan(&filemeta.FileName, &filemeta.FileSize, &filemeta.Location, &filemeta.CreateAt)
	if err != nil {
		log.Println(err)
	}
	return &filemeta
}

func (f *FileMeta) DeleteFile(hash string) bool {
	return true
}
