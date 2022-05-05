package mq

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
)

func New(redis *redis.Pool, maxRetries int) MessageQueue {
	return &QueueArgs{
		redis:      redis,
		maxRetries: maxRetries,
	}
}

func (q *QueueArgs) Publish(topic string, msg string, delay int) error {
	_redis := q.redis.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	data := &MsgArgs{
		Msg:   msg,
		Delay: delay,
		Retry: 0,
	}
	if _data, err := json.Marshal(data); err != nil {
		return err
	} else {
		_, err := _redis.Do("LPUSH", "rmq:"+topic, string(_data))
		return err
	}
}

func (q *QueueArgs) Subscribe(topic string, handler func(data string)) {
	go func() {
		_redis := q.redis.Get()
		defer func() {
			// handle error
			if err := recover(); err != nil {
				q.Subscribe(topic, handler)
			}
		}()

		for {
			_data, err := redis.Strings(_redis.Do("BRPOP", "rmq:"+topic, 1))
			if err != nil || _data == nil {
				continue
			}

			var _msg MsgArgs
			if err := json.Unmarshal([]byte(_data[1]), &_msg); err != nil {
				continue
			}
			// execute handler
			go func(_msg MsgArgs) {
				defer func() {
					// retry if error
					if err := recover(); err != nil {
						if _data, err := json.Marshal(_msg); err != nil {
							return
						} else {
							// max retry
							if _msg.Retry > q.maxRetries {
								_, _ = _redis.Do("LPUSH", "rmq:"+topic+"failed", _data)
								return
							}
							_, _ = _redis.Do("LPUSH", "rmq:"+topic, _data)
						}
						_msg.Retry++
					}
				}()
				time.Sleep(time.Duration(_msg.Delay) * time.Second)
				handler(_msg.Msg)
			}(_msg)
		}
	}()
}
