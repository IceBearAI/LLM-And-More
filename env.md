#### 系统公共环境变量配置

可以修改`.env`调整相关配置

##### 数据库配置

目前支持两类数据库的配置，默认是使用**sqlite**,如果是使用的sqlite那么默认会存储在`AIGC_ADMIN_SERVER_STORAGE_PATH`
所配置的路径下的`storage/database/aigc.db`。

如果使用的是`mysql`驱动，则按照下面配置设置。

| 变量名                   | 描述               | 值       |
|-----------------------|------------------|---------|
| `AIGC_DB_DRIVER`      | 数据库驱动类型（可能是遗留错误） | `mysql` |
| `AIGC_MYSQL_DRIVE`    | 数据库驱动类型          | `mysql` |
| `AIGC_MYSQL_HOST`     | 数据库主机地址          | `mysql` |
| `AIGC_MYSQL_PORT`     | 数据库端口号           | `3306`  |
| `AIGC_MYSQL_USER`     | 数据库用户名           | `aigc`  |
| `AIGC_MYSQL_PASSWORD` | 数据库密码            | `admin` |
| `AIGC_MYSQL_DATABASE` | 数据库名             | `aigc`  |

##### Tracer 链路追踪配置

如果想看整个调用链以下是相关配置，默认不开启。

| 变量名                            | 描述          | 值        |
|--------------------------------|-------------|----------|
| `AIGC_TRACER_ENABLE`           | 是否启用链路追踪    | `false`  |
| `AIGC_TRACER_DRIVE`            | 链路追踪驱动类型    | `jaeger` |
| `AIGC_TRACER_JAEGER_HOST`      | Jaeger 服务地址 |          |
| `AIGC_TRACER_JAEGER_PARAM`     | Jaeger 采样参数 | `1`      |
| `AIGC_TRACER_JAEGER_TYPE`      | Jaeger 采样类型 | `const`  |
| `AIGC_TRACER_JAEGER_LOG_SPANS` | 是否记录追踪日志    | `false`  |

##### 跨域配置

跨域配置，默认不开启

| 变量名                           | 描述        | 值                                                                                                   |
|-------------------------------|-----------|-----------------------------------------------------------------------------------------------------|
| `AIGC_ENABLE_CORS`            | 是否启用CORS  | `true`                                                                                              |
| `AIGC_CORS_ALLOW_METHODS`     | 允许的HTTP方法 | `GET,POST,PUT,DELETE,OPTIONS`                                                                       |
| `AIGC_CORS_ALLOW_HEADERS`     | 允许的HTTP头  | `Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization,x-tenant-id,x-token` |
| `AIGC_CORS_ALLOW_CREDENTIALS` | 是否允许携带凭证  | `true`                                                                                              |
| `AIGC_CORS_ALLOW_ORIGINS`     | 允许的源      | `*`                                                                                                 |

##### 外部服务调用配置

如果还有使用OpenAI的相关模型，则设置OpenAI的相关信息。

| 变量名                           | 描述            | 值                           |
|-------------------------------|---------------|-----------------------------|
| `AIGC_SERVICE_CHAT_API_HOST`  | 聊天API服务地址     | `http://localhost:8000/v1`  |
| `AIGC_SERVICE_CHAT_API_TOKEN` | 聊天API服务访问令牌   |                             |
| `AIGC_SERVICE_OPENAI_ORG_ID`  | OpenAI 组织ID   |                             |
| `AIGC_SERVICE_OPENAI_HOST`    | OpenAI 服务地址   | `https://api.openai.com/v1` |
| `AIGC_SERVICE_OPENAI_TOKEN`   | OpenAI 服务访问令牌 |                             |

##### LDAP 配置

如果是企业使用可以配置LDAP地址。

| 变量名                   | 描述          | 值                                                        |
|-----------------------|-------------|----------------------------------------------------------|
| `AIGC_LDAP_HOST`      | LDAP 服务器地址  | `ldap`                                                   |
| `AIGC_LDAP_BASE_DN`   | LDAP 基础DN   | `OU=HABROOT,DC=corp`                                     |
| `AIGC_LDAP_BIND_USER` | LDAP 绑定用户   |                                                          |
| `AIGC_LDAP_BIND_PASS` | LDAP 绑定用户密码 |                                                          |
| `AIGC_LDAP_USER_ATTR` | LDAP 用户属性   | `name,mail,userPrincipalName,displayName,sAMAccountName` |

##### aigc-server 环境变量配置

本系统基础配置，通常不需要修改。

需要注意的是`AIGC_ADMIN_SERVER_ADMIN_USER`和`AIGC_ADMIN_SERVER_ADMIN_PASS`是系统初始化的管理员账号密码，只有在系统第一次启动初始化的时候配置有效，后续调整将不会生效。

| 变量名                                     | 描述                            | 值                       |
|-----------------------------------------|-------------------------------|-------------------------|
| `AIGC_ADMIN_SERVER_HTTP_PORT`           | 服务HTTP端口                      | `:8080`                 |
| `AIGC_ADMIN_SERVER_LOG_DRIVE`           | 日志驱动类型(默认term)                | `term`                  |
| `AIGC_ADMIN_SERVER_NAME`                | 服务名称                          | `aigc-server`           |
| `AIGC_ADMIN_SERVER_DEBUG`               | 是否开启调试模式(开启后控制台显示所有Debug信息)   | `true`                  |
| `AIGC_ADMIN_SERVER_LOG_LEVEL`           | 日志级别(debug,info,warn,error)   | `all`                   |
| `AIGC_ADMIN_SERVER_LOG_PATH`            | 日志路径(设置之后会写入文件)               |                         |
| `AIGC_ADMIN_SERVER_LOG_NAME`            | 日志文件名称                        | `aigc-server.log`       |
| `AIGC_ADMIN_SERVER_DEFAULT_CHANNEL_KEY` | 默认渠道密钥                        | `sk-001`                |
| `AIGC_ADMIN_SERVER_STORAGE_PATH`        | 上传文件所存储的路径                    | `./storage/`            |
| `AIGC_ADMIN_SERVER_DOMAIN`              | 本服务的域名(容器回调传输数据，需要保证容器网络可以访问) | `http://localhost:8080` |
| `AIGC_ADMIN_SERVER_ADMIN_USER`          | 初始化默认账号                       | `admin`                 |
| `AIGC_ADMIN_SERVER_ADMIN_PASS`          | 初始化默认密码                       | `admin`                 |
