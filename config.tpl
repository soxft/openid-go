Server:
  Address: 127.0.0.1:8080
  Debug: true
  Log: true
  Title: "X openID"
  ServerName: "X openID"
  Url: http://127.0.0.1:8080
Redis:
  Address: "127.0.0.1:6379"
  Password:
  Database: 0
  Prefix: openid
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
Aliyun:
  AccessKey: "aliyun_access_key"
  AccessSecret: "aliyun_access_secret"
  Email: "aliyun_email"
Jwt:
  Secret: "jwt_secret"
Developer:
  AppLimit: 10
Github:
  ClientID: "github_client_id"
  ClientSecret: "github_client_secret"