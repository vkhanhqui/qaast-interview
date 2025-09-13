package errors

type RenderError struct {
	Message string
	Status  int
}

type RenderedError interface {
	IsRendered() bool
}

func (r *RenderError) Error() string {
	return r.Message
}

func (r *RenderError) ResponseBody() ([]byte, error) {
	return []byte(r.Message), nil
}

func (r *RenderError) ResponseHeaders() (int, map[string]string) {
	return r.Status, map[string]string{"Content-Type": "text/html; charset=utf-8"}
}

func (r *RenderError) IsRendered() bool {
	return true
}
