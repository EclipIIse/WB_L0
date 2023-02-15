package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Client interface {
	Set(key string, value interface{}) (err error)
	WriteKeyValues(pairs ...interface{}) (err error)
	WriteGroup(data map[string]interface{}) (err error)
	ReadGroup(readProperty string) (messages []map[string]interface{}, err error)
	GetValue(key string) (value interface{}, err error)
	MGet(keys []string) (value []interface{}, err error)
	Ack(id string) (err error)
	DeleteKeyValue(keys ...string) (err error)
	IsNotFound(err error) bool
}

type client struct {
	rClient   *redis.Client
	argsWrite *redis.XAddArgs
	argsRead  *redis.XReadGroupArgs
}

func (c *client) WriteKeyValues(pairs ...interface{}) (err error) {
	cmd := c.rClient.MSet(pairs...)
	err = cmd.Err()
	return
}
func (c *client) Set(key string, value interface{}) (err error) {
	cmd := c.rClient.Set(key, value, 0)
	err = cmd.Err()
	return
}

func (c *client) WriteGroup(data map[string]interface{}) (err error) {
	var args = *c.argsWrite
	args.Values = data

	cmd := c.rClient.XAdd(&args)
	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}
	return
}

// ReadGroup ...
func (c *client) ReadGroup(readProperty string) (messages []map[string]interface{}, err error) {
	// копируем аргументы чтения
	var args = *c.argsRead
	if readProperty != "" && readProperty != " " {
		// добавляем поток и параметр чтения
		args.Streams = make([]string, 2)
		args.Streams[0] = c.argsRead.Streams[0]
		args.Streams[1] = readProperty
	}
	cmd := c.rClient.XReadGroup(&args)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		err = cmd.Err()
		return
	}

	for _, value := range cmd.Val() {
		for _, msg := range value.Messages {
			messages = append(messages, msg.Values)
		}

	}
	return
}

// GetValue ...
func (c *client) GetValue(key string) (value interface{}, err error) {
	cmd := c.rClient.Get(key)
	if cmd.Err() != nil && cmd.Err() == redis.Nil {
		err = cmd.Err()
		return
	}
	return cmd.Val(), nil
}

func (c *client) MGet(keys []string) (value []interface{}, err error) {
	cmd := c.rClient.MGet(keys...)
	fmt.Println(cmd.Result())
	if cmd.Err() != nil && cmd.Err() == redis.Nil {
		err = cmd.Err()
		return
	}

	return cmd.Val(), nil
}

func (c *client) Ack(id string) (err error) {
	if len(c.argsRead.Streams) == 0 {
		err = fmt.Errorf("требуется указать поток для чтения")
		return
	}
	// Подтверждаем обработку сообщения
	cmd := c.rClient.XAck(c.argsRead.Group, c.argsRead.Streams[0], id)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return
}

// DeleteKeyValue ...
func (c *client) DeleteKeyValue(keys ...string) (err error) {
	cmd := c.rClient.Del(keys...)
	return cmd.Err()
}

// IsNotFound ...
func (c *client) IsNotFound(err error) bool {
	return err == redis.Nil
}

func New(rClient *redis.Client,

//argsWrite *redis.XAddArgs,
//argsRead *redis.XReadGroupArgs,
) Client {
	return &client{
		rClient: rClient,
		//argsWrite: argsWrite,
		//argsRead:  argsRead,
	}
}
