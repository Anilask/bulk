package server

import (
	"context"
	sqlPackage "database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"cloud.google.com/go/storage"

	gcs "bulk/pkg/google-cloud-storage"

	"google.golang.org/grpc"

	"bulk/api/v1/handlers"
	"bulk/api/v1/services"
	"bulk/config"
	"bulk/logger"
	disbursementpb "bulk/pb/disbursement"
	filepb "bulk/pb/file"
	"bulk/pkg/database"
	"bulk/pkg/pubsub"
	"bulk/server/middleware"
	"bulk/store/agent"
	"bulk/tracer"
)

// channel to listen for signals
var signalChan = make(chan os.Signal, 1)

// Server is struct contain config and logger
type Server struct {
	cfg       *config.Config
	logger    logger.ILogger
	service   services.Repository
	handler   handlers.BulkDisbursementHandler
	bulkDB    *sqlPackage.DB
	gcsClient *storage.Client
	gcsRepo   gcs.Repository
}

var SVR *Server
var pricingDB *sqlPackage.DB

// New is used to create a new server
func New(cfg *config.Config, logger logger.ILogger) *Server {
	if SVR != nil {
		return SVR
	}
	SVR = &Server{
		cfg:    cfg,
		logger: logger,
	}

	//register bulkDB
	//MYSQL
	bulkDbConnection := database.NewDatabaseConnection(logger, cfg.MySQL)
	if bulkDbConnection == nil {
		logger.Fatal("Expecting bulkDB connection object but received nil")
	}

	bulkDB := bulkDbConnection.DBConnect()
	if bulkDB == nil {
		logger.Fatal("Expecting bulkDB connection object but received nil")
	}

	SVR.bulkDB = bulkDB

	bulkRepo := agent.NewBulkFileRepository(bulkDB, logger)

	bulkFileUploadRepo, err := agent.NewFileUploadRepository(
		&cfg.GCSBulk,
	)
	if err != nil {
		logger.Fatal("Expecting bulkFileUploadRepo connection object but received nil")
	}

	// Setup GCS Client
	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Error while creating Google Cloud Storage client : %v", err)
	}
	SVR.gcsClient = gcsClient

	SVR.gcsRepo = gcs.New().WithConfig(*cfg).
		WithLogger(logger).WithClient(gcsClient)
	pubsubClient := pubsub.NewPubSubClient(logger, context.Background(), cfg)
	if pubsubClient == nil {
		log.Fatal("Expecting pubsub connection but received nil")
	}
	// Register the Services
	SVR.service = services.New().
		WithConfig(*cfg).
		WithLogger(logger).
		WithBulkFileRepo(bulkRepo).
		WithBulkFileUploadRepo(bulkFileUploadRepo).
		WithGCSRepo(SVR.gcsRepo).
		WithPubSub(pubsubClient)

	// Register the Handler
	SVR.handler = *handlers.New().
		WithLogger(logger).
		WithService(SVR.service).WithConfig(SVR.cfg)

	return SVR
}

// Start the Server ...
func (s *Server) Start() {
	// init tracer
	tracer.New(s.cfg.Tracer.Enable, s.cfg.ProjectId, s.cfg.Tracer.TracerName, s.logger)
	defer tracer.Shutdown()

	addr := "0.0.0.0"
	if len(s.cfg.GRPCAddress) > 0 {
		if _, err := strconv.Atoi(s.cfg.GRPCAddress); err == nil {
			addr = fmt.Sprintf(":%v", s.cfg.GRPCAddress)
		} else {
			addr = s.cfg.GRPCAddress
		}
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.TraceUnaryServerInterceptor),
	)
	//Register BulkFileHandler
	filepb.RegisterBulkFileHandlerServer(grpcServer, &s.handler)
	//Register BulkDisbursementHandler
	disbursementpb.RegisterBulkDisbursementHandlerServer(grpcServer, &s.handler)
	// SIGINT handles Ctrl+C locally
	// SIGTERM handles Cloud Run termination signal
	// SIGKILL handles inappropriate/abrupt stopping of application
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		s.logger.Infof("[Start] starting server at port %v...", s.cfg.Port)

		lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", addr, s.cfg.Port))
		if err != nil {
			s.logger.Errorf("failed to listen: %v", err)
		}

		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Errorf("failed to serve: %s", err)
		}
		s.logger.Debugf("HTTP server start at %v", addr)

	}()

	// Receive output from signalChan.
	sig := <-signalChan
	s.logger.Infof("%s signal caught", sig)

	if err := s.Shutdown(context.Background()); err != nil {
		s.logger.Infof("server shutdown failed: %v", err)
	}
	s.logger.Info("server exited")
	s.logger.Debugf("HTTP server start %v", addr)
}

func (s *Server) Shutdown(_ context.Context) error {
	s.logger.Debugf("HTTP server stop %v", s.cfg.GRPCAddress)
	s.bulkDB.Close()
	s.gcsClient.Close()
	return nil
}
