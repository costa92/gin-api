package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/costa92/errors"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

type Config struct {
	Host                  string
	Port                  int
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	EnableCluster         bool
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

var ErrRedisIsDown = errors.New("db: Redis is either down or ws not configured")

var (
	singlePool      atomic.Value
	singleCachePool atomic.Value
	redisUp         atomic.Value
)

var disableRedis atomic.Value

func DisableRedis(ok bool) {
	if ok {
		redisUp.Store(false)
		disableRedis.Store(true)
		return
	}
	redisUp.Store(true)
	disableRedis.Store(false)
}

func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}
	return true
}

func Connected() bool {
	if v := redisUp.Load(); v != nil {
		return v.(bool)
	}
	return false
}

func singleton(cache bool) redis.UniversalClient {
	if cache {
		v := singleCachePool.Load()
		if v != nil {
			return v.(redis.UniversalClient)
		}
		return nil
	}
	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}
	return nil
}

func connectSingleton(cache bool, config *Config) bool {
	if singleton(cache) == nil {
		if cache {
			singleCachePool.Store(NewRedisClusterPool(cache, config))
			return true
		}

		singlePool.Store(NewRedisClusterPool(cache, config))
		return true
	}
	return true
}

type RedisCluster struct {
	KeyPrefix string
	HasKeys   bool
	IsCache   bool
}

func clusterConnectionIsOpen(cluster RedisCluster) bool {
	c := singleton(cluster.IsCache)
	testKey := "redis-test-" + uuid.Must(uuid.NewV4(), nil).String()
	if err := c.Set(testKey, "test", time.Second).Err(); err != nil {
		return false
	}

	if _, err := c.Get(testKey).Result(); err != nil {
		return false
	}
	return true
}

func ConnectToRedis(ctx context.Context, config *Config) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	c := []RedisCluster{
		{}, {IsCache: true},
	}
	var ok bool
	for _, v := range c {
		if !connectSingleton(v.IsCache, config) {
			break
		}

		if !clusterConnectionIsOpen(v) {
			redisUp.Store(false)
			break
		}
		ok = true
	}
	redisUp.Store(ok)
again:
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if !shouldConnect() {
				continue
			}
			for _, v := range c {
				if !connectSingleton(v.IsCache, config) {
					redisUp.Store(false)
					goto again
				}
				if !clusterConnectionIsOpen(v) {
					redisUp.Store(false)
					goto again
				}
			}
			redisUp.Store(true)
		}
	}
}

func NewRedisClusterPool(isCache bool, config *Config) redis.UniversalClient {
	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	timeout := 5 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}
	var tlsConfig *tls.Config

	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient

	opts := &RedisOpts{
		Addrs:        getRedisAddrs(config),
		MasterName:   config.MasterName,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  240 * timeout,
		PoolSize:     poolSize,
		TLSConfig:    tlsConfig,
	}
	if opts.MasterName != "" {
		client = redis.NewFailoverClient(opts.failover())
	} else if config.EnableCluster {
		client = redis.NewClusterClient(opts.cluster())
	} else {
		client = redis.NewClient(opts.simple())
	}

	return client
}

func getRedisAddrs(config *Config) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}
	return addrs
}

type RedisOpts redis.UniversalOptions

func (o *RedisOpts) cluster() *redis.ClusterOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:6379"}
	}

	return &redis.ClusterOptions{
		Addrs:     o.Addrs,
		OnConnect: o.OnConnect,

		Password: o.Password,

		MaxRedirects:   o.MaxRedirects,
		ReadOnly:       o.ReadOnly,
		RouteByLatency: o.RouteByLatency,
		RouteRandomly:  o.RouteRandomly,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:        o.DialTimeout,
		ReadTimeout:        o.ReadTimeout,
		WriteTimeout:       o.WriteTimeout,
		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) simple() *redis.Options {
	addr := "127.0.0.1:6379"
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) failover() *redis.FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}
	return &redis.FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		OnConnect:     o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,
		TLSConfig:          o.TLSConfig,
	}
}

func (r *RedisCluster) Connect() bool {
	return true
}

func (r *RedisCluster) singleton() redis.UniversalClient {
	return singleton(r.IsCache)
}

func (r *RedisCluster) hashKey(in string) string {
	if !r.HasKeys {
		return in
	}
	return HashStr(in)
}

func (r *RedisCluster) fixKey(keyName string) string {
	return r.KeyPrefix + r.hashKey(keyName)
}

func (r *RedisCluster) cleanKey(keyName string) string {
	return strings.Replace(keyName, r.KeyPrefix, "", 1)
}

func (r *RedisCluster) up() error {
	if !Connected() {
		return ErrRedisIsDown
	}

	return nil
}

// GetKey will retrieve a key from the database.
func (r *RedisCluster) GetKey(keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}

	cluster := r.singleton()

	value, err := cluster.Get(r.fixKey(keyName)).Result()
	if err != nil {
		return "", ErrKeyNotFound
	}

	return value, nil
}

// GetMultiKey gets multiple keys from the database.
func (r *RedisCluster) GetMultiKey(keys []string) ([]string, error) {
	if err := r.up(); err != nil {
		return nil, err
	}
	cluster := r.singleton()
	keyNames := make([]string, len(keys))
	copy(keyNames, keys)
	for index, val := range keyNames {
		keyNames[index] = r.fixKey(val)
	}

	result := make([]string, 0)

	switch v := cluster.(type) {
	case *redis.ClusterClient:
		{
			getCmds := make([]*redis.StringCmd, 0)
			pipe := v.Pipeline()
			for _, key := range keyNames {
				getCmds = append(getCmds, pipe.Get(key))
			}
			_, err := pipe.Exec()
			if err != nil && !errors.Is(err, redis.Nil) {
				return nil, ErrKeyNotFound
			}
			for _, cmd := range getCmds {
				result = append(result, cmd.Val())
			}
		}
	case *redis.Client:
		{
			values, err := cluster.MGet(keyNames...).Result()
			if err != nil {
				return nil, ErrKeyNotFound
			}
			for _, val := range values {
				strVal := fmt.Sprint(val)
				if strVal == "<nil>" {
					strVal = ""
				}
				result = append(result, strVal)
			}
		}
	}

	for _, val := range result {
		if val != "" {
			return result, nil
		}
	}

	return nil, ErrKeyNotFound
}

// GetKeyTTL return ttl of the given key.
func (r *RedisCluster) GetKeyTTL(keyName string) (ttl int64, err error) {
	if err = r.up(); err != nil {
		return 0, err
	}
	duration, err := r.singleton().TTL(r.fixKey(keyName)).Result()

	return int64(duration.Seconds()), err
}

// GetRawKey return the value of the given key.
func (r *RedisCluster) GetRawKey(keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}
	value, err := r.singleton().Get(keyName).Result()
	if err != nil {
		return "", ErrKeyNotFound
	}

	return value, nil
}

// GetExp return the expiry of the given key.
func (r *RedisCluster) GetExp(keyName string) (int64, error) {
	if err := r.up(); err != nil {
		return 0, err
	}

	value, err := r.singleton().TTL(r.fixKey(keyName)).Result()
	if err != nil {
		return 0, ErrKeyNotFound
	}

	return int64(value.Seconds()), nil
}

// SetExp set expiry of the given key.
func (r *RedisCluster) SetExp(keyName string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Expire(r.fixKey(keyName), timeout).Err()
	if err != nil {
	}
	return err
}

// SetKey will create (or update) a key value in the store.
func (r *RedisCluster) SetKey(keyName, session string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Set(r.fixKey(keyName), session, timeout).Err()
	if err != nil {
		return err
	}

	return nil
}

// SetRawKey set the value of the given key.
func (r *RedisCluster) SetRawKey(keyName, session string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Set(keyName, session, timeout).Err()
	if err != nil {
		return err
	}

	return nil
}

// Decrement will decrement a key in redis.
func (r *RedisCluster) Decrement(keyName string) {
	keyName = r.fixKey(keyName)
	if err := r.up(); err != nil {
		return
	}
	err := r.singleton().Decr(keyName).Err()
	if err != nil {
	}
}

// IncrememntWithExpire will increment a key in redis.
func (r *RedisCluster) IncrememntWithExpire(keyName string, expire int64) int64 {
	if err := r.up(); err != nil {
		return 0
	}
	// This function uses a raw key, so we shouldn't call fixKey
	fixedKey := keyName
	val, err := r.singleton().Incr(fixedKey).Result()

	if err != nil {
		fmt.Errorf("error trying to increment value: %s", err.Error())
	} else {
		fmt.Errorf("incremented key: %s, val is: %d", fixedKey, val)
	}

	if val == 1 && expire > 0 {
		fmt.Print("--> Setting Expire")
		r.singleton().Expire(fixedKey, time.Duration(expire)*time.Second)
	}

	return val
}

// GetKeys will return all keys according to the filter (filter is a prefix - e.g. tyk.keys.*).
func (r *RedisCluster) GetKeys(filter string) []string {
	if err := r.up(); err != nil {
		return nil
	}
	client := r.singleton()

	filterHash := ""
	if filter != "" {
		filterHash = r.hashKey(filter)
	}
	searchStr := r.KeyPrefix + filterHash + "*"

	fnFetchKeys := func(client *redis.Client) ([]string, error) {
		values := make([]string, 0)

		iter := client.Scan(0, searchStr, 0).Iterator()
		for iter.Next() {
			values = append(values, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return nil, err
		}

		return values, nil
	}

	var err error
	var values []string
	sessions := make([]string, 0)

	switch v := client.(type) {
	case *redis.ClusterClient:
		ch := make(chan []string)

		go func() {
			err = v.ForEachMaster(func(client *redis.Client) error {
				values, err = fnFetchKeys(client)
				if err != nil {
					return err
				}

				ch <- values

				return nil
			})
			close(ch)
		}()

		for res := range ch {
			sessions = append(sessions, res...)
		}
	case *redis.Client:
		sessions, err = fnFetchKeys(v)
	}

	if err != nil {
		return nil
	}

	for i, v := range sessions {
		sessions[i] = r.cleanKey(v)
	}

	return sessions
}

// GetKeysAndValuesWithFilter will return all keys and their values with a filter.
func (r *RedisCluster) GetKeysAndValuesWithFilter(filter string) map[string]string {
	if err := r.up(); err != nil {
		return nil
	}
	keys := r.GetKeys(filter)
	if keys == nil {
		return nil
	}

	if len(keys) == 0 {
		return nil
	}

	for i, v := range keys {
		keys[i] = r.KeyPrefix + v
	}

	client := r.singleton()
	values := make([]string, 0)

	switch v := client.(type) {
	case *redis.ClusterClient:
		{
			getCmds := make([]*redis.StringCmd, 0)
			pipe := v.Pipeline()
			for _, key := range keys {
				getCmds = append(getCmds, pipe.Get(key))
			}
			_, err := pipe.Exec()
			if err != nil && !errors.Is(err, redis.Nil) {
				return nil
			}

			for _, cmd := range getCmds {
				values = append(values, cmd.Val())
			}
		}
	case *redis.Client:
		{
			result, err := v.MGet(keys...).Result()
			if err != nil {
				fmt.Errorf("Error trying to get client keys: %s", err.Error())

				return nil
			}

			for _, val := range result {
				strVal := fmt.Sprint(val)
				if strVal == "<nil>" {
					strVal = ""
				}
				values = append(values, strVal)
			}
		}
	}

	m := make(map[string]string)
	for i, v := range keys {
		m[r.cleanKey(v)] = values[i]
	}

	return m
}

// GetKeysAndValues will return all keys and their values - not to be used lightly.
func (r *RedisCluster) GetKeysAndValues() map[string]string {
	return r.GetKeysAndValuesWithFilter("")
}

// DeleteKey will remove a key from the database.
func (r *RedisCluster) DeleteKey(keyName string) bool {
	if err := r.up(); err != nil {
		// log.Debug(err)
		return false
	}

	n, err := r.singleton().Del(r.fixKey(keyName)).Result()
	if err != nil {
		fmt.Errorf("error trying to delete key: %s", err.Error())
	}

	return n > 0
}

// DeleteAllKeys will remove all keys from the database.
func (r *RedisCluster) DeleteAllKeys() bool {
	if err := r.up(); err != nil {
		return false
	}
	n, err := r.singleton().FlushAll().Result()
	if err != nil {
		fmt.Errorf("error trying to delete keys: %s", err.Error())
	}

	if n == "OK" {
		return true
	}

	return false
}

// DeleteRawKey will remove a key from the database without prefixing, assumes user knows what they are doing.
func (r *RedisCluster) DeleteRawKey(keyName string) bool {
	if err := r.up(); err != nil {
		return false
	}
	n, err := r.singleton().Del(keyName).Result()
	if err != nil {
		fmt.Errorf("Error trying to delete key: %s", err.Error())
	}

	return n > 0
}

// DeleteScanMatch will remove a group of keys in bulk.
func (r *RedisCluster) DeleteScanMatch(pattern string) bool {
	if err := r.up(); err != nil {
		return false
	}
	client := r.singleton()

	fnScan := func(client *redis.Client) ([]string, error) {
		values := make([]string, 0)

		iter := client.Scan(0, pattern, 0).Iterator()
		for iter.Next() {
			values = append(values, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return nil, err
		}

		return values, nil
	}

	var err error
	var keys []string
	var values []string

	switch v := client.(type) {
	case *redis.ClusterClient:
		ch := make(chan []string)
		go func() {
			err = v.ForEachMaster(func(client *redis.Client) error {
				values, err = fnScan(client)
				if err != nil {
					return err
				}

				ch <- values

				return nil
			})
			close(ch)
		}()

		for vals := range ch {
			keys = append(keys, vals...)
		}
	case *redis.Client:
		keys, err = fnScan(v)
	}

	if err != nil {
		fmt.Errorf("SCAN command field with err: %s", err.Error())

		return false
	}

	if len(keys) > 0 {
		for _, name := range keys {

			err := client.Del(name).Err()
			if err != nil {
				fmt.Errorf("error trying to delete key: %s - %s", name, err.Error())
			}
		}
		fmt.Printf("Deleted: %d records", len(keys))
	} else {
		fmt.Printf("RedisCluster called DEL - Nothing to delete")
	}

	return true
}

// DeleteKeys will remove a group of keys in bulk.
func (r *RedisCluster) DeleteKeys(keys []string) bool {
	if err := r.up(); err != nil {
		return false
	}
	if len(keys) > 0 {
		for i, v := range keys {
			keys[i] = r.fixKey(v)
		}

		fmt.Printf("Deleting: %v", keys)
		client := r.singleton()
		switch v := client.(type) {
		case *redis.ClusterClient:
			{
				pipe := v.Pipeline()
				for _, k := range keys {
					pipe.Del(k)
				}

				if _, err := pipe.Exec(); err != nil {
					fmt.Errorf("error trying to delete keys: %s", err.Error())
				}
			}
		case *redis.Client:
			{
				_, err := v.Del(keys...).Result()
				if err != nil {
					fmt.Errorf("error trying to delete keys: %s", err.Error())
				}
			}
		}
	} else {
		fmt.Printf("RedisCluster called DEL - Nothing to delete")
	}

	return true
}

// StartPubSubHandler will listen for a signal and run the callback for
// every subscription and message event.
func (r *RedisCluster) StartPubSubHandler(channel string, callback func(interface{})) error {
	if err := r.up(); err != nil {
		return err
	}
	client := r.singleton()
	if client == nil {
		return errors.New("redis connection failed")
	}

	pubsub := client.Subscribe(channel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(); err != nil {
		fmt.Errorf("error while receiving pubsub message: %s", err.Error())

		return err
	}

	for msg := range pubsub.Channel() {
		callback(msg)
	}

	return nil
}

// Publish publish a message to the specify channel.
func (r *RedisCluster) Publish(channel, message string) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Publish(channel, message).Err()
	if err != nil {
		fmt.Errorf("error trying to set value: %s", err.Error())

		return err
	}

	return nil
}

// GetAndDeleteSet get and delete a key.
func (r *RedisCluster) GetAndDeleteSet(keyName string) []interface{} {
	fmt.Printf("Getting raw key set: %s", keyName)
	if err := r.up(); err != nil {
		return nil
	}
	fmt.Printf("keyName is: %s", keyName)
	fixedKey := r.fixKey(keyName)
	fmt.Printf("Fixed keyname is: %s", fixedKey)

	client := r.singleton()

	var lrange *redis.StringSliceCmd
	_, err := client.TxPipelined(func(pipe redis.Pipeliner) error {
		lrange = pipe.LRange(fixedKey, 0, -1)
		pipe.Del(fixedKey)

		return nil
	})
	if err != nil {
		fmt.Errorf("multi command failed: %s", err.Error())

		return nil
	}

	vals := lrange.Val()
	fmt.Printf("Analytics returned: %d", len(vals))
	if len(vals) == 0 {
		return nil
	}

	fmt.Printf("Unpacked vals: %d", len(vals))
	result := make([]interface{}, len(vals))
	for i, v := range vals {
		result[i] = v
	}

	return result
}

// AppendToSet append a value to the key set.
func (r *RedisCluster) AppendToSet(keyName, value string) {
	fixedKey := r.fixKey(keyName)
	if err := r.up(); err != nil {
		return
	}
	if err := r.singleton().RPush(fixedKey, value).Err(); err != nil {
		fmt.Printf("Error trying to append to set keys: %s", err.Error())
	}
}

// Exists check if keyName exists.
func (r *RedisCluster) Exists(keyName string) (bool, error) {
	fixedKey := r.fixKey(keyName)

	exists, err := r.singleton().Exists(fixedKey).Result()
	if err != nil {
		fmt.Errorf("Error trying to check if key exists: %s", err.Error())

		return false, err
	}
	if exists == 1 {
		return true, nil
	}

	return false, nil
}

// RemoveFromList delete an value from a list idetinfied with the keyName.
func (r *RedisCluster) RemoveFromList(keyName, value string) error {
	fixedKey := r.fixKey(keyName)

	if err := r.singleton().LRem(fixedKey, 0, value).Err(); err != nil {
		return err
	}

	return nil
}

// GetListRange gets range of elements of list identified by keyName.
func (r *RedisCluster) GetListRange(keyName string, from, to int64) ([]string, error) {
	fixedKey := r.fixKey(keyName)

	elements, err := r.singleton().LRange(fixedKey, from, to).Result()
	if err != nil {
		return nil, err
	}

	return elements, nil
}

// AppendToSetPipelined append values to redis pipeline.
func (r *RedisCluster) AppendToSetPipelined(key string, values [][]byte) {
	if len(values) == 0 {
		return
	}

	fixedKey := r.fixKey(key)
	if err := r.up(); err != nil {
		fmt.Printf(err.Error())

		return
	}
	client := r.singleton()

	pipe := client.Pipeline()
	for _, val := range values {
		pipe.RPush(fixedKey, val)
	}

	if _, err := pipe.Exec(); err != nil {
		fmt.Errorf("Error trying to append to set keys: %s", err.Error())
	}

	// if we need to set an expiration time
	if storageExpTime := int64(viper.GetDuration("analytics.storage-expiration-time")); storageExpTime != int64(-1) {
		// If there is no expiry on the analytics set, we should set it.
		exp, _ := r.GetExp(key)
		if exp == -1 {
			_ = r.SetExp(key, time.Duration(storageExpTime)*time.Second)
		}
	}
}

// GetSet return key set value.
func (r *RedisCluster) GetSet(keyName string) (map[string]string, error) {
	fmt.Printf("Getting from key set: %s", keyName)
	fmt.Printf("Getting from fixed key set: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		return nil, err
	}
	val, err := r.singleton().SMembers(r.fixKey(keyName)).Result()
	if err != nil {
		fmt.Errorf("Error trying to get key set: %s", err.Error())

		return nil, err
	}

	result := make(map[string]string)
	for i, value := range val {
		result[strconv.Itoa(i)] = value
	}

	return result, nil
}

// AddToSet add value to key set.
func (r *RedisCluster) AddToSet(keyName, value string) {
	fmt.Printf("Pushing to raw key set: %s", keyName)
	fmt.Printf("Pushing to fixed key set: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		return
	}
	err := r.singleton().SAdd(r.fixKey(keyName), value).Err()
	if err != nil {
		fmt.Errorf("Error trying to append keys: %s", err.Error())
	}
}

// RemoveFromSet remove a value from key set.
func (r *RedisCluster) RemoveFromSet(keyName, value string) {
	fmt.Printf("Removing from raw key set: %s", keyName)
	fmt.Printf("Removing from fixed key set: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		fmt.Printf(err.Error())

		return
	}
	err := r.singleton().SRem(r.fixKey(keyName), value).Err()
	if err != nil {
		fmt.Errorf("Error trying to remove keys: %s", err.Error())
	}
}

// IsMemberOfSet return whether the given value belong to key set.
func (r *RedisCluster) IsMemberOfSet(keyName, value string) bool {
	if err := r.up(); err != nil {
		fmt.Printf(err.Error())

		return false
	}
	val, err := r.singleton().SIsMember(r.fixKey(keyName), value).Result()
	if err != nil {
		fmt.Errorf("Error trying to check set member: %s", err.Error())

		return false
	}

	fmt.Printf("SISMEMBER %s %s %v %v", keyName, value, val, err)

	return val
}

// SetRollingWindow will append to a sorted set in redis and extract a timed window of values.
func (r *RedisCluster) SetRollingWindow(
	keyName string,
	per int64,
	valueOverride string,
	pipeline bool,
) (int, []interface{}) {
	fmt.Printf("Incrementing raw key: %s", keyName)
	if err := r.up(); err != nil {
		fmt.Printf(err.Error())

		return 0, nil
	}
	fmt.Printf("keyName is: %s", keyName)
	now := time.Now()
	fmt.Printf("Now is: %v", now)
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)
	fmt.Printf("Then is: %v", onePeriodAgo)

	client := r.singleton()
	var zrange *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		zrange = pipe.ZRange(keyName, 0, -1)

		element := redis.Z{
			Score: float64(now.UnixNano()),
		}

		if valueOverride != "-1" {
			element.Member = valueOverride
		} else {
			element.Member = strconv.Itoa(int(now.UnixNano()))
		}

		pipe.ZAdd(keyName, &element)
		pipe.Expire(keyName, time.Duration(per)*time.Second)

		return nil
	}

	var err error
	if pipeline {
		_, err = client.Pipelined(pipeFn)
	} else {
		_, err = client.TxPipelined(pipeFn)
	}

	if err != nil {
		fmt.Errorf("Multi command failed: %s", err.Error())

		return 0, nil
	}

	values := zrange.Val()

	// Check actual value
	if values == nil {
		return 0, nil
	}

	intVal := len(values)
	result := make([]interface{}, len(values))

	for i, v := range values {
		result[i] = v
	}

	fmt.Printf("Returned: %d", intVal)

	return intVal, result
}

// GetRollingWindow return rolling window.
func (r RedisCluster) GetRollingWindow(keyName string, per int64, pipeline bool) (int, []interface{}) {
	if err := r.up(); err != nil {
		fmt.Printf(err.Error())

		return 0, nil
	}
	now := time.Now()
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)

	client := r.singleton()
	var zrange *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		zrange = pipe.ZRange(keyName, 0, -1)

		return nil
	}

	var err error
	if pipeline {
		_, err = client.Pipelined(pipeFn)
	} else {
		_, err = client.TxPipelined(pipeFn)
	}
	if err != nil {
		fmt.Errorf("multi command failed: %s", err.Error())

		return 0, nil
	}

	values := zrange.Val()

	// Check actual value
	if values == nil {
		return 0, nil
	}

	intVal := len(values)
	result := make([]interface{}, intVal)
	for i, v := range values {
		result[i] = v
	}

	fmt.Printf("Returned: %d", intVal)

	return intVal, result
}

// GetKeyPrefix returns storage key prefix.
func (r *RedisCluster) GetKeyPrefix() string {
	return r.KeyPrefix
}

// AddToSortedSet adds value with given score to sorted set identified by keyName.
func (r *RedisCluster) AddToSortedSet(keyName, value string, score float64) {
	fixedKey := r.fixKey(keyName)

	if err := r.up(); err != nil {
		fmt.Printf(err.Error())

		return
	}
	member := redis.Z{Score: score, Member: value}
	if err := r.singleton().ZAdd(fixedKey, &member).Err(); err != nil {
	}
}

// GetSortedSetRange gets range of elements of sorted set identified by keyName.
func (r *RedisCluster) GetSortedSetRange(keyName, scoreFrom, scoreTo string) ([]string, []float64, error) {
	fixedKey := r.fixKey(keyName)

	args := redis.ZRangeBy{Min: scoreFrom, Max: scoreTo}
	values, err := r.singleton().ZRangeByScoreWithScores(fixedKey, &args).Result()
	if err != nil {
		return nil, nil, err
	}

	if len(values) == 0 {
		return nil, nil, nil
	}

	elements := make([]string, len(values))
	scores := make([]float64, len(values))

	for i, v := range values {
		elements[i] = fmt.Sprint(v.Member)
		scores[i] = v.Score
	}
	return elements, scores, nil
}

// RemoveSortedSetRange removes range of elements from sorted set identified by keyName.
func (r *RedisCluster) RemoveSortedSetRange(keyName, scoreFrom, scoreTo string) error {
	fixedKey := r.fixKey(keyName)
	if err := r.singleton().ZRemRangeByScore(fixedKey, scoreFrom, scoreTo).Err(); err != nil {
		return err
	}

	return nil
}
