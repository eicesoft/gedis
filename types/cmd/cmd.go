package cmd

const (
	// System commands
	Auth         = "auth"
	Subscribe    = "subscribe"
	Publish      = "publish"
	UnSubscribe  = "unsubscribe"
	BgRewriteAof = "bgrewriteaof"
	RewriteAof   = "rewriteaof"
	FlushAll     = "flushall"
	Select       = "select"
	Ping         = "ping"
	Info         = "info"

	// Server commands
	Multi   = "multi"
	Discard = "discard"
	Exec    = "exec"
	Watch   = "watch"

	// Hash commands
	HSet    = "HSet"
	HSetNX  = "HSetNX"
	HExists = "HExists"
	HGet    = "HGet"
	HDel    = "HDel"
	HLen    = "HLen"
	HMSet   = "HMSet"
	HMGet   = "HMGet"
	HKeys   = "HKeys"
	HVals   = "HVals"
	HGetAll = "HGetAll"
	HIncrBy = "HIncrBy"

	// Key commands
	Del       = "Del"
	Expire    = "Expire"
	ExpireAt  = "ExpireAt"
	PExpire   = "PExpire"
	PExpireAt = "PExpireAt"
	TTL       = "TTL"
	PTTL      = "PTTL"
	Persist   = "Persist"
	Exists    = "Exists"
	Type      = "Type"
	Rename    = "Rename"
	RenameNx  = "RenameNx"
	FlushDB   = "FlushDB"
	Keys      = "Keys"
	Scan      = "Scan"

	LPush     = "LPush"
	LPushX    = "LPushX"
	RPush     = "RPush"
	RPushX    = "RPushX"
	LPop      = "LPop"
	RPop      = "RPop"
	RPopLPush = "RPopLPush"
	LRem      = "LRem"
	LLen      = "LLen"
	LIndex    = "LIndex"
	LSet      = "LSet"
	LRange    = "LRange"

	//Set commands
	SAdd        = "SAdd"
	SIsMember   = "SIsMember"
	SRem        = "SRem"
	SCard       = "SCard"
	SMembers    = "SMembers"
	SInter      = "SInter"
	SInterStore = "SInterStore"
	SUnion      = "SUnion"
	SUnionStore = "SUnionStore"
	SDiff       = "SDiff"
	SDiffStore  = "SDiffStore"
	SRandMember = "SRandMember"

	//String commands
	Set         = "Set"
	SetNx       = "SetNx"
	SetEX       = "SetEX"
	PSetEX      = "PSetEX"
	MSet        = "MSet"
	MGet        = "MGet"
	MSetNX      = "MSetNX"
	Get         = "Get"
	GetSet      = "GetSet"
	Incr        = "Incr"
	IncrBy      = "IncrBy"
	IncrByFloat = "IncrByFloat"
	Decr        = "Decr"
	DecrBy      = "DecrBy"
	StrLen      = "StrLen"
	Append      = "Append"
	SetRange    = "SetRange"
	GetRange    = "GetRange"

	//SortedSet commands
	ZAdd             = "ZAdd"
	ZScore           = "ZScore"
	ZIncrBy          = "ZIncrBy"
	ZRank            = "ZRank"
	ZCount           = "ZCount"
	ZRevRank         = "ZRevRank"
	ZCard            = "ZCard"
	ZRange           = "ZRange"
	ZRangeByScore    = "ZRangeByScore"
	ZRevRange        = "ZRevRange"
	ZRevRangeByScore = "ZRevRangeByScore"
	ZRem             = "ZRem"
	ZRemRangeByScore = "ZRemRangeByScore"
	ZRemRangeByRank  = "ZRemRangeByRank"
)
