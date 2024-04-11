package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/services/fastchat"
	"github.com/IceBearAI/aigc/src/services/ldapcli"
	runtime2 "github.com/IceBearAI/aigc/src/services/runtime"
	"github.com/IceBearAI/aigc/src/util"
	"github.com/sashabaranov/go-openai"
	gormlogger "gorm.io/gorm/logger"
	"net"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	kittracing "github.com/go-kit/kit/tracing/opentracing"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"

	"github.com/IceBearAI/aigc/src/encode"
	"github.com/IceBearAI/aigc/src/logging"
	"github.com/IceBearAI/aigc/src/repository"
)

const (
	DefaultHttpPort   = ":8080"
	DefaultConfigPath = "/usr/local/aigc-server/etc/app.cfg"
	DefaultWebPath    = "web"

	// [DB相关]
	EnvNameDbDrive       = "AIGC_DB_DRIVE"
	EnvNameMysqlHost     = "AIGC_MYSQL_HOST"
	EnvNameMysqlPort     = "AIGC_MYSQL_PORT"
	EnvNameMysqlUser     = "AIGC_MYSQL_USER"
	EnvNameMysqlPassword = "AIGC_MYSQL_PASSWORD"
	EnvNameMysqlDatabase = "AIGC_MYSQL_DATABASE"
	//EnvNameRedisHosts    = "AIGC_REDIS_HOSTS"
	//EnvNameRedisDb       = "AIGC_REDIS_DB"
	//EnvNameRedisPassword = "AIGC_REDIS_PASSWORD"
	//EnvNameRedisPrefix   = "AIGC_REDIS_PREFIX"

	// [跨域]
	EnvNameEnableCORS           = "AIGC_ENABLE_CORS"
	EnvNameCORSAllowMethods     = "AIGC_CORS_ALLOW_METHODS"
	EnvNameCORSAllowHeaders     = "AIGC_CORS_ALLOW_HEADERS"
	EnvNameCORSAllowCredentials = "AIGC_CORS_ALLOW_CREDENTIALS"
	EnvNameCORSAllowOrigins     = "AIGC_CORS_ALLOW_ORIGINS"

	// [Trace相关]
	EnvNameTracerEnable         = "AIGC_TRACER_ENABLE"
	EnvNameTracerDrive          = "AIGC_TRACER_DRIVE"
	EnvNameTracerJaegerHost     = "AIGC_TRACER_JAEGER_HOST"
	EnvNameTracerJaegerParam    = "AIGC_TRACER_JAEGER_PARAM"
	EnvNameTracerJaegerType     = "AIGC_TRACER_JAEGER_TYPE"
	EnvNameTracerJaegerLogSpans = "AIGC_TRACER_JAEGER_LOG_SPANS"

	// [外部Service相关]
	EnvNameServiceAlarmHost    = "AIGC_SERVICE_ALARM_HOST"     // 告警相关
	EnvNameServiceLocalAIHost  = "AIGC_SERVICE_CHAT_API_HOST"  // chat-api 地址
	EnvNameServiceLocalAIToken = "AIGC_SERVICE_CHAT_API_TOKEN" // chat-api token
	EnvNameServiceOpenAiEnable = "AIGC_SERVICE_OPENAI_ENABLE"  // openai相关
	EnvNameServiceOpenAiHost   = "AIGC_SERVICE_OPENAI_HOST"
	EnvNameServiceOpenAiToken  = "AIGC_SERVICE_OPENAI_TOKEN"
	EnvNameServiceOpenAiModel  = "AIGC_SERVICE_OPENAI_MODEL"
	EnvNameServiceOpenAiOrgId  = "AIGC_SERVICE_OPENAI_ORG_ID"
	//EnvNameServiceS3Host         = "AIGC_SERVICE_S3_HOST" // S3对象存储相当
	//EnvNameServiceS3AccessKey    = "AIGC_SERVICE_S3_ACCESS_KEY"
	//EnvNameServiceS3SecretKey    = "AIGC_SERVICE_S3_SECRET_KEY"
	//EnvNameServiceS3S3Url        = "AIGC_SERVICE_S3_S3URL"
	//EnvNameServiceS3Region       = "AIGC_SERVICE_S3_REGION"
	//EnvNameServiceS3Bucket       = "AIGC_SERVICE_S3_BUCKET"
	//EnvNameServiceS3BucketPublic = "AIGC_SERVICE_S3_BUCKET_PUBLIC"
	//EnvNameServiceS3DownloadUrl  = "AIGC_SERVICE_S3_DOWNLOAD_URL"
	//EnvNameServiceS3ProjectName  = "AIGC_SERVICE_S3_PROJECT_NAME"
	//EnvNameServiceS3Cluster      = "AIGC_SERVICE_S3_CLUSTER"

	// [LDAP 相关]
	EnvNameLdapHost        = "AIGC_LDAP_HOST"
	EnvNameLdapPort        = "AIGC_LDAP_PORT"
	EnvNameLdapUseSSL      = "AIGC_LDAP_USE_SSL"
	EnvNameLdapBaseDN      = "AIGC_LDAP_BASE_DN"
	EnvNameLdapBindUser    = "AIGC_LDAP_BIND_USER"
	EnvNameLdapBindPass    = "AIGC_LDAP_BIND_PASS"
	EnvNameLdapUserFilter  = "AIGC_LDAP_USER_FILTER"
	EnvNameLdapGroupFilter = "AIGC_LDAP_GROUP_FILTER"
	EnvNameLdapUserAttr    = "AIGC_LDAP_USER_ATTR"

	// [以下是aigc-server模块配置]
	EnvHttpPort                 = "AIGC_ADMIN_HTTP_PORT"
	EnvNameServerLogDrive       = "AIGC_ADMIN_SERVER_LOG_DRIVE"
	EnvNameServerLogPath        = "AIGC_ADMIN_SERVER_LOG_PATH"
	EnvNameServerName           = "AIGC_ADMIN_SERVER_NAME"
	EnvNameServerDebug          = "AIGC_ADMIN_SERVER_DEBUG"
	EnvNameServerKey            = "AIGC_ADMIN_SERVER_KEY"
	EnvNameServerLogLevel       = "AIGC_ADMIN_SERVER_LOG_LEVEL"
	EnvNameServerLogName        = "AIGC_ADMIN_SERVER_LOG_NAME"
	EnvNameServerAigcChannelKey = "AIGC_ADMIN_SERVER_AIGC_CHANNEL_KEY"
	EnvNameServerAdminUser      = "AIGC_ADMIN_SERVER_ADMIN_USER"
	EnvNameServerAdminPass      = "AIGC_ADMIN_SERVER_ADMIN_PASS"
	EnvNameServerStoragePath    = "AIGC_ADMIN_SERVER_STORAGE_PATH"
	EnvNameServerDomain         = "AIGC_ADMIN_SERVER_DOMAIN"

	// [datasets]
	EnvNameDatasetsImage         = "AIGC_DATASETS_IMAGE"
	EnvNameDatasetsModelName     = "AIGC_DATASETS_MODEL_NAME"
	EnvNameDatasetsDevice        = "AIGC_DATASETS_DEVICE"
	EnvNameDatasetsGpuToleration = "AIGC_DATASETS_GPU_TOLERATION"

	// [runtime]
	EnvNameRuntimePlatform        = "AIGC_RUNTIME_PLATFORM"
	EnvNameRuntimeShmSize         = "AIGC_RUNTIME_SHM_SIZE"
	EnvNameRuntimeK8sHost         = "AIGC_RUNTIME_K8S_HOST"
	EnvNameRuntimeK8sToken        = "AIGC_RUNTIME_K8S_TOKEN"
	EnvNameRuntimeK8sConfigPath   = "AIGC_RUNTIME_K8S_CONFIG_PATH"
	EnvNameRuntimeK8sNamespace    = "AIGC_RUNTIME_K8S_NAMESPACE"
	EnvNameRuntimeK8sInsecure     = "AIGC_RUNTIME_K8S_INSECURE"
	EnvNameRuntimeK8sVolumeName   = "AIGC_RUNTIME_K8S_VOLUME_NAME"
	EnvNameRuntimeDockerWorkspace = "AIGC_RUNTIME_DOCKER_WORKSPACE"
	EnvNameRuntimeGpuNum          = "AIGC_RUNTIME_GPU_NUM"

	// [local]
	EnvNameStorageType = "AIGC_STORAGE_TYPE"
	//EnvNameLocalDataPath = "AIGC_LOCAL_DATA_PATH"

	// [fschat]
	EnvNameFsChatControllerAddress = "AIGC_FSCHAT_CONTROLLER_ADDRESS"
	EnvNameFsChatApiAddress        = "AIGC_FSCHAT_API_ADDRESS"

	DefaultRuntimePlatform      = "docker"
	DefaultRuntimeShmSize       = "16G"
	DefaultRuntimeK8sHost       = ""
	DefaultRuntimeK8sToken      = ""
	DefaultRuntimeK8sInsecure   = false
	DefaultRuntimeK8sConfigPath = ""
	DefaultRuntimeK8sNamespace  = "default"
	DefaultRuntimeK8sVolumeName = ""

	// [cronjob]
	EnvNameCronJobAuto = "AIGC_CRONJOB_AUTO"

	DefaultDbDrive       = "sqlite"
	DefaultMysqlHost     = "mysql"
	DefaultMysqlPort     = 3306
	DefaultMysqlUser     = "aigc"
	DefaultMysqlPassword = ""
	DefaultMysqlDatabase = "aigc"
	DefaultRedisHosts    = "redis:6379"
	DefaultRedisDb       = 0
	DefaultRedisPassword = ""
	DefaultRedisPrefix   = "aigc"

	DefaultServerName      = "aigc-server"
	DefaultServerKey       = "Aigcfj@202401"
	DefaultServerLogLevel  = "all"
	DefaultServerLogDrive  = "term"
	DefaultServerLogPath   = ""
	DefaultServerLogName   = "aigc-server.log"
	DefaultServerDebug     = false
	DefaultEnableCORS      = false
	DefaultServerAdminUser = "admin"
	DefaultServerAdminPass = "admin"

	DefaultCORSAllowOrigins     = "*"
	DefaultCORSAllowMethods     = "GET,POST,PUT,DELETE,OPTIONS"
	DefaultCORSAllowHeaders     = "Accept,Content-Type,Content-Length,Accept-Encoding,X-CSRF-Token,Authorization"
	DefaultCORSAllowCredentials = true
	DefaultCORSExposeHeaders    = "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type"

	DefaultJaegerEnable           = false
	DefaultJaegerDrive            = "jaeger"
	DefaultJaegerHost             = "jaeger:6832"
	DefaultJaegerParam    float64 = 1
	DefaultJaegerType             = "const"
	DefaultJaegerLogSpans         = false

	// [chat]相关
	DefaultServiceChatApiHost  = "http://fschat-api:8000/v1"
	DefaultServiceChatApiToken = "sk-001"
	DefaultServiceOpenAiHost   = "https://api.openai.com/v1"
	DefaultServiceOpenAiToken  = "sk-001"
	DefaultServiceOpenAiModel  = openai.GPT3Dot5Turbo
	DefaultServiceOpenAiOrgId  = ""

	// [ldap相关]
	DefaultLdapHost        = "ldap://ldap"
	DefaultLdapPort        = 389
	DefaultLdapBaseDn      = "OU=HABROOT,DC=ORG,DC=corp"
	DefaultLdapBindUser    = "aigc_ldap"
	DefaultLdapBindPass    = ""
	DefaultLdapUserFilter  = "(userPrincipalName=%s)"
	DefaultLdapGroupFilter = ""
	DefaultLdapAttributes  = "name,mail,userPrincipalName,displayName,sAMAccountName"

	// [s3]
	//DefaultServiceS3Host         = "http://s3"
	//DefaultServiceS3AccessKey    = ""
	//DefaultServiceS3SecretKey    = ""
	//DefaultServiceS3Bucket       = "aigc"
	//DefaultServiceS3BucketPublic = "aigc"
	//DefaultServiceS3Region       = "default"
	//DefaultServiceS3Cluster      = "ceph-c2"

	// [datasets]
	DefaultDatasetsImage     = "dudulu/llmops:latest"
	DefaultDatasetsModelName = "uer/sbert-base-chinese-nli"
	DefaultDatasetsDevice    = ""
)

var (
	httpAddr, configPath string
	webPath              string
	logger               log.Logger
	gormDB               *gorm.DB
	db                   *sql.DB
	err                  error
	store                repository.Repository
	namespace            string
	webEmbed             bool

	rootCmd = &cobra.Command{
		Use:               "aigc-server",
		Short:             "",
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Long: `# Aigc Admin服务
有关本系统的相关概述，请参阅 http://github.com/IceBearAI/aigc
`, Version: version,
	}
)

var (
	//rdb    redis.UniversalClient
	apiSvc services.Service
	//hashId   hashids.HashIds
	dbDrive, mysqlHost, mysqlUser, mysqlPassword, mysqlDatabase                                            string
	mysqlPort, redisDb, ormPort                                                                            int
	redisAuth, redisHosts, redisPrefix                                                                     string
	serverName, serverKey, serverLogLevel, serverLogDrive, serverLogPath, serverLogName, serverStoragePath string
	defaultStoragePath, serverDomain                                                                       string
	serverAdminUser, serverAdminPass                                                                       string
	corsAllowOrigins, corsAllowMethods, corsAllowHeaders, corsExposeHeaders                                string
	serverDebug, enableCORS, corsAllowCredentials, tracerEnable, tracerJaegerLogSpans, mysqlOrmMetrics     bool
	tracerDrive, tracerJaegerHost, tracerJaegerType                                                        string
	tracerJaegerParam                                                                                      float64
	serverChannelKey                                                                                       string

	// [gpt]
	serviceOpenAiEnable                                                           bool
	serviceLocalAiHost, serviceLocalAiToken                                       string
	serviceOpenAiHost, serviceOpenAiToken, serviceOpenAiModel, serviceOpenAiOrgId string

	// [s3]
	//serviceS3Host, serviceS3AccessKey, serviceS3SecretKey, serviceS3Bucket, serviceS3BucketPublic, serviceS3Region, serviceS3ProjectName string

	// [ldap]相关
	ldapHost, ldapBaseDn, ldapBindUser, ldapBindPass, ldapUserFilter, ldapGroupFilter string
	ldapPort                                                                          int
	ldapUserAttr                                                                      []string
	ldapUseSsl                                                                        bool

	// [chat]
	defaultServiceChatHost = "http://chat-api:8080"

	// datasets
	datasetsImage, datasetsModelName, datasetsDevice, datasetsGpuToleration string

	// local
	storageType string

	// [runtime]
	runtimePlatform, runtimeShmSize, runtimeK8sHost, runtimeK8sToken, runtimeK8sConfigPath, runtimeK8sNamespace, runtimeK8sVolumeName string
	runtimeDockerWorkspace                                                                                                            string
	runtimeK8sInsecure                                                                                                                bool
	runtimeGpuNum                                                                                                                     int

	// [fschat]
	fsChatControllerAddress, fsChatApiAddress string

	channelId     int
	corsHeaders   = make(map[string]string, 3)
	rateBucketNum = 5000000
	traceId       = logging.TraceId

	// [cronjob]
	cronJobAuto bool

	goOS                                     = runtime.GOOS
	goArch                                   = runtime.GOARCH
	goVersion                                = runtime.Version()
	compiler                                 = runtime.Compiler
	version, buildDate, gitCommit, gitBranch string
)

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
GitCommit: ` + gitCommit + `
GitBranch: ` + gitBranch + `
BuildDate: ` + buildDate + `
Compiler: ` + compiler + `
GoVersion: ` + goVersion + `
Platform: ` + goOS + "/" + goArch + `
`)

	pwd, _ := os.Getwd()

	defaultStoragePath = path.Join(pwd, "storage")

	startCmd.PersistentFlags().StringVarP(&httpAddr, "http.port", "p", DefaultHttpPort, "服务启动的http端口")
	startCmd.PersistentFlags().BoolVar(&webEmbed, "web.embed", true, "是否使用embed.FS")
	startCmd.PersistentFlags().StringVar(&serverDomain, "server.domain", fmt.Sprintf("http://localhost%s", httpAddr), "启动服务的域名")
	// [cors]
	startCmd.PersistentFlags().BoolVar(&enableCORS, "cors.enable", DefaultEnableCORS, "是否开启跨域访问")
	startCmd.PersistentFlags().StringVar(&corsAllowOrigins, "cors.allow.origins", DefaultCORSAllowOrigins, "允许跨域访问的域名")
	startCmd.PersistentFlags().StringVar(&corsAllowMethods, "cors.allow.methods", DefaultCORSAllowMethods, "允许跨域访问的方法")
	startCmd.PersistentFlags().StringVar(&corsAllowHeaders, "cors.allow.headers", DefaultCORSAllowHeaders, "允许跨域访问的头部")
	startCmd.PersistentFlags().StringVar(&corsExposeHeaders, "cors.expose.headers", DefaultCORSExposeHeaders, "允许跨域访问的头部")
	startCmd.PersistentFlags().BoolVar(&corsAllowCredentials, "cors.allow.credentials", DefaultCORSAllowCredentials, "是否允许跨域访问的凭证")
	// [tracer]
	startCmd.PersistentFlags().BoolVar(&tracerEnable, "tracer.enable", DefaultJaegerEnable, "是否启用Tracer")
	startCmd.PersistentFlags().StringVar(&tracerDrive, "tracer.drive", DefaultJaegerDrive, "Tracer驱动")
	startCmd.PersistentFlags().StringVar(&tracerJaegerHost, "tracer.jaeger.host", DefaultJaegerHost, "Tracer Jaeger Host")
	startCmd.PersistentFlags().Float64Var(&tracerJaegerParam, "tracer.jaeger.param", DefaultJaegerParam, "Tracer Jaeger Param")
	startCmd.PersistentFlags().StringVar(&tracerJaegerType, "tracer.jaeger.type", DefaultJaegerType, "采样器的类型 const: 固定采样, probabilistic: 随机取样, ratelimiting: 速度限制取样, remote: 基于Jaeger代理的取样")
	startCmd.PersistentFlags().BoolVar(&tracerJaegerLogSpans, "tracer.jaeger.log.spans", DefaultJaegerLogSpans, "Tracer Jaeger Log Spans")

	// [database]
	rootCmd.PersistentFlags().StringVar(&dbDrive, "db.drive", DefaultDbDrive, "数据库驱动")
	rootCmd.PersistentFlags().StringVar(&mysqlHost, "db.mysql.host", DefaultMysqlHost, "mysql数据库地址: mysql")
	rootCmd.PersistentFlags().IntVar(&mysqlPort, "db.mysql.port", DefaultMysqlPort, "mysql数据库端口")
	rootCmd.PersistentFlags().StringVar(&mysqlUser, "db.mysql.user", DefaultMysqlUser, "mysql数据库用户")
	rootCmd.PersistentFlags().StringVar(&mysqlPassword, "db.mysql.password", DefaultMysqlPassword, "mysql数据库密码")
	rootCmd.PersistentFlags().StringVar(&mysqlDatabase, "db.mysql.database", DefaultMysqlDatabase, "mysql数据库")
	rootCmd.PersistentFlags().BoolVar(&mysqlOrmMetrics, "db.mysql.metrics", false, "是否启GORM的Metrics")
	// [redis]
	//rootCmd.PersistentFlags().StringVar(&redisHosts, "redis.hosts", DefaultRedisHosts, "连接Redis地址")
	//rootCmd.PersistentFlags().IntVar(&redisDb, "redis.db", DefaultRedisDb, "连接Redis DB")
	//rootCmd.PersistentFlags().StringVar(&redisAuth, "redis.auth", DefaultRedisPassword, "连接Redis密码")
	//rootCmd.PersistentFlags().StringVar(&redisPrefix, "redis.prefix", DefaultRedisPrefix, "Redis写入Cache的前缀")
	// [server]
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", DefaultRedisPrefix, "命名空间")
	rootCmd.PersistentFlags().StringVarP(&serverName, "server.name", "a", DefaultServerName, "本系统服务名称")
	rootCmd.PersistentFlags().StringVar(&serverKey, "server.key", DefaultServerKey, "本系统服务密钥")
	rootCmd.PersistentFlags().StringVar(&serverLogLevel, "server.log.level", DefaultServerLogLevel, "本系统日志级别")
	rootCmd.PersistentFlags().StringVar(&serverLogDrive, "server.log.drive", DefaultServerLogDrive, "本系统日志驱动, 支持syslog,term")
	rootCmd.PersistentFlags().StringVar(&serverLogPath, "server.log.path", DefaultServerLogPath, "本系统日志路径")
	rootCmd.PersistentFlags().StringVar(&serverLogName, "server.log.name", DefaultServerLogName, "本系统日志名称")
	rootCmd.PersistentFlags().BoolVar(&serverDebug, "server.debug", DefaultServerDebug, "是否开启Debug模式")
	rootCmd.PersistentFlags().StringVar(&serverAdminUser, "server.admin.user", DefaultServerAdminUser, "系统管理员账号")
	rootCmd.PersistentFlags().StringVar(&serverAdminPass, "server.admin.pass", DefaultServerAdminPass, "系统管理员密码")
	rootCmd.PersistentFlags().StringVar(&serverStoragePath, "server.storage.path", defaultStoragePath, "文件存储绝对路径")
	// [service]
	rootCmd.PersistentFlags().StringVarP(&configPath, "config.path", "c", "", "配置文件路径，如果没有传入配置文件路径则默认使用环境变量")
	//rootCmd.PersistentFlags().StringVar(&serviceAlarmHost, "service.alarm.token", DefaultServiceAlarmHost, "告警中心服务地址")
	// [gpt]
	rootCmd.PersistentFlags().StringVar(&serviceLocalAiHost, "service.local.ai.host", DefaultServiceChatApiHost, "Chat-Api 地址")
	rootCmd.PersistentFlags().StringVar(&serviceLocalAiToken, "service.local.ai.token", DefaultServiceChatApiToken, "Chat-Api Token")
	rootCmd.PersistentFlags().BoolVar(&serviceOpenAiEnable, "service.openai.enable", false, "是否启用OpenAI服务")
	rootCmd.PersistentFlags().StringVar(&serviceOpenAiHost, "service.openai.host", DefaultServiceOpenAiHost, "OpenAI服务地址")
	rootCmd.PersistentFlags().StringVar(&serviceOpenAiModel, "service.openai.model", DefaultServiceOpenAiModel, "OpenAI模型名称")
	rootCmd.PersistentFlags().StringVar(&serviceOpenAiToken, "service.openai.token", "", "OpenAI Token")
	rootCmd.PersistentFlags().StringVar(&serviceOpenAiOrgId, "service.openai.org.id", DefaultServiceOpenAiOrgId, "OpenAI OrgId")
	rootCmd.PersistentFlags().StringVar(&fsChatControllerAddress, "service.fschat.controller.host", "http://fschat-controller:21001", "fastchat controller address")
	rootCmd.PersistentFlags().StringVar(&fsChatApiAddress, "service.fschat.api.host", "http://fschat-api:8000", "fastchat api address")

	// [s3]
	//rootCmd.PersistentFlags().StringVar(&serviceS3Host, "service.s3.host", DefaultServiceS3Host, "S3服务地址")
	//rootCmd.PersistentFlags().StringVar(&serviceS3AccessKey, "service.s3.access.key", DefaultServiceS3AccessKey, "S3 AccessKey")
	//rootCmd.PersistentFlags().StringVar(&serviceS3SecretKey, "service.s3.secret.key", DefaultServiceS3SecretKey, "S3 SecretKey")
	//rootCmd.PersistentFlags().StringVar(&serviceS3Bucket, "service.s3.bucket", DefaultServiceS3Bucket, "S3 Bucket")
	//rootCmd.PersistentFlags().StringVar(&serviceS3BucketPublic, "service.s3.bucket.public", DefaultServiceS3BucketPublic, "S3 Bucket Public")
	//rootCmd.PersistentFlags().StringVar(&serviceS3Region, "service.s3.region", DefaultServiceS3Region, "S3 Bucket")
	//rootCmd.PersistentFlags().StringVar(&serviceS3Cluster, "service.s3.cluster", DefaultServiceS3Cluster, "S3 集群")
	//rootCmd.PersistentFlags().StringVar(&serviceS3ProjectName, "service.s3.project.name", namespace, "S3 项目名称")

	// [ldap]
	startCmd.PersistentFlags().StringVar(&ldapHost, "ldap.host", DefaultLdapHost, "LDAP地址")
	startCmd.PersistentFlags().IntVar(&ldapPort, "ldap.port", DefaultLdapPort, "LDAP端口")
	startCmd.PersistentFlags().StringVar(&ldapBaseDn, "ldap.base.dn", DefaultLdapBaseDn, "LDAP Base DN")
	startCmd.PersistentFlags().BoolVar(&ldapUseSsl, "ldap.use.ssl", false, "LDAP Base DN")
	startCmd.PersistentFlags().StringVar(&ldapBindUser, "ldap.bind.user", DefaultLdapBindUser, "LDAP Bind User")
	startCmd.PersistentFlags().StringVar(&ldapBindPass, "ldap.bind.pass", DefaultLdapBindPass, "LDAP Bind Password")
	startCmd.PersistentFlags().StringVar(&ldapUserFilter, "ldap.user.filter", DefaultLdapUserFilter, "LDAP User Filter")
	startCmd.PersistentFlags().StringVar(&ldapGroupFilter, "ldap.group.filter", DefaultLdapGroupFilter, "LDAP Group Filter")
	startCmd.PersistentFlags().StringSliceVar(&ldapUserAttr, "ldap.user.attr", []string{"name", "mail", "userPrincipalName", "displayName", "sAMAccountName"}, "LDAP Attributes")

	// [runtime]
	rootCmd.PersistentFlags().StringVar(&runtimePlatform, "runtime.platform", DefaultRuntimePlatform, "运行时平台")
	rootCmd.PersistentFlags().StringVar(&runtimeShmSize, "runtime.shm.size", DefaultRuntimeShmSize, "运行时共享内存大小")
	rootCmd.PersistentFlags().StringVar(&runtimeK8sHost, "runtime.k8s.host", DefaultRuntimeK8sHost, "K8s地址")
	rootCmd.PersistentFlags().StringVar(&runtimeK8sToken, "runtime.k8s.token", DefaultRuntimeK8sToken, "K8s Token")
	rootCmd.PersistentFlags().StringVar(&runtimeK8sConfigPath, "runtime.k8s.config.path", DefaultRuntimeK8sConfigPath, "K8s配置文件路径")
	rootCmd.PersistentFlags().StringVar(&runtimeK8sNamespace, "runtime.k8s.namespace", DefaultRuntimeK8sNamespace, "K8s命名空间")
	rootCmd.PersistentFlags().StringVar(&runtimeK8sVolumeName, "runtime.k8s.volume.name", DefaultRuntimeK8sVolumeName, "K8s挂载的存储名")
	rootCmd.PersistentFlags().BoolVar(&runtimeK8sInsecure, "runtime.k8s.insecure", DefaultRuntimeK8sInsecure, "K8s是否不安全")
	rootCmd.PersistentFlags().StringVar(&runtimeDockerWorkspace, "runtime.docker.workspace", defaultStoragePath, "Docker工作目录")
	rootCmd.PersistentFlags().IntVar(&runtimeGpuNum, "runtime.gpu.num", 8, "GPU数量")

	// [dataset]
	startCmd.PersistentFlags().StringVar(&datasetsImage, "datasets.image", DefaultDatasetsImage, "datasets image")
	startCmd.PersistentFlags().StringVar(&datasetsModelName, "datasets.model.name", DefaultDatasetsModelName, "datasets model name")
	startCmd.PersistentFlags().StringVar(&datasetsDevice, "datasets.device", DefaultDatasetsDevice, "datasets device")
	startCmd.PersistentFlags().StringVar(&datasetsGpuToleration, "datasets.gpu.toleration", "", "datasets gpu toleration")

	// [local]
	startCmd.PersistentFlags().StringVar(&storageType, "storage.type", "local", "storage type")

	startCmd.PersistentFlags().BoolVar(&cronJobAuto, "cronjob.auto", true, "是否自动执行定时任务")
	cronJobStartCmd.PersistentFlags().BoolVar(&cronJobAuto, "cronjob.auto", true, "是否自动执行定时任务")

	//jobClearCmd.AddCommand(jobClearAudioTaggedCmd)
	generateCmd.AddCommand(genTableCmd)

	jobFineTuningCmd.AddCommand(jobFineTuningJobRunWaitingTrainCmd, jobFineTuningJobRunningJobLogCmd)
	jobCmd.AddCommand(jobFineTuningCmd)

	cronJobCmd.AddCommand(cronJobStartCmd)

	addFlags(rootCmd)
	rootCmd.AddCommand(startCmd, generateCmd, jobCmd, cronJobCmd)

}

func prepare(ctx context.Context) error {
	logger = log.NewLogfmtLogger(os.Stdout)
	initConfigFile(configPath)

	logger = logging.SetLogging(logger, serverLogPath, serverLogName, serverLogLevel, serverName, serverLogDrive)

	// 连接数据库
	var dbErr error
	if strings.EqualFold(dbDrive, "mysql") {
		dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=20m&collation=utf8mb4_unicode_ci",
			mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
		sqlDB, err := sql.Open("mysql", dbUrl)
		if err != nil {
			_ = level.Error(logger).Log("sql", "Open", "err", err.Error())
			return err
		}
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn:              sqlDB,
			DefaultStringSize: 255,
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if dbErr != nil {
			_ = level.Error(logger).Log("db", "connect", "err", dbErr.Error())
			dbErr = encode.ErrServerStartDbConnect.Wrap(dbErr)
			return dbErr
		}
		_ = level.Debug(logger).Log("mysql", "connect", "success", true)
		//gormDB.Statement.Clauses["soft_delete_enabled"] = clause.Clause{}
	} else if strings.EqualFold(dbDrive, "sqlite") {
		_ = os.MkdirAll(fmt.Sprintf("%s/database", serverStoragePath), 0755)
		gormDB, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/database/aigc.db", serverStoragePath)), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			_ = level.Error(logger).Log("sqlite", "connect", "err", err.Error())
			return err
		}
		_ = level.Debug(logger).Log("sqlite", "connect", "success", true)
	} else {
		err = fmt.Errorf("db drive not support: %s", dbDrive)
		_ = level.Error(logger).Log("db", "drive", "err", err.Error())
		return err
	}
	db, dbErr = gormDB.DB()
	if dbErr != nil {
		_ = level.Error(logger).Log("gormDB", "DB", "err", dbErr.Error())
		dbErr = encode.ErrServerStartDbConnect.Wrap(dbErr)
		return dbErr
	}
	if !strings.EqualFold(serverLogPath, "") {
		gormDB.Logger = logging.NewGormLogging(logger)
	} else {
		gormDB.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxIdleConns(20)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxLifetime(time.Hour)

	if mysqlOrmMetrics {
		//if err = gormDB.Use(gormprometheus.New(gormprometheus.Config{
		//	DBName:          mysqlDatabase,
		//	RefreshInterval: 15,
		//	//PushAddr:        prometheusHost,  // 如果配置了 `PushAddr`，则推送指标
		//	StartServer: false, // 启用一个 http 服务来暴露指标
		//	//HTTPServerPort: uint32(ormPort), // 配置 http 服务监听端口，默认端口为 8080 （如果您配置了多个，只有第一个 `HTTPServerPort` 会被使用）
		//	MetricsCollector: []gormprometheus.MetricsCollector{
		//		&gormprometheus.MySQL{
		//			VariableNames: []string{"Threads_running"},
		//		},
		//	}, // 用户自定义指标
		//})); err != nil {
		//	_ = level.Error(logger).Log("gormDB", "Use", "plugin", "prometheus", "err", err.Error())
		//}
	}

	// 链路追踪
	if tracerEnable {
		tracer, _, err = newJaegerTracer()
		if err != nil {
			_ = level.Warn(logger).Log("jaegerTracer", "connect", "err", err.Error())
		}
		if tracer != nil {
			gormTracingErr := gormDB.Use(gormopentracing.New(gormopentracing.WithTracer(tracer)))
			if gormTracingErr != nil {
				_ = level.Warn(logger).Log("gormDB", "Use", "err", gormTracingErr)
			}
		}
	}

	// 实例化redis
	//rdb, err = redisclient.NewRedisClient(redisHosts, redisAuth, redisPrefix, redisDb, tracer)
	//if err != nil {
	//	_ = level.Error(logger).Log("redis", "connect", "err", err.Error())
	//	return err
	//}
	//_ = level.Debug(logger).Log("redis", "connect", "success", true)

	//go func() {
	//	for {
	//		if pingErr := rdb.Ping(ctx).Err(); pingErr != nil {
	//			_ = level.Error(logger).Log("redis", "ping", "err", pingErr)
	//		}
	//		time.Sleep(time.Second)
	//	}
	//}()

	var clientOpts []kithttp.ClientOption

	dialer := &net.Dialer{
		Timeout:   10 * time.Minute,
		KeepAlive: 10 * time.Minute,
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Minute,
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	clientOpts = []kithttp.ClientOption{
		kithttp.SetClient(httpClient),
		kithttp.ClientBefore(kittracing.ContextToHTTP(tracer, logger)),
	}

	// 实例化仓库
	store = repository.New(gormDB, logger, traceId, tracer)

	var runtimeOpts []runtime2.CreationOption
	runtimeOpts = append(runtimeOpts,
		runtime2.WithHttpClientOptions(clientOpts),
		runtime2.WithK8sConfigPath(runtimeK8sConfigPath),
		runtime2.WithK8sToken(runtimeK8sHost, runtimeK8sToken, runtimeK8sInsecure),
		runtime2.WithShmSize(runtimeShmSize),
		runtime2.WithK8sVolumeName(runtimeK8sVolumeName),
		runtime2.WithNamespace(runtimeK8sNamespace),
		runtime2.WithWorkspace(runtimeDockerWorkspace),
		runtime2.WithGpuNum(runtimeGpuNum),
	)
	apiSvc = services.NewApi(ctx, logger, traceId, serverDebug, tracer, &services.Config{
		Namespace: namespace, ServiceName: serverName,
		FastChat: fastchat.Config{
			OpenAiEndpoint:  serviceOpenAiHost,
			OpenAiToken:     serviceOpenAiToken,
			OpenAiModel:     serviceOpenAiModel,
			OpenAiOrgId:     serviceOpenAiOrgId,
			LocalAiEndpoint: serviceLocalAiHost,
			LocalAiToken:    serviceLocalAiToken,
			Debug:           serverDebug,
		},
		Ldap: ldapcli.Config{
			Host:         ldapHost,
			Port:         ldapPort,
			UseSSL:       ldapUseSsl,
			BindUser:     ldapBindUser,
			BindPassword: ldapBindPass,
			BindDN:       ldapBaseDn,
			Attributes:   ldapUserAttr,
			Filter:       ldapUserFilter,
		},
		Runtime:         runtimeOpts,
		RuntimePlatform: runtimePlatform,
	}, clientOpts)

	// 如果是docker来台，查询fschat-controller 和 fschat-api是否启动，如果没有则创建
	if strings.EqualFold(runtimePlatform, "docker") {
		if fschatErr := runFastChat(ctx); fschatErr != nil {
			_ = level.Warn(logger).Log("fschat", "run", "err", fschatErr)
		}
	}

	return err
}

func runFastChat(ctx context.Context) (err error) {
	status, err := apiSvc.Runtime().GetDeploymentStatus(ctx, "fschat-controller")
	if err != nil {
		_ = level.Error(logger).Log("fschat-controller", "status", "err", err)
		return err
	}
	_ = level.Info(logger).Log("fschat-controller", "status", status)
	if !util.StringInArray([]string{
		"Failed",
		"Unknown",
	}, status) {
		_ = level.Info(logger).Log("fschat-controller", "status", "running", true)
		return
	}

	// 创建fschat-controller
	deploymentName, err := apiSvc.Runtime().CreateDeployment(ctx, runtime2.Config{
		ServiceName: "fschat-controller",
		Image:       "dudulu/fschat:latest",
		Command: []string{
			"python3",
			"-m",
			"fastchat.serve.controller",
			"--host",
			"0.0.0.0",
			"--port",
			"21001",
		},
		Ports: map[string]string{
			"21001": "21001",
		},
		Replicas: 1,
	})
	if err != nil {
		_ = level.Error(logger).Log("fschat-controller", "create", "err", err)
		return err
	}

	_ = level.Info(logger).Log("fschat-controller", "create", "success", deploymentName)
	status, err = apiSvc.Runtime().GetDeploymentStatus(ctx, "fschat-api")
	if err != nil {
		_ = level.Error(logger).Log("fschat-api", "status", "err", err)
		return err
	}
	_ = level.Info(logger).Log("fschat-api", "status", status)
	if !util.StringInArray([]string{
		"Failed",
		"Unknown",
	}, status) {
		_ = level.Info(logger).Log("fschat-api", "status", "running", true)
	}
	// 创建fschat-api
	deploymentName, err = apiSvc.Runtime().CreateDeployment(ctx, runtime2.Config{
		ServiceName: "fschat-api",
		Image:       "dudulu/fschat:latest",
		Command: []string{
			"python3",
			"-m",
			"fastchat.serve.openai_api_server",
			"--host",
			"0.0.0.0",
			"--port",
			"8000",
			"--controller-address",
			"http://$(hostname -I | awk '{print $1}'):21001",
		},
		Ports: map[string]string{
			"8000": "8000",
		},
		Replicas: 1,
	})
	if err != nil {
		_ = level.Error(logger).Log("fschat-api", "create", "err", err)
		return err
	}

	_ = level.Info(logger).Log("fschat-api", "create", "success", deploymentName)
	return
}

func Run() {
	webPath = envString("WEB_PATH", DefaultWebPath)

	httpAddr = envString(EnvHttpPort, DefaultHttpPort)
	namespace = envString("POD_NAMESPACE", envString("NAMESPACE", namespace))

	// [database]
	dbDrive = envString(EnvNameDbDrive, DefaultDbDrive)
	mysqlHost = envString(EnvNameMysqlHost, DefaultMysqlHost)
	mysqlPort, _ = strconv.Atoi(envString(EnvNameMysqlPort, strconv.Itoa(DefaultMysqlPort)))
	mysqlUser = envString(EnvNameMysqlUser, DefaultMysqlUser)
	mysqlPassword = envString(EnvNameMysqlPassword, DefaultMysqlPassword)
	mysqlDatabase = envString(EnvNameMysqlDatabase, DefaultMysqlDatabase)

	// [redis]
	//redisHosts = envString(EnvNameRedisHosts, DefaultRedisHosts)
	//redisDb, _ = strconv.Atoi(envString(EnvNameRedisDb, strconv.Itoa(DefaultRedisDb)))
	//redisAuth = envString(EnvNameRedisPassword, DefaultRedisPassword)
	//redisPrefix = envString(EnvNameRedisPrefix, DefaultRedisPrefix)

	// [cors]
	enableCORS, _ = strconv.ParseBool(envString(EnvNameEnableCORS, strconv.FormatBool(DefaultEnableCORS)))
	corsAllowMethods = envString(EnvNameCORSAllowMethods, DefaultCORSAllowMethods)
	corsAllowHeaders = envString(EnvNameCORSAllowHeaders, DefaultCORSAllowHeaders)
	corsAllowOrigins = envString(EnvNameCORSAllowOrigins, DefaultCORSAllowOrigins)
	corsAllowCredentials, _ = strconv.ParseBool(envString(EnvNameCORSAllowCredentials, strconv.FormatBool(DefaultCORSAllowCredentials)))

	// [trace]
	tracerEnable, _ = strconv.ParseBool(envString(EnvNameTracerEnable, strconv.FormatBool(DefaultJaegerEnable)))
	tracerDrive = envString(EnvNameTracerDrive, DefaultJaegerDrive)
	tracerJaegerParam, _ = strconv.ParseFloat(envString(EnvNameTracerJaegerParam, strconv.FormatFloat(tracerJaegerParam, 'f', -1, 64)), 64)
	tracerJaegerHost = envString(EnvNameTracerJaegerHost, DefaultJaegerHost)
	tracerJaegerType = envString(EnvNameTracerJaegerType, DefaultJaegerType)
	tracerJaegerLogSpans, _ = strconv.ParseBool(envString(EnvNameTracerJaegerLogSpans, strconv.FormatBool(DefaultJaegerLogSpans)))

	// [server]
	serverName = envString(EnvNameServerName, DefaultServerName)
	serverKey = envString(EnvNameServerKey, DefaultServerKey)
	serverLogLevel = envString(EnvNameServerLogLevel, DefaultServerLogLevel)
	serverLogDrive = envString(EnvNameServerLogDrive, DefaultServerLogDrive)
	serverLogPath = envString(EnvNameServerLogPath, DefaultServerLogPath)
	serverLogName = envString(EnvNameServerLogName, DefaultServerLogName)
	serverChannelKey = envString(EnvNameServerAigcChannelKey, "sk-001")
	serverDebug, _ = strconv.ParseBool(envString(EnvNameServerDebug, strconv.FormatBool(DefaultServerDebug)))
	serverAdminUser = envString(EnvNameServerAdminUser, DefaultServerAdminUser)
	serverAdminUser = envString(EnvNameServerAdminPass, DefaultServerAdminPass)
	serverStoragePath = envString(EnvNameServerStoragePath, defaultStoragePath)
	serverDomain = envString(EnvNameServerDomain, fmt.Sprintf("http://localhost%s", httpAddr))
	cronJobAuto, _ = strconv.ParseBool(envString(EnvNameCronJobAuto, "true"))

	// [service.gpt]
	serviceOpenAiEnable, _ = strconv.ParseBool(envString(EnvNameServiceOpenAiEnable, "false"))
	serviceOpenAiHost = envString(EnvNameServiceOpenAiHost, DefaultServiceOpenAiHost)
	serviceOpenAiToken = envString(EnvNameServiceOpenAiToken, DefaultServiceOpenAiToken)
	serviceOpenAiModel = envString(EnvNameServiceOpenAiModel, DefaultServiceOpenAiModel)
	serviceOpenAiOrgId = envString(EnvNameServiceOpenAiOrgId, DefaultServiceOpenAiOrgId)
	serviceLocalAiHost = envString(EnvNameServiceLocalAIHost, DefaultServiceChatApiHost)
	serviceLocalAiToken = envString(EnvNameServiceLocalAIToken, DefaultServiceChatApiToken)

	// [ldap]
	ldapHost = envString(EnvNameLdapHost, DefaultLdapHost)
	ldapPort, _ = strconv.Atoi(envString(EnvNameLdapPort, strconv.Itoa(DefaultLdapPort)))
	ldapUseSsl, _ = strconv.ParseBool(envString(EnvNameLdapUseSSL, "false"))
	ldapBindUser = envString(EnvNameLdapBindUser, DefaultLdapBindUser)
	ldapBindPass = envString(EnvNameLdapBindPass, DefaultLdapBindPass)
	ldapBaseDn = envString(EnvNameLdapBaseDN, DefaultLdapBaseDn)
	ldapUserFilter = envString(EnvNameLdapUserFilter, DefaultLdapUserFilter)
	ldapUserAttr = strings.Split(envString(EnvNameLdapUserAttr, DefaultLdapAttributes), ",")

	// [service.s3]
	//serviceS3Host = envString(EnvNameServiceS3Host, DefaultServiceS3Host)
	//serviceS3AccessKey = envString(EnvNameServiceS3AccessKey, DefaultServiceS3AccessKey)
	//serviceS3SecretKey = envString(EnvNameServiceS3SecretKey, DefaultServiceS3SecretKey)
	//serviceS3Bucket = envString(EnvNameServiceS3Bucket, DefaultServiceS3Bucket)
	//serviceS3BucketPublic = envString(EnvNameServiceS3BucketPublic, DefaultServiceS3BucketPublic)
	//serviceS3Region = envString(EnvNameServiceS3Region, DefaultServiceS3Region)
	//serviceS3ProjectName = envString(EnvNameServiceS3ProjectName, namespace)

	// [dataset]
	datasetsImage = envString(EnvNameDatasetsImage, DefaultDatasetsImage)
	datasetsModelName = envString(EnvNameDatasetsModelName, DefaultDatasetsModelName)
	datasetsDevice = envString(EnvNameDatasetsDevice, DefaultDatasetsDevice)
	datasetsGpuToleration = envString(EnvNameDatasetsGpuToleration, "")

	// [local]
	storageType = envString(EnvNameStorageType, "local")
	//localDataPath = envString(EnvNameLocalDataPath, DefaultLocalDataPath)

	// [runtime]
	runtimePlatform = envString(EnvNameRuntimePlatform, DefaultRuntimePlatform)
	runtimeShmSize = envString(EnvNameRuntimeShmSize, DefaultRuntimeShmSize)
	runtimeK8sHost = envString(EnvNameRuntimeK8sHost, DefaultRuntimeK8sHost)
	runtimeK8sToken = envString(EnvNameRuntimeK8sToken, DefaultRuntimeK8sToken)
	runtimeK8sConfigPath = envString(EnvNameRuntimeK8sConfigPath, DefaultRuntimeK8sConfigPath)
	runtimeK8sNamespace = envString(EnvNameRuntimeK8sNamespace, DefaultRuntimeK8sNamespace)
	runtimeK8sVolumeName = envString(EnvNameRuntimeK8sVolumeName, DefaultRuntimeK8sVolumeName)
	runtimeK8sInsecure, _ = strconv.ParseBool(envString(EnvNameRuntimeK8sInsecure, strconv.FormatBool(DefaultRuntimeK8sInsecure)))
	runtimeDockerWorkspace = envString(EnvNameRuntimeDockerWorkspace, defaultStoragePath)
	runtimeGpuNum, _ = strconv.Atoi(envString(EnvNameRuntimeGpuNum, "8"))

	// [fschat]
	fsChatControllerAddress = envString(EnvNameFsChatControllerAddress, "http://fschat-controller:21001")
	fsChatApiAddress = envString(EnvNameFsChatControllerAddress, "http://fschat-api:8000")

	if err = rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute", err.Error())
		os.Exit(-1)
	}
}

func envString(env, fallback string) string {
	e, _ := os.LookupEnv(env)
	if e == "" {
		return fallback
	}
	return e
}

func getStringDefault(value, defaultValue string) string {
	if strings.EqualFold(value, "") {
		return defaultValue
	}
	return value
}

func addFlags(rootCmd *cobra.Command) {
	flag.CommandLine.VisitAll(func(gf *flag.Flag) {
		rootCmd.PersistentFlags().AddGoFlag(gf)
	})
}

// 如果传入了配置文件说明走配置文件否则使用环境变量
func initConfigFile(configPath string) {
	if strings.EqualFold(configPath, "") {
		_ = level.Debug(logger).Log("config", "usage", "Environment", "环境变量")
		return
	}
	_ = level.Debug(logger).Log("config", "usage", "configPath", configPath)
	//cfg, err = config.NewConfig(configPath)
	//if err != nil {
	//	_ = level.Error(logger).Log("config", "NewConfig", "err", err.Error())
	//	panic(err)
	//}
	//
	//// [database]
	//dbDrive = getStringDefault(cfg.GetString("database", "drive"), dbDrive)
	//mysqlHost = getStringDefault(cfg.GetString("database", "host"), mysqlHost)
	//mysqlPort, _ = strconv.Atoi(cfg.GetString("database", "port"))
	//mysqlUser = getStringDefault(cfg.GetString("database", "user"), mysqlUser)
	//mysqlPassword = getStringDefault(cfg.GetString("database", "password"), mysqlPassword)
	//mysqlDatabase = getStringDefault(cfg.GetString("database", "database"), mysqlDatabase)
	//
	//// [tracer]
	//tracerEnable, _ = strconv.ParseBool(cfg.GetString("tracer", "enable"))
	//tracerDrive = getStringDefault(cfg.GetString("tracer", "drive"), tracerDrive)
	//tracerJaegerParam, _ = strconv.ParseFloat(cfg.GetString("tracer", "jaeger.param"), 64)
	//tracerJaegerHost = getStringDefault(cfg.GetString("tracer", "jaeger.host"), tracerJaegerHost)
	//tracerJaegerType = getStringDefault(cfg.GetString("tracer", "jaeger.type"), tracerJaegerType)
	//tracerJaegerLogSpans, _ = strconv.ParseBool(cfg.GetString("tracer", "jaeger.logSpans"))
}

func closeConnection(ctx context.Context) {
	if db != nil {
		_ = level.Debug(logger).Log("db", "close", "err", db.Close())
	}
	//if rdb != nil {
	//	_ = level.Debug(logger).Log("rdb", "close", "err", rdb.Close())
	//}
}
