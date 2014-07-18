package lsproc

import (
	"fmt"
	"github.com/elasticsearch/kriterium/component/process"
	"github.com/elasticsearch/kriterium/errors"
	"github.com/elasticsearch/kriterium/panics"
	"log"
	"lsf"
	"lsf/fs"
	"lsf/lsfun"
	"lsf/lslib"
	"lsf/schema"
	"lsf/system"
	"time"
)

type TrackConfig struct {
	Env          *lsf.Environment // any lsf func/proc config TODO refactor
	Debug        bool             // any lsf func/proc config TODO refactor
	Verbose      bool             // any lsf func/proc config TODO refactor
	StreamId     string
	Delay        time.Duration
	MaxCacheSize uint16
	MaxCacheAge  fs.InfoAge
	// REVU: TODO: log rotate config here -- BEGIN
	MaxRotations uint8
	LogFileSize  uint64
	// REVU: TODO: log rotate config here -- END
}

func NewDefaultConfig(env *lsf.Environment, streamId string) *TrackConfig {
	return &TrackConfig{
		Env:          env,
		StreamId:     streamId,
		Delay:        time.Second,
		MaxCacheSize: 16,
		MaxRotations: uint8(1),
		LogFileSize:  uint64(1024 * 128),
	}
}

func (c *TrackConfig) SetMaxCacheAge(infoAge fs.InfoAge) error {
	c.MaxCacheAge = infoAge
	return nil
}

// TODO: move to env or system
func docIdForStream(id string) string {
	return fmt.Sprintf("stream.%s.stream", id)
}

// async function : use go lsproc.TrackProcess(..)
func TrackProcess(controller process.Controller, config *TrackConfig) {
	// TODO: various asserts here and panic if necessary (or just for fun :)
	if controller == nil {
		panic(errors.IllegalArgument("controller is nil"))
	}
	if config == nil {
		panic(errors.IllegalArgument("config is nil"))
	}

	defer panics.AsyncRecover(controller.Respond())
	env := config.Env

	// Load stream doc and get LogStream instance

	docId := docIdForStream(config.StreamId)
	doc, e := env.LoadDocument(docId)
	panics.OnError(e, "no such stream:", config.StreamId)
	panics.OnNil(doc, fmt.Sprintf("BUG: system doc for stream %q", config.StreamId))

	logStream := schema.DecodeLogStream(doc)

	// Run command in exclusive mode and lockout
	// other track ops for this specific stream

	opLock, opLockId, e := env.ExclusiveResourceOp(system.Op.StreamTrack, config.StreamId, "lsf-command-track")
	panics.OnError(e)
	defer func() {
		panics.OnError(opLock.Unlock())
		if config.Debug {
			log.Printf("DEBUG: unlocked %s\n", opLockId)
		}
	}()

	if config.Debug {
		log.Printf("DEBUG: locked %s\n", opLockId)
	}

	/// run the track process ///////////////////////////////////

	var scout lsfun.TrackScout = lsfun.NewTrackScout(logStream.Path, logStream.Pattern, config.MaxCacheSize, config.MaxCacheAge)
	eventlogBasepath := env.Port()
	eventlogBasename := config.StreamId + ".trackscout.event.log" // REVU: don't like this

	rotator, e := lslib.NewRotatingFileWriter(eventlogBasepath, eventlogBasename, config.MaxRotations, config.LogFileSize)
	panics.OnError(e, "NewFileRotator")

	//	defer func() { // WHY IS THIS CAUSING PROBLEMS? (the unlock defer above won't run if this does ..)
	//		rotator.Close()
	//		if config.Debug {
	//			log.Printf("DEBUG: closing event log rotator at %s/%s", eventlogBasepath, eventlogBasename)
	//		}
	//	}()

	var cmd process.CommandCode
	cmd = <-controller.Command()

	// REVU: this looks like boiler plate to me .. BEGIN
	switch cmd {
	case process.Start: //
		controller.Respond() <- process.Start
	case process.Stop, process.Abort:
		controller.Respond() <- process.Stop
		return
	default:
		panic(errors.IllegalState("unexpected commmand from controller", cmd))
	}
	for {
		select {
		case cmd = <-controller.Command():
			switch cmd {
			case process.Stop:
				controller.Respond() <- process.Stop
				return
			}
		case <-time.After(config.Delay):
			report, e := scout.Report()
			panics.OnError(e, "main", "scout.Report") // REVU: wrong. send error via channel and close

			log.Println("--- events -------------------------------------------")
			for _, event := range report.Events {
				if event.Code != lsfun.TrackEvent.KnownFile { // printing NOP events gets noisy
					log.Println(event)
					rotator.Write([]byte(event.String() + "\n"))
				}
			}

			objects := scout.ObjectMap()

			log.Println("--- objects ------------------------------------------")
			objectsByAge := fs.AsObjectMap(objects).Sort(fs.ObjectIterationOrder.ByAge, fs.IterationDirection.Ascending)
			for _, fsobj := range objectsByAge {
				log.Println(fsobj.String())
			}
			log.Println()
		}
	}
	// REVU: this looks like boiler plate to me .. END

	return
}
