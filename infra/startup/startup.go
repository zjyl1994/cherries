package startup

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zjyl1994/cherries/infra/util"
	"github.com/zjyl1994/cherries/infra/util/figlet"
	"github.com/zjyl1994/cherries/infra/vars"
	"github.com/zjyl1994/cherries/server"
	"github.com/iancoleman/strcase"
	_ "github.com/joho/godotenv/autoload"
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Startup() (err error) {
	envPrefix := strcase.ToScreamingSnake(vars.APP_NAME)
	vars.LISTEN_ADDR = util.COALESCE(os.Getenv(envPrefix+"_LISTEN"), vars.DEFAULT_LISTEN)
	vars.DEBUG_MODE, _ = strconv.ParseBool(os.Getenv(envPrefix + "_DEBUG"))
	if vars.DEBUG_MODE {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("Debug mode enabled")
	} else {
		if banner, e := figlet.Render(vars.APP_NAME); e == nil {
			fmt.Println(banner)
		}
	}

	databaseFile := strcase.ToSnake(vars.APP_NAME) + ".db"
	logrus.Debugln("Database file:", databaseFile)
	vars.DB, err = gorm.Open(sqlite.Open(databaseFile), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		return err
	}
	err = vars.DB.Exec("PRAGMA journal_mode=WAL;").Error
	if err != nil {
		return err
	}
	err = vars.DB.AutoMigrate()
	if err != nil {
		return err
	}

	err = server.Run(vars.LISTEN_ADDR)
	if err != nil {
		return err
	}
	// cleanup
	sqlDB, _ := vars.DB.DB()
	sqlDB.Close()
	logrus.Infoln("Server stopped")
	return nil
}
