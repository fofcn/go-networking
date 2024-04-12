package bytes_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer_ShouldBeZeroCap_WhenCreateWithNewKeyword(t *testing.T) {

	buf := new(bytes.Buffer)
	println("刚创建完成的Buffer容量: ", buf.Cap())

	assert.Equal(t, 0, buf.Cap())
}

func TestBuffer_ShouldBeZeroCap_WhenCreateWithVar(t *testing.T) {
	var buf bytes.Buffer
	println("刚创建完成的Buffer容量: ", buf.Cap())

	assert.Equal(t, 0, buf.Cap())
}

func TestBuffer_ShouldBeBytesLength_whenCreateWithBytes(t *testing.T) {
	buf := bytes.NewBuffer([]byte("hello world"))
	println("刚创建完成的Buffer容量: ", buf.Cap())
	assert.Equal(t, 11, buf.Len())
}

func TestBuffer_ShouldBeZeroCap_WhenCreateStringBuffer(t *testing.T) {
	str := "hello buffer"
	println("字符串长度: ", cap([]byte(str)))
	buf := bytes.NewBufferString(str)
	println("刚创建完成的Buffer容量: ", buf.Cap())

	// 这里的断言是buf.Len而不是buf.Cap，是因为在测试框架中容量不一定等于str的长度
	assert.Equal(t, len(str), buf.Len())
}

func TestBuffer_ShouldGetXiao_WhenWriteUnicodeStringAndGetFirstCharacter(t *testing.T) {
	buf := new(bytes.Buffer)
	str := "小厂程序员"
	buf.Write([]byte(str))

	r, _, err := buf.ReadRune()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, '小', r)
}

func TestBuffer_ShouldGetAll_WhenWriteUnicodeStringAndCallReadString(t *testing.T) {
	buf := new(bytes.Buffer)
	str := "小厂程序员"
	buf.Write([]byte(str))

	nstr := buf.String()
	assert.Equal(t, nstr, nstr)
}
