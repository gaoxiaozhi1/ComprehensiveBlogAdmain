package images_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gvb_server/global"
	"gvb_server/models/res"
	"os"
	"path"
)

// 图片上传的响应
type FileUploadResponse struct {
	FileName  string `json:"file_name"`  // 文件名
	IsSuccess bool   `json:"is_success"` // 是否上传成功
	Msg       string `json:"msg"`        // 消息
}

// 上传图片，返回图片的URL
func (ImagesApi) ImageUploadView(c *gin.Context) {
	// func (c *Context) MultipartForm() (*multipart.Form, error)
	//
	// type Form struct {
	//	Value map[string][]string
	//	File  map[string][]*FileHeader
	// }
	form, err := c.MultipartForm()
	if err != nil {
		res.FailWithMessage(err.Error(), c)
		return
	}

	fileList, ok := form.File["images"]
	if !ok {
		res.FailWithMessage("不存在的文件", c)
		return
	}

	// 判断路径是否存在
	// 截断（逐级判断）
	basePath := global.Config.Upload.Path
	// 在Go语言中，os.ReadDir(basePath)函数用于读取指定路径basePath下的所有目录条目，并将它们以文件名排序后返回。
	// 这个函数返回一个FileInfo结构的数组，每个元素代表一个目录条目。
	// 这个函数可以用来列出一个目录下的所有文件和子目录。
	_, err = os.ReadDir(basePath)
	if err != nil {
		// 在Go语言中，os.MkdirAll(basePath, os.ModePerm)函数用于创建指定路径basePath下的所有目录，包括任何必要的父目录。
		// 这个函数返回nil，或者在出错时返回一个错误。
		// 第二个参数os.ModePerm是一个权限位，用于设置所有被该函数创建的目录的读/写/执行权限。
		// 在Unix-like系统中，os.ModePerm等价于07772，意味着用户有权列出、修改和搜索目录中的文件。
		// 如果目录已经存在，os.MkdirAll()函数不会做任何事情，而是返回nil。
		err = os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			global.Log.Error(err)
		}
	}

	// 不存在就创建
	var resList []FileUploadResponse

	for _, file := range fileList {
		// 存数据库。。。。

		// 图片存储
		filePath := path.Join(basePath, file.Filename) // 存储路径

		// 判断大小
		// 加float64是因为要进行浮点数除法
		size := float64(file.Size) / float64(1024*1024)
		if size > float64(global.Config.Upload.Size) {
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       fmt.Sprintf("图片大小超过设定大小，当前大小为：%.2fMB，设定大小为：%dMB", size, global.Config.Upload.Size),
			})
			continue
		}

		// 上传
		err := c.SaveUploadedFile(file, filePath)

		// 上传失败
		if err != nil {
			global.Log.Error(err)
			resList = append(resList, FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       err.Error(),
			})
			continue
		}

		// 上传成功
		resList = append(resList, FileUploadResponse{
			FileName:  filePath,
			IsSuccess: true,
			Msg:       "上传成功",
		})

	}

	res.OKWithData(resList, c)
}