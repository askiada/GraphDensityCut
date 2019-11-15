package session

import (
	"testing"

	"github.com/askiada/GraphDensityCut/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SessionTestSuite struct {
	suite.Suite
}

func (suite *SessionTestSuite) SetupSuite() {
}

func (suite *SessionTestSuite) SetupTest() {

}

func (suite *SessionTestSuite) TestNodeChainRemoveEmpty() {
	c := &NodeChain{}
	c.Remove()
	assert.Nil(suite.T(), c.Prev)
	assert.Nil(suite.T(), c.Next)
}

func (suite *SessionTestSuite) TestNodeChainEmptyPushFront() {
	c := &NodeChain{}
	second := c.PushFront(&model.Node{Index: 1, Value: "1"})
	assert.Nil(suite.T(), c.Prev)
	assert.NotEqual(suite.T(), second, c)
	assert.Equal(suite.T(), second, c.Next)
	assert.Equal(suite.T(), c, second.Prev)
	assert.Nil(suite.T(), second.Next)

	between := c.PushFront(&model.Node{Index: 2, Value: "2"})
	assert.Nil(suite.T(), c.Prev)
	assert.NotEqual(suite.T(), between, c)
	assert.NotEqual(suite.T(), between, second)
	assert.Equal(suite.T(), between, c.Next)
	assert.Equal(suite.T(), c, between.Prev)
	assert.Equal(suite.T(), second, between.Next)
	assert.Equal(suite.T(), between, second.Prev)

	last := second.PushFront(&model.Node{Index: 3, Value: "3"})
	assert.NotEqual(suite.T(), last, c)
	assert.NotEqual(suite.T(), last, second)
	assert.NotEqual(suite.T(), last, between)
	assert.Nil(suite.T(), last.Next)
	assert.Equal(suite.T(), second, last.Prev)
	assert.Equal(suite.T(), last, second.Next)

	head := &NodeChain{}
	tail := head.PushFront(&model.Node{Index: 1, Value: "1"}).PushFront(&model.Node{Index: 2, Value: "2"}).PushFront(&model.Node{Index: 3, Value: "3"})
	count := 0
	for {
		if head != nil {
			head = head.Next
		} else {
			break
		}
		count++
	}
	assert.Equal(suite.T(), 4, count)

	count = 0
	for {
		if tail != nil {
			tail = tail.Prev
		} else {
			break
		}
		count++
	}
	assert.Equal(suite.T(), 4, count)
}

func (suite *SessionTestSuite) TestNodeChainRemove() {

	head := &NodeChain{}
	first := head.PushFront(&model.Node{Index: 1, Value: "1"})
	second := first.PushFront(&model.Node{Index: 2, Value: "2"})
	tail := second.PushFront(&model.Node{Index: 3, Value: "3"})
	count := 0

	tmp := head
	for {
		if tmp != nil {
			tmp = tmp.Next
		} else {
			break
		}
		count++
	}
	assert.Equal(suite.T(), 4, count)

	new := head.Remove()

	assert.Equal(suite.T(), first, new)
	count = 0
	for {
		if new != nil {
			new = new.Next
		} else {
			break
		}
		count++
	}
	assert.Equal(suite.T(), 3, count)

	new = second.Remove()
	assert.Equal(suite.T(), first, new)
	count = 0
	for {
		if new != nil {
			new = new.Next
		} else {
			break
		}
		count++
	}
	assert.Equal(suite.T(), 2, count)

	new = first.Remove()
	assert.Equal(suite.T(), tail, new)
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}
