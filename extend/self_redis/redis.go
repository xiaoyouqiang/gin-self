package self_redis

import (
	"context"
	"gin-self/extend/config"
	"gin-self/extend/self_loger"
	"github.com/go-redis/redis/v7"
	"time"
)

type resource struct {
	ctx    context.Context //记录日志时 传入 上下文
	client *redis.Client
}

var connPool = make(map[string]*resource) //redis 连接对象map 可存储多个 redis连接实例

//Open 建立连接
func Open(instanceName string) {
	if _, ok := connPool[instanceName]; ok {
		return
	}

	section := "redis-" + instanceName

	if config.Section(section) == nil {
		//没有配置redis章节
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Get(section, "addr").String(),
		Password:     config.Get(section, "password").String(),
		DB:           config.Get(section, "db").MustInt(0),
		MaxRetries:   config.Get(section, "max_try").MustInt(3),
		PoolSize:     config.Get(section, "pool_size").MustInt(10),
		MinIdleConns: config.Get(section, "min_conn").MustInt(10),
		ReadTimeout:  time.Duration(config.Get(section, "read_timeout").MustInt(2)) * time.Second,
		WriteTimeout:  time.Duration(config.Get(section, "write_timeout").MustInt(2)) * time.Second,
	})

	if err := client.Ping().Err(); err != nil {
		panic("redis connect error:" + err.Error())
	}

	connPool[instanceName] = &resource{
		client: client,
	}
}

//GetConn 获得连接
func GetConn(instanceName string) *resource {
	if _, ok := connPool[instanceName]; !ok {
		panic("redis not init,maybe no redis config")
	}

	conn := connPool[instanceName]

	//end

	return conn
}

//WithContext 打开日志记录模式
func WithContext(ctx context.Context)  {
	for k, _ := range connPool {
		connPool[k].client = connPool[k].client.WithContext(ctx)
	}
}

// Set set some <key,value> into redis
func (c *resource) Set(key, value string, ttl time.Duration) error {
	ts := time.Now()
	defer func() {
		ctx := c.client.Context()
		if ctx.Value("trace") == nil {
			return
		}

		ctx.Value("trace").(*self_loger.TraceData).AddRedisLog(
				time.Now().Format("2006/01/02 15:04:05"),
				"set",
				key,
				value,
				ttl.Seconds(),
				time.Since(ts).Seconds(),
			)
	}()

	if err := c.client.Set(key, value, ttl).Err(); err != nil {
		return err
	}

	return nil
}

//Get some key from redis
func (c *resource) Get(key string) (string, error) {
	ts := time.Now()
	defer func() {
		ctx := c.client.Context()
		if ctx.Value("trace") == nil {
			return
		}

		ctx.Value("trace").(*self_loger.TraceData).AddRedisLog(
				time.Now().Format("2006/01/02 15:04:05"),
				"get",
				key,
				"",
				0,
				time.Since(ts).Seconds(),
				)
	}()

	value, err := c.client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}


//TTL
func (c *resource) TTL(key string) (time.Duration, error) {
	ttl, err := c.client.TTL(key).Result()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

// Expire expire some key
func (c *resource) Expire(key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *resource) ExpireAt(key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(key, ttl).Result()
	return ok
}

func (c *resource) Exists(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := c.client.Exists(keys...).Result()
	return value > 0
}

func (c *resource) Del(key string) bool {
	ts := time.Now()

	defer func() {
		defer func() {
			ctx := c.client.Context()
			if ctx.Value("trace") == nil {
				return
			}

			ctx.Value("trace").(*self_loger.TraceData).AddRedisLog(
				time.Now().Format("2006/01/02 15:04:05"),
				"del",
				key,
				"",
				0,
				time.Since(ts).Seconds(),
			)
		}()
	}()

	if key == "" {
		return true
	}

	value, _ := c.client.Del(key).Result()
	return value > 0
}

func (c *resource) Incr(key string) int64 {
	ts := time.Now()

	defer func() {
		ctx := c.client.Context()
		if ctx.Value("trace") == nil {
			return
		}

		ctx.Value("trace").(*self_loger.TraceData).AddRedisLog(
			time.Now().Format("2006/01/02 15:04:05"),
			"Incr",
			key,
			"",
			0,
			time.Since(ts).Seconds(),
		)
	}()

	value, _ := c.client.Incr(key).Result()
	return value
}

//GetClient 没有封装的命令 可用该方法 直接调用 client的原生方法
func (c *resource) GetClient() *redis.Client {
	return c.client
}
