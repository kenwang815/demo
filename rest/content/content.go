package content

type Content interface {
	Code(int) Content
	Msg(string) Content
	Data(interface{}) Content
}

type content map[string]interface{}

func (c content) Code(code int) Content {
	c["code"] = code
	return c
}

func (c content) Msg(msg string) Content {
	c["msg"] = msg
	return c
}

func (c content) Data(d interface{}) Content {
	c["data"] = d
	return c
}

func NewContent() Content {
	c := make(content)
	return c
}
