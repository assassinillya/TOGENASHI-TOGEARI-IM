// Code generated by goctl. DO NOT EDIT.
package types

type FileRequest struct {
	UserID uint `header:"User-ID"`
}

type FileResponse struct {
	Src string `json:"src"`
}

type ImageRequest struct {
	UserID uint `header:"User-ID"`
}

type ImageResponse struct {
	Url string `json:"url"`
}

type ImageShowRequest struct {
	ImageType string `path:"imageType"`
	ImageName string `path:"imageName"`
}
