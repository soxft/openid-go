Server:
  Address: 127.0.0.1:8080
  Debug: true
  Log: true
  Title: "X openID"
  ServerName: "X openID"
  FrontUrl: http://127.0.0.1:8080
Redis:
  Address: "127.0.0.1:6379"
  Password:
  Database: 0
  Prefix: openid
  MinIdle: 10
  MaxIdle: 50
  MaxActive: 500
  MaxRetries: 3
Mysql:
  Address: 127.0.0.1:3306
  Username: openid
  Password: openid
  Database: openid
  Charset: utf8mb4
  MaxOpen: 200
  MaxIdle: 100
  MaxLifetime: 240
Aliyun: # Aliyun 邮件推送
  Domain: dm.aliyuncs.com
  Region: cn-hangzhou
  Version: 2015-11-23
  AccessKey: AccessKey
  AccessSecret: AccessSecret
  Email: no-reply@mail.example.com
Smtp: # SMTP配置
  Host: smtp.example.com
  Port: 465
  Secure: true
  Username: username
  Password: password
Jwt:
  Secret: "jwt_secret"
Developer:
  AppLimit: 10
Github:
  ClientID: "github_client_id"
  ClientSecret: "github_client_secret"