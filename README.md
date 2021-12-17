# Gedis

一个Redis的服务端实现, 完全兼容redis客户端, 包含基本的命令, 以及数据结构(string, list, map, set, zset等)的实现, 文件存储实现了Aof和leveldb作为数据落地的方式. 改项目不是为了替代Redis, 只是为了测试GO的内存和网络编程的上限.

实现了一个线程安全的dict, 以及双向链表实现的list, 以及基于dict的set结构, sortedset基于skiplist

基本上实现了主流命令, 包含: 
```
subscribe,publish,unsubscribe,bgrewriteaof,rewriteaof,flushall,select,ping,info,multi,discard,exec,watch,HSet,HSetNX,HExists,HGet,HDel,HLen,HMSet,HMGet,HKeys,HVals,HGetAll,HIncrBy,Del,Expire,ExpireAt,PExpire,PExpireAt,TTL,PTTL,Persist,Exists,Type,Rename,RenameNx,FlushDB,Keys,Scan,LPush,LPushX,RPush,RPushX,LPop,RPop,RPopLPush,LRem,LLen,LIndex,LSet,LRange,SAdd,SIsMember,SRem,SCard,SMembers,SInter,SInterStore,SUnion,SUnionStore,SDiff,SDiffStore,SRandMember,Set,SetNx,SetEX,PSetEX,MSet,MGet,MSetNX,Get,GetSet,Incr,IncrBy,IncrByFloat,Decr,DecrBy,StrLen,Append,SetRange,GetRange,ZAdd,ZScore,ZIncrBy,ZRank,ZCount,ZRevRank,ZCard,ZRange,ZRangeByScore,ZRevRange,ZRevRangeByScore,ZRem,ZRemRangeByScore,ZRemRangeByRank
```
对于日常应用基本上够用了.

