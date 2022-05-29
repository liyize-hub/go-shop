package datasource

import (
	"go-shop/config"
	"go-shop/models"
	"go-shop/utils"
	"net/http"
	"net/url"

	"github.com/kataras/iris/v12"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
)

var Imgc *cos.Client
var Htmlc *cos.Client

func GeImg(ctx iris.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.Logger.Error("获取图片失败", zap.Any("Error", err))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "获取图片失败", "")
		return
	}

	//上传文件到cos
	_, err = Imgc.Object.Put(ctx, fileHeader.Filename, file, nil)
	if err != nil {
		utils.Logger.Error("上传图片到cos失败", zap.Any("Error", err))
		utils.SendJSON(ctx, models.ErrorCode.ERROR, "上传图片到cos失败", "")
		return
	}

	utils.SendJSON(ctx, models.ErrorCode.SUCCESS, "图片上传成功", iris.Map{
		"value": config.ServerConfig.ImgURL + "/" + fileHeader.Filename,
	})
}

// cos img存储桶客户端
func NewImgCOS() {
	u, _ := url.Parse(config.ServerConfig.ImgURL)
	b := &cos.BaseURL{BucketURL: u}
	Imgc = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.ServerConfig.SecretID,
			SecretKey: config.ServerConfig.SecretKey,
		},
	})
}

// cos 静态页面存储桶客户端
func NewHTMLCOS() {
	u, _ := url.Parse(config.ServerConfig.HtmlURL)
	b := &cos.BaseURL{BucketURL: u}
	Htmlc = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.ServerConfig.SecretID,
			SecretKey: config.ServerConfig.SecretKey,
		},
	})
}

func init() {
	NewImgCOS()
	NewHTMLCOS()
}

/*func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}*/
