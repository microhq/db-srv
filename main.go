package main

import (
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/db-srv/db"
	_ "github.com/micro/db-srv/db/mysql"
	"github.com/micro/db-srv/handler"

	proto "github.com/micro/db-srv/proto/db"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.db"),
		micro.Version("latest"),

		micro.Flags(
			cli.StringFlag{
				Name:   "database_service_namespace",
				EnvVar: "DATABASE_SERVICE_NAMESPACE",
				Usage:  "The namespace used when looking up databases in registry e.g go.micro.db",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_service_namespace")) > 0 {
				db.DBServiceNamespace = c.String("database_service_namespace")
			}
		}),
	)

	service.Init(
		// init the db
		micro.BeforeStart(func() error {
			return db.Init(service.Client().Options().Selector)
		}),
	)

	proto.RegisterDBHandler(service.Server(), new(handler.DB))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
