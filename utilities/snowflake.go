package utilities

import (
	"math/rand"

	"github.com/bwmarrin/snowflake"
)

//go:generate mockery --name=SnowflakeIDGenerator --output=mocks
type SnowflakeIDGenerator interface {
	Next() snowflake.ID
}

type snowflakeIDGeneratorImpl struct {
	node *snowflake.Node
}

func NewSnowflakeIDGenerator() (SnowflakeIDGenerator, error) {
	node, err := snowflake.NewNode(int64(rand.Int31n(1024)))
	if err != nil {
		return nil, err
	}

	return &snowflakeIDGeneratorImpl{node}, nil
}

func (impl *snowflakeIDGeneratorImpl) Next() snowflake.ID {
	return impl.node.Generate()
}
