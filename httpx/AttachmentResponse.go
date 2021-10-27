package httpx

import (
	"io/ioutil"
)

type AttachmentResponse struct {
	buf                []byte
	attachmentFileName string
}

func NewAttachmentResponseFromFile(fpath, attachmentFileName string) AttachmentResponse {
	buf, _ := ioutil.ReadFile(fpath)

	return AttachmentResponse{
		buf:                buf,
		attachmentFileName: attachmentFileName,
	}
}

func NewAttachmentResponseFromBuffer(buf []byte, attachmentFileName string) AttachmentResponse {
	return AttachmentResponse{
		buf:                buf,
		attachmentFileName: attachmentFileName,
	}
}

func (p AttachmentResponse) GetContentType() string {
	return "application/octet-stream"
}

func (p AttachmentResponse) GetContents() string {
	return ""
}

func (p AttachmentResponse) Buffer() []byte {
	return p.buf
}

func (p AttachmentResponse) AttachmentFileName() string {
	return p.attachmentFileName
}
