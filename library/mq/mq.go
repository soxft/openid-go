package mq

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

func New(c context.Context, redis *redis.Client, maxRetries int) MessageQueue {
	return &QueueArgs{
		redis:      redis,
		maxRetries: maxRetries,
		ctx:        c,
	}
}

func (q *QueueArgs) Publish(topic string, msg string, delay int64) error {
	_redis := q.redis

	data := &MsgArgs{
		Msg:     msg,
		DelayAt: delay + time.Now().Unix(),
		Retry:   0,
	}
	if _data, err := json.Marshal(data); err != nil {
		return err
	} else {
		return _redis.LPush(q.ctx, "rmq:"+topic, string(_data)).Err()
	}
}

func (q *QueueArgs) Subscribe(topic string, processes int, handler func(data string)) {
	for i := 0; i < processes; i++ {
		go func() {
			_redis := q.redis
			defer func() {
				// handle error
				if err := recover(); err != nil {
					q.Subscribe(topic, 1, handler)
				}
			}()

			// 阻塞
			wg := sync.WaitGroup{}
			for {
				_data, err := _redis.BRPop(q.ctx, 1*time.Second, "rmq:"+topic).Result()
				if err != nil || _data == nil {
					continue
				}

				var _dataString = _data[1]
				var _msg MsgArgs
				if err := json.Unmarshal([]byte(_dataString), &_msg); err != nil {
					continue
				}

				wg.Add(1)
				// execute handler
				go func(_msg MsgArgs) {
					defer func() {
						// retry if error
						wg.Done()
						if err := recover(); err != nil {
							log.Printf("[ERROR] mq handler: %s", err)
							if _data, err := json.Marshal(_msg); err != nil {
								return
							} else {
								// max retry
								if _msg.Retry > q.maxRetries {
									_redis.LPush(q.ctx, "rmq:"+topic+"failed", _data)
									return
								}
								_redis.LPush(q.ctx, "rmq:"+topic, _data)
							}
							_msg.Retry++
						}
					}()
					// delay 重新放入队列
					if _msg.DelayAt > time.Now().Unix() {
						if err := _redis.LPush(q.ctx, "rmq:"+topic, _dataString).Err(); err != nil {
							log.Printf("[ERROR] mq delay lpush: %s", err)
						}
						return
					}
					handler(_msg.Msg)

					// prevent loop
					time.Sleep(time.Millisecond * 500)
				}(_msg)
				// wait for handler
				wg.Wait()
			}
		}()
	}
}
