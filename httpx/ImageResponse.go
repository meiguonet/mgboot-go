package httpx

import "io/ioutil"

type ImageResponse struct {
	buf      []byte
	mimeType string
}

func NewImageResponseFromFile(fpath, mimeType string) ImageResponse {
	buf, _ := ioutil.ReadFile(fpath)

	return ImageResponse{
		buf:      buf,
		mimeType: mimeType,
	}
}

func NewImageResponseFromBuffer(buf []byte, mimeType string) ImageResponse {
	return ImageResponse{
		buf:      buf,
		mimeType: mimeType,
	}
}

func (p ImageResponse) GetContentType() string {
	return p.mimeType
}

func (p ImageResponse) GetContents() string {
	return ""
}

func (p ImageResponse) Buffer() []byte {
	return p.buf
}
