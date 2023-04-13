# openid-go

The source code of https://9420.ltd backend

# project structure

```
├── README.md
├── config
│   ├── config.go   # config
│   └── config.yaml # config demo
├── main.go        # entry
├── app
│   ├── controller # controller
│   ├── model      # model
│   ├── middleware # middleware
├── library
│   ├── apiutil   # api format
│   ├── apputil   # app related tools 
│   ├── codeutil  # send verification code
│   ├── mailutil  # send mail
│   ├── mq        # redis based message queue
│   ├── toolutil  # tool like "hash" "randStr" "regex"
│   ├── userutil  # user management
├── process
│   ├── dbutil    # database related tools
│   ├── queueutil # message queue
│   ├── redisutil # redis 
│   ├── webutil   # gin
```

# copyright

Copyright xcsoft 2023