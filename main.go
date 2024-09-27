package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("icon.ico"))
	systray.SetTooltip("Toggle between light and dark app mode")

	mToggle := systray.AddMenuItem("Toggle app mode", "Toggle between light and dark app mode")
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				toggleTheme()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
}

func toggleTheme() {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatalf("Failed to open registry key: %v", err)
	}
	defer func(key registry.Key) {
		err := key.Close()
		if err != nil {
			log.Printf("Failed to close registry key: %v", err)
		}
	}(key)

	lightTheme, _, err := key.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		log.Fatalf("Failed to read AppsUseLightTheme: %v", err)
	}

	var newTheme uint32
	if lightTheme == 0 {
		newTheme = 1
		fmt.Println("Switching to light app mode...")
	} else {
		newTheme = 0
		fmt.Println("Switching to dark app mode...")
	}

	err = key.SetDWordValue("AppsUseLightTheme", newTheme)
	if err != nil {
		log.Fatalf("Failed to set AppsUseLightTheme: %v", err)
	}
	fmt.Println("Mode switched successfully")
}

func getIcon(s string) []byte {
	icon, err := os.ReadFile(s)
	if err != nil {
		log.Fatalf("Failed to load icon: %v", err)
	}
	return icon
}
