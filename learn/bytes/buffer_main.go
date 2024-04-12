package bytes

import (
	"bytes"
)

func main() {
	str := "hello buffer"
	strbytes := []byte(str)
	print(len(strbytes))
	buf := bytes.NewBufferString(str)
	println("刚创建完成的Buffer容量: ", buf.Cap())
}
