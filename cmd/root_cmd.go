package cmd

import (
	"cos-storager/config"
	"cos-storager/pkg/database"
	"cos-storager/router"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const (
	ProjectName = "cos-storage"
	Version     = "v0.0.1"
)

var rootCMD = &cobra.Command{
	Use: ProjectName,
	Run: runRootCMD,
}

func runRootCMD(cmd *cobra.Command, args []string) {
	database.Init()
	defer database.POSTGRES.Close()
	config.EnvInit()

	var storage router.Routers
	g := gin.Default()
	storage.Load(g)
	g.Run(":9099")
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
