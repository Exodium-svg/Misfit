package main

import (
	"MisFitsDiscord/Discord"
	"MisFitsDiscord/Game"
	"MisFitsDiscord/Utils"
	"MisFitsDiscord/Website"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func autoSave() {
	for {
		time.Sleep(time.Minute * 5)
		fmt.Printf("Auto saved triggered...\n")
		Utils.SaveStrings()
		Game.SaveAscensions()
		Game.SaveTaoists()
	}
}

func main() {
	Utils.LoadStrings()
	Game.LoadTaoists()
	Game.StartZodiacCycle()

	go autoSave()
	env, err := Utils.NewEnv("environment.env")

	if nil != err {
		fmt.Printf("failed to read env: %s\n", err.Error())
		return
	}

	session := Discord.Start(&env)
	defer session.Close()

	if nil == session {
		return
	}

	Website.Start(&env, session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("Shutting down...")

	Game.SaveTaoists()
	Utils.SaveStrings()
}
