package meta

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

//通过sha1获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta{
	return fileMetas[fileSha1]
}

//删除文件元信息
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}