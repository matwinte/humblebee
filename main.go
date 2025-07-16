package main

import (
	"fmt"

	_ "ariga.io/atlas-provider-gorm/gormschema"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/matwinte/humblebee/models"
)

func main() {
	testing := models.Plants{
		Name:              "Test Plant",
		OutdoorSowDate:    "05-01",
		DaysToGermination: 10,
		DaysToHarvest:     60,
	}
	fmt.Printf("Congrats, created example: %+v\n", testing)

	starterApp := app.New()
	welcomeWindow := starterApp.NewWindow("Humblebee")

	welcomeWindow.SetContent(widget.NewLabel("Welcome to Humblebee!"))
	welcomeWindow.ShowAndRun()
}
