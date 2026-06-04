package http

import (
	"fmt"
	"net"

	"user/internal/service"

	"github.com/UnderTreeTech/waterdrop/pkg/server/http/config"

	"github.com/UnderTreeTech/waterdrop/pkg/server/http/server"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xnet"

	"github.com/UnderTreeTech/waterdrop/pkg/registry"

	"user/internal/server/http/middlewares/cors"
	"user/internal/server/http/middlewares/sanitizer"
	"user/internal/utils"
	"user/internal/utils/reply"

	"strconv"

	"net/http"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/UnderTreeTech/waterdrop/pkg/trace"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerInfo struct {
	Server      *server.Server
	ServiceInfo *registry.ServiceInfo
}

var (
	svc        *service.Service
	sanitize   *sanitizer.Sanitizer
	emptyReply = &emptypb.Empty{}
)

func init() {
	sanitize = sanitizer.NewSanitizer()
}

func New(s *service.Service) *ServerInfo {
	srvConfig := &config.ServerConfig{}
	parseConfig("server.http", srvConfig)
	if srvConfig.WatchConfig {
		conf.OnChange(func(config *conf.Config) {
			parseConfig("server.http", srvConfig)
		})
	}

	server := server.New(srvConfig)
	gin.EnableJsonDecoderUseNumber()
	svc = s

	registerMiddlewares(server)
	router(server)

	addr := server.Start()
	_, port, _ := net.SplitHostPort(addr.String())
	serviceInfo := &registry.ServiceInfo{
		Name:    "server.http.example",
		Scheme:  "http",
		Addr:    fmt.Sprintf("%s://%s:%s", "http", xnet.InternalIP(), port),
		Version: "1.0.0",
	}
	binding.Validator.Engine().(*validator.Validate).SetTagName("validate")
	return &ServerInfo{Server: server, ServiceInfo: serviceInfo}
}

// ShouldBind 此方法返回的error不需要在输出到客户端,内部已经有了ctx.JSON() 处理
func ShouldBind(ctx *gin.Context, req interface{}) (err error) {
	if err = ctx.ShouldBind(req); err != nil {
		var errmsg string
		if errs, ok := err.(validator.ValidationErrors); ok {
			errmsg = utils.TranslateError(ctx, errs)
		} else {
			errmsg = err.Error()
		}

		// Reply StatusBadRequest
		validateResp := &reply.Response{
			Code:    strconv.Itoa(http.StatusBadRequest),
			Message: errmsg,
			TransId: trace.TraceID(ctx.Request.Context()),
		}
		ctx.JSON(http.StatusOK, validateResp)
		return
	}

	log.Info(ctx.Request.Context(), "request params", log.String("path", ctx.Request.URL.Path), log.Any("req", req))
	// clean xss attack
	sanitize.Sanitize(ctx, req)
	return
}

func parseConfig(configName string, srvConfig *config.ServerConfig) {
	if err := conf.Unmarshal(configName, srvConfig); err != nil {
		panic(fmt.Sprintf("unmarshal http server config fail, err msg %s", err.Error()))
	}
}

func registerMiddlewares(s *server.Server) {
	//	register middleware here
	s.Use(cors.Cors())
}

func router(s *server.Server) {
	registerAPI(s.Engine)
}
