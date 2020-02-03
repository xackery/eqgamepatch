package main

import (
	"context"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/rs/zerolog/log"
	"github.com/xackery/eqgamepatch/server"
)

var (
	// Version of build
	Version = "EXPERIMENTAL"
	ctx     context.Context
	cancel  context.CancelFunc
)

func main() {
	ctx, cancel = context.WithCancel(context.Background())
	log.Info().Msgf("starting eqgamepatch %s", Version)

	onExit := func() {
		log.Info().Msg("exiting")
		cancel()
		return
	}
	defer func() {
		cancel()
		systray.Quit()
	}()
	// Should be called at the very beginning of main().
	systray.RunWithAppWindow("EQGamePatch", 1024, 768, onReady, onExit)

	s, err := server.New(ctx, cancel)
	if err != nil {
		log.Error().Err(err).Msg("new server")
		return
	}
	fmt.Println(s)

	return
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("EQGamePatch")
	systray.SetTooltip("EQGamePatch")

	settings := systray.AddMenuItem("Settings", "Show Settings")
	launch := systray.AddMenuItem("Launch Everquest", "Launch Everquest")
	verify := systray.AddMenuItem("Verify Everquest", "Verify Everquest")
	quit := systray.AddMenuItem("Quit", "Quit EQGamePatch")

	// Sets the icon of a menu item. Only available on Mac.
	quit.SetIcon(icon.Data)

	systray.AddSeparator()
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("exiting via systray")
			systray.Quit()
			return
		case <-settings.ClickedCh:
			log.Info().Msg("todo: show settings")
		case <-launch.ClickedCh:
			log.Info().Msg("launching everquest")
		case <-verify.ClickedCh:
			log.Info().Msg("verifying everquest")
		case <-quit.ClickedCh:
			log.Info().Msg("quitting eqgamepatch")
			systray.Quit()
			cancel()
			return
		}
	}
}
