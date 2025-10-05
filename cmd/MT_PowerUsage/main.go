package main

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/olegbilovus/MT_PowerUsage/internal/database"
	"github.com/olegbilovus/MT_PowerUsage/pkg/plug"
)

var (
	envVars = []string{
		"PLUG_IP",
		"PLUG_TYPE",
		"FREQ",
		"SQLITE_DB",
	}
)

//goland:noinspection D
func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true,
		PadLevelText: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", filename
		}})
	log.SetLevel(log.TraceLevel)

	if _, foundEnv := os.LookupEnv("PLUG_IP"); !foundEnv {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("Error loading .env file")
		}
		log.Info("Loaded Env from .env")
	} else {
		log.Info("Skipped loading .env because found PLUG_IP already in the Env")
	}

	if missingEnvVars := checkEnvVars(); len(missingEnvVars) > 0 {
		log.Fatalln("The following Env vars are missing: ", strings.Join(missingEnvVars, ", "))
	}

	plugType, err := strconv.Atoi(os.Getenv("PLUG_TYPE"))
	if err != nil {
		log.Fatalf("Invalid PLUG_TYPE: %v", err)
	}
	ip := os.Getenv("PLUG_IP")
	client := &http.Client{Timeout: 5 * time.Second}
	var powerPlug plug.Plug
	switch plugType {
	case 1:
		powerPlug = plug.ShellyPlugS{
			Ip:     ip,
			Client: client,
		}
	case 2:
		powerPlug = plug.ShellyPlugSv2{
			Ip:     ip,
			Client: client,
		}
	default:
		log.Fatalf("PLUG_TYPE must be a value in : 1, 2. %d is not valid", plugType)
	}

	dbPath := path.Join("dbs", os.Getenv("SQLITE_DB"))
	log.Infof("saving data in %s", dbPath)

	dbSQLite, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer dbSQLite.Close()
	db := database.SQLite{DB: dbSQLite}
	if err := db.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	if os.Getenv("RESET_DB") == "true" {
		if err := db.Reset(); err != nil {
			log.Fatalf("Error resetting the db: %v", err)
		}
		log.Warning("DB reset")
	}

	if err := powerPlug.TurnOn(); err != nil {
		log.Fatalf("Could not turn on the powerPlug: %v", err)
	}

	freq, err := strconv.Atoi(os.Getenv("FREQ"))
	if err != nil {
		log.Fatalf("Invalid FREQ: %v", err)
	}
	freqDur := time.Duration(freq) * time.Second

	// Manage signals
	var receivedSignal os.Signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	var i int32
	for ; receivedSignal == nil; receivedSignal = sleep(freqDur, sigs) {
		i++
		go func(i int32) {
			start := time.Now().UTC()
			log.Info("Started check #", i)

			if load, err := powerPlug.Load(); err != nil {
				log.Errorf("Error getting load: %v", err)
			} else {
				if err := db.Write(start, load); err != nil {
					log.Errorf("Error writing to db: %v", err)
				} else {
					log.WithFields(log.Fields{
						"name":  powerPlug.Name(),
						"power": load,
					}).Info("Plug")
				}
			}

			log.WithFields(log.Fields{
				"duration": time.Since(start),
			}).Info("Finished check #", i)
		}(i)
	}

}

func checkEnvVars() []string {
	var missingEnvVars []string
	for _, envVar := range envVars {
		if _, found := os.LookupEnv(envVar); !found {
			missingEnvVars = append(missingEnvVars, envVar)
		}
	}

	return missingEnvVars
}

func sleep(duration time.Duration, sigChan <-chan os.Signal) os.Signal {
	select {
	case sig := <-sigChan:
		return sig
	case <-time.After(duration):
		return nil
	}
}
