package utils

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"path"
	"io"
	"strings"
	"errors"
	"github.com/dustin/go-humanize"
	"github.com/disintegration/imaging"
)

type Size interface {
	Size() int64 //字节
}

// 图片上传器
type ImageUploader struct {
	UploadPath string //图片保存的磁盘目录
	UrlPath string //上传图片的URL路径
	FileName string //文件名(不含扩展名)
	AllowExtName string //允许的扩展名
	Rename bool //是否重命名文件
	Cut bool //是否裁剪图片
	Resize bool //是否重新调整大小
	ToWidth int //要操作的图片宽度
	ToHeight int //要操作的图片高度
	MaxSize uint64 //图片大小,字节
}

// 创建一个新的图片上传器
func NewImageUploader(path string) *ImageUploader {
	return &ImageUploader{
		UploadPath: path,
		UrlPath: "/upload",
		AllowExtName: "gif,jpg,png,jpeg",
		Rename: true,
		MaxSize: 2000000, //2M
	}
}

// 上传图片
func (this *ImageUploader) Upload(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		file, handler, err := r.FormFile("file")
		if err != nil {
			return "", err
		}
		defer file.Close()

		//检查图片大小
		if sizeInterface, ok := file.(Size); ok {
			size := uint64(sizeInterface.Size())
			if size > this.MaxSize {
				return "", errors.New("图片过大，请换一张小于 " + humanize.Bytes(size) + " 的图片重新上传")
			}
		} else {
			return "", errors.New("图片尺寸超限或未能取得图片大小，请尝试换一张图片")
		}

		if !DirIsExist(this.UploadPath) {
			err := os.MkdirAll(this.UploadPath, os.ModePerm)
			if err != nil {
				return "", err
			}
		}

		extname := strings.ToLower(path.Ext(handler.Filename))
		allowExt := strings.Split(this.AllowExtName, ",")
		if InArray(Substr(extname, 1), allowExt) == -1 {
			return "", errors.New("不允许上传该格式的文件")
		}

		filename := this.FileName
		if filename == "" {
			filename = path.Base(handler.Filename)
			//唯一文件名
			if this.Rename {
				filename = Md5Sum(strconv.FormatInt(time.Now().UnixNano(), 10))
			}
		}

		if (this.Cut || this.Resize) && this.ToWidth > 0 && this.ToHeight > 0 {
			srcImg, err := imaging.Decode(file)
			if err != nil {
				return "", err
			}

			if this.Cut {
				srcImg = imaging.CropAnchor(srcImg, this.ToWidth, this.ToHeight, imaging.Center)
			}

			if this.Resize {
				srcImg = imaging.Resize(srcImg, this.ToWidth, this.ToHeight, imaging.Lanczos)
			}
			err = imaging.Save(srcImg, this.UploadPath + filename + extname)
			if err != nil {
				return "", err
			}
		} else {
			f, err := os.OpenFile(this.UploadPath + filename + extname, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return "", err
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				return "", err
			}
		}
		return filename + extname, nil;
	}
	return "", errors.New("没有发现上传的图片")
}