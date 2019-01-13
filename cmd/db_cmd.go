package cmd

import (
	"cos-storager/model"
	"cos-storager/pkg/database"

	"github.com/spf13/cobra"
	"qiniupkg.com/x/log.v7"
)

func init() {
	rootCMD.AddCommand(dbCMD)
}

const (
	dbCommand = "db"
)

var dbCMD = &cobra.Command{
	Use: dbCommand,
	Run: runCMD,
}

func runCMD(cmd *cobra.Command, args []string) {
	database.Init()
	defer database.POSTGRES.Close()
	if len(args) == 0 {
		log.Println("you should add one args at least")
		return
	}
	if args[0] == "migrate" {
		migration()
	}
}

func migration() {
	collections := model.GetAllModels()
	for _, collection := range collections {
		database.POSTGRES.AutoMigrate(collection)
	}
}
