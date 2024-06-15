package entity

// 定义字段结构体
type Field struct {
	Name     string
	Type     string
	Tag      string
	Validate string
}

// 定义消息结构体
type Message struct {
	Parent *Message
	Name   string
	Fields []*Field
	Method string
}

func (msg Message) AttachTag(tag string) {
	for _, field := range msg.Fields {
		field.Tag = tag
	}
}

// 定义服务处理函数结构体
type HandlerFunc struct {
	Method   string
	Name     string
	Request  string
	Response string
}

// 定义服务结构体
type Service struct {
	Name     string
	Handlers []HandlerFunc
}

type CodeCompnent struct {
	Module   string
	Imports  []string
	Services map[string]*Service
	Messages map[string]*Message
}
