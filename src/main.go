package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/repository/cache"
	crsRepo "github.com/bookpanda/mygraderlist-backend/src/app/repository/course"
	emjRepo "github.com/bookpanda/mygraderlist-backend/src/app/repository/emoji"
	lkRepo "github.com/bookpanda/mygraderlist-backend/src/app/repository/like"
	prblmRepo "github.com/bookpanda/mygraderlist-backend/src/app/repository/problem"
	rtngRepo "github.com/bookpanda/mygraderlist-backend/src/app/repository/rating"
	usrRepo "github.com/bookpanda/mygraderlist-backend/src/app/repository/user"
	crsService "github.com/bookpanda/mygraderlist-backend/src/app/service/course"
	emjService "github.com/bookpanda/mygraderlist-backend/src/app/service/emoji"
	lkService "github.com/bookpanda/mygraderlist-backend/src/app/service/like"
	prblmService "github.com/bookpanda/mygraderlist-backend/src/app/service/problem"
	rtngService "github.com/bookpanda/mygraderlist-backend/src/app/service/rating"
	usrService "github.com/bookpanda/mygraderlist-backend/src/app/service/user"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	"github.com/bookpanda/mygraderlist-backend/src/database"
	seed "github.com/bookpanda/mygraderlist-backend/src/database/seeds"
	crs_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/course"
	emj_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/emoji"
	lk_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/like"
	prblm_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/problem"
	rtng_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/rating"
	usr_proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func handleArgs(db *gorm.DB) {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			err := seed.Execute(db, args[1:]...)
			if err != nil {
				log.Fatal().
					Str("service", "seeder").
					Msg("Not found seed")
			}
			os.Exit(0)
		}
	}
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		sig := <-s

		log.Info().
			Str("service", "graceful shutdown").
			Msgf("got signal \"%v\" shutting down service", sig)

		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Error().
				Str("service", "graceful shutdown").
				Msgf("timeout %v ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("cleaning up: %v", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Error().
						Str("service", "graceful shutdown").
						Err(err).
						Msgf("%v: clean up failed: %v", innerKey, err.Error())
					return
				}

				log.Info().
					Str("service", "graceful shutdown").
					Msgf("%v was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()
		close(wait)
	}()

	return wait
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend load config").
			Msg("Failed to start service")
	}

	db, err := database.InitDatabase(&conf.Database)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend database").
			Msg("Failed to start service")
	}

	cacheDB, err := database.InitRedisConnect(&conf.Redis)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend redis cache").
			Msg("Failed to start service")
	}

	handleArgs(db)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.App.Port))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend server").
			Msg("Failed to start service")
	}

	grpcServer := grpc.NewServer()

	cacheRepo := cache.NewRepository(cacheDB)

	userRepo := usrRepo.NewRepository(db)
	userService := usrService.NewService(userRepo, conf.App)

	courseRepo := crsRepo.NewRepository(db)
	courseService := crsService.NewService(courseRepo, conf.App)

	problemRepo := prblmRepo.NewRepository(db)
	problemService := prblmService.NewService(problemRepo, cacheRepo, conf.App)

	likeRepo := lkRepo.NewRepository(db)
	likeService := lkService.NewService(likeRepo, conf.App)

	emojiRepo := emjRepo.NewRepository(db)
	emojiService := emjService.NewService(emojiRepo, conf.App)

	ratingRepo := rtngRepo.NewRepository(db)
	ratingService := rtngService.NewService(ratingRepo, conf.App)

	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	usr_proto.RegisterUserServiceServer(grpcServer, userService)
	crs_proto.RegisterCourseServiceServer(grpcServer, courseService)
	prblm_proto.RegisterProblemServiceServer(grpcServer, problemService)
	lk_proto.RegisterLikeServiceServer(grpcServer, likeService)
	emj_proto.RegisterEmojiServiceServer(grpcServer, emojiService)
	rtng_proto.RegisterRatingServiceServer(grpcServer, ratingService)
	reflection.Register(grpcServer)

	go func() {
		log.Info().
			Str("service", "backend").
			Msgf("MyGraderList backend starting at port %v", conf.App.Port)

		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal().
				Err(err).
				Str("service", "backend").
				Msg("Failed to start service")
		}
	}()

	wait := gracefulShutdown(context.Background(), 2*time.Second, map[string]operation{
		"database": func(ctx context.Context) error {
			sqlDb, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDb.Close()
		},
		"cache": func(ctx context.Context) error {
			return cacheDB.Close()
		},
		"server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	<-wait

	grpcServer.GracefulStop()
	log.Info().
		Str("service", "backend").
		Msg("Closing the listener")
	lis.Close()
	log.Info().
		Str("service", "backend").
		Msg("End of Program")
}
