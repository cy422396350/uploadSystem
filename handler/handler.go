package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"uploadSystem/db"
	"uploadSystem/meta"
	"uploadSystem/util"
)

func UploadHandle(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		abs, _ := filepath.Abs("./static/view/index.html")
		file, err := ioutil.ReadFile(abs)
		if err != nil {
			log.Printf("文件读取出错 %v \n", err)
		}
		io.WriteString(writer, string(file))
	} else if r.Method == "POST" {
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("上传的文件读取出错 %v \n", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		//初始化结构体
		fileMeta := meta.FileMeta{
			FileSha1: "",
			FileName: header.Filename,
			FileSize: 0,
			Location: "./upload/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:15"),
		}
		defer file.Close()
		create, err := os.Create("./upload/" + header.Filename)
		defer create.Close()
		fileMeta.FileSize, err = io.Copy(create, file)
		create.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(create)
		if err != nil {
			panic(err)
		}
		meta.UpdateFilemetas(fileMeta)
		finish := db.OnfileFinish(fileMeta)
		if finish {
			http.Redirect(writer, r, "upload/success", http.StatusFound)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("存储数据库出错"))
		}

	}
}

func UploadSuccess(writer http.ResponseWriter, r *http.Request) {
	io.WriteString(writer, "Upload Success !")
}

//根据hash来找文件
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := r.Form["filehash"][0]
	fileMeta := meta.GetFileMeta(hash)
	marshal, err := json.Marshal(fileMeta)
	if err != nil {
		log.Println("转换出错")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(marshal)
}

// todo:下载
func Download(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(hash)
	open, err := os.Open(fileMeta.Location)
	if err != nil {
		log.Println("地址错误,", fileMeta.Location)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer open.Close()
	all, err := ioutil.ReadAll(open)
	if err != nil {
		log.Println("读取出错,", fileMeta.Location)
		w.WriteHeader(http.StatusBadRequest)
	}

	// TODO:设置头部
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\""+fileMeta.FileName+"\"")
	w.Write(all)
}

// 重命名文件
func RenameFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	hash := r.Form.Get("filehash")
	newName := r.Form.Get("new")
	opType := r.Form.Get("op")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	filemeta := meta.FileMetas[hash]
	os.Rename(filemeta.Location, "./upload/"+newName)
	filemeta.FileName = newName
	filemeta.Location = "./upload/" + newName
	meta.UpdateFilemetas(filemeta)

	data, err := json.Marshal(filemeta)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	hash := r.Form.Get("filehash")
	ok := meta.DeleteFile(hash)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
}
