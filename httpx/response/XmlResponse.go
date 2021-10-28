package BuiltintResponse

type XmlResponse struct {
	contents string
}

func NewXmlResponse(contents string) XmlResponse {
	return XmlResponse{contents: contents}
}

func (p XmlResponse) GetContentType() string {
	return "text/xml; charset=utf-8"
}

func (p XmlResponse) GetContents() interface{} {
	return p.contents
}
