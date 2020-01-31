package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/gobuffalo/packr"
	"github.com/lxn/walk"
	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/xackery/eqgamepatch/checksum"
	"github.com/xackery/eqgamepatch/version"
)

// UI Wraps the entire UI
type UI struct {
	progressValue float64
	box           packr.Box
	title         string
	version       string
}

// New creates a new UI
func New(title string) (*UI, error) {
	ui := new(UI)

	ui.box = packr.NewBox("./assets")
	var err error
	ui.version, err = version.Detect()
	if err != nil {
		return nil, errors.Wrap(err, "version detect")
	}
	log.Debug().Msgf("detected %s client", version.Name(ui.version))
	ui.title = fmt.Sprintf("EQGamePatch %s", version.Name(ui.version))
	if len(title) > 0 {
		ui.title = title
	}
	return ui, nil
}

// Start starts a new UI instance
func (ui *UI) Start() error {
	var inTE, outTE *walk.TextEdit

	declarative.MainWindow{
		Title:   "SCREAMO",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &inTE},
					TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			PushButton{
				Text: "SCREAM",
				OnClicked: func() {
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
	}.Run()
	return nil
}

func (ui *UI) scanDirectory() {

	blacklist := []string{
		"texture.txt", //logs texture errors
		"UIErrors.txt",
		"notes.txt",
		"eqOptions1.opt",
		"checksum.dat",
	}

	iniWhiteList := []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		//if info.IsDir() && info.Name() == subDirToSkip {
		//	fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		//	return filepath.SkipDir
		//}

		if info.IsDir() {
			if info.Name() == "Logs" {
				return filepath.SkipDir
			}
			if info.Name() == "patch" {
				return filepath.SkipDir
			}
			if info.Name() == "userdata" {
				return filepath.SkipDir
			}
			return nil
		}

		for _, f := range blacklist {
			if path == f {
				return nil
			}
		}

		if filepath.Ext(path) == ".ini" {
			for _, f := range iniWhiteList {
				if path == f {
					return nil
				}
			}
		}
		sum, err := checksum.Get(path)
		if err != nil {
			return err
		}

		err = checksum.Add("patch", "rebuildeq", path, sum)
		if err != nil {
			return err
		}

		ui.log(fmt.Sprintf("%s: %s %t", path, sum, checksum.IsDirty("rebuildeq", path)))

		//fmt.Printf("visited file or dir: %q and got sum %s and dirty %t\n", path, sum, checksum.IsDirty("rebuildeq", path))
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to walk")
		return
	}
}

func (ui *UI) onServerSelect(value string) {
	log.Debug().Msgf("selected %s", value)
}

func (ui *UI) log(message string) {
	newMsg := fmt.Sprintf("%s\n%s", ui.textLog.Text, message)
	if len(ui.textLog.Text) == 0 {
		newMsg = message
	}
	ui.textLog.SetText(newMsg)
}

func isEverquestRunning() bool {
	processes, err := ps.Processes()
	if err != nil {
		log.Warn().Err(err).Msg("could not get process list")
		return false
	}
	for _, proc := range processes {
		if proc.Executable() == "eqgame.exe" {
			return true
		}
	}
	return false
}
