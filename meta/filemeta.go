package meta

import mydb "Distributed-fileserver/db"

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}

//新增/更新文件元信息
func UploadFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1] = fmeta
}

// 新增与更新元信息到数据库
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//通过sha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}

func GetFileMetaDB(filesha1 string) (FileMeta, error){
	tfile, err := mydb.GetFileMeta(filesha1)

	if err != nil{
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1:tfile.FileHash,
		FileName:tfile.FileName.String,
		FileSize:tfile.FileSize.Int64,
		Location:tfile.FileAddr.String,
	}
	return fmeta, nil
}

//删除文件元信息
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}