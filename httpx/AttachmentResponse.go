package httpx

import (
	"github.com/meiguonet/mgboot-go-common/util/mimex"
	"io/ioutil"
)

type AttachmentResponse struct {
	buf                []byte
	mimeType           string
	attachmentFileName string
}

func NewAttachmentResponseFromFile(fpath, attachmentFileName string, mimeType ...string) AttachmentResponse {
	buf, _ := ioutil.ReadFile(fpath)
	var _mimeType string

	if len(mimeType) > 0 {
		_mimeType = mimeType[0]
	}

	if _mimeType == "" {
		_mimeType = mimex.GetMimeType(buf)
	}

	return AttachmentResponse{
		buf:                buf,
		mimeType:           _mimeType,
		attachmentFileName: attachmentFileName,
	}
}

func NewAttachmentResponseFromBuffer(buf []byte, attachmentFileName string, mimeType ...string) AttachmentResponse {
	var _mimeType string

	if len(mimeType) > 0 {
		_mimeType = mimeType[0]
	}

	if _mimeType == "" {
		_mimeType = mimex.GetMimeType(buf)
	}

	return AttachmentResponse{
		buf:                buf,
		mimeType:           _mimeType,
		attachmentFileName: attachmentFileName,
	}
}

func (p AttachmentResponse) GetContentType() string {
	if p.mimeType == "" {
		return "application/octet-stream"
	}

	return p.mimeType
}

func (p AttachmentResponse) GetContents() interface{} {
	return ""
}

func (p AttachmentResponse) Buffer() []byte {
	return p.buf
}

func (p AttachmentResponse) AttachmentFileName() string {
	return p.attachmentFileName
}
