package sharding

type ShardRouter struct {
	NumShards int
}

func NewShardRouter(numShard int) *ShardRouter {
	return &ShardRouter{
		NumShards: numShard,
	}
}

func (sr *ShardRouter) GetShard(key int64) int {
	return int(key % int64(sr.NumShards))
}
