package meta

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var FileMetas map[string]FileMeta

func init() {
	FileMetas = make(map[string]FileMeta)
}

func UpdateFilemetas(fileMeta FileMeta) {
	FileMetas[fileMeta.FileSha1] = fileMeta
}

func GetFileMeta(hash string) FileMeta {
	return FileMetas[hash]
}

func DeleteFile(hash string) bool {
	_, ok := FileMetas[hash]
	delete(FileMetas, hash)
	return ok
}
