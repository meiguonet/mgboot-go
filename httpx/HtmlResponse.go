package httpx

type HtmlResponse struct {
	contents string
}

func NewHtmlResponse(contents string) HtmlResponse {
	return HtmlResponse{contents: contents}
}

func (p HtmlResponse) GetContentType() string {
	return "text/html; charset=utf-8"
}

func (p HtmlResponse) GetContents() string {
	return p.contents
}
