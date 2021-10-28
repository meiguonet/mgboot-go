package BuiltintResponse

import (
	"github.com/meiguonet/mgboot-go-common/util/mimex"
	"io/ioutil"
)

type ImageResponse struct {
	buf      []byte
	mimeType string
}

func NewImageResponseFromFile(fpath string, mimeType ...string) ImageResponse {
	buf, _ := ioutil.ReadFile(fpath)
	var _mimeType string

	if len(mimeType) > 0 {
		_mimeType = mimeType[0]
	}

	if _mimeType == "" {
		_mimeType = mimex.GetMimeType(buf)
	}

	return ImageResponse{
		buf:      buf,
		mimeType: _mimeType,
	}
}

func NewImageResponseFromBuffer(buf []byte, mimeType ...string) ImageResponse {
	var _mimeType string

	if len(mimeType) > 0 {
		_mimeType = mimeType[0]
	}

	if _mimeType == "" {
		_mimeType = mimex.GetMimeType(buf)
	}

	return ImageResponse{
		buf:      buf,
		mimeType: _mimeType,
	}
}

func (p ImageResponse) GetContentType() string {
	return p.mimeType
}

func (p ImageResponse) GetContents() interface{} {
	return ""
}

func (p ImageResponse) Buffer() []byte {
	return p.buf
}
