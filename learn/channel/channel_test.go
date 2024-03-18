package channel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChan_ShouldRecvValues_WhenSendValuesToChannel(t *testing.T) {

	c := make(chan int)

	// given
	s := []int{1, 2, 3, 4, 5, 6}

	// when
	go sum(s[:len(s)/2], c)
	ret1 := <-c

	// should
	assert.Equal(t, 6, ret1)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}
