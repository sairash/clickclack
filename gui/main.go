package main

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CustomTheme struct{}

var (
	soundPacks     = []string{"sound-pack-1200000000001", "sound-pack-1200000000002", "sound-pack-1200000000003", "sound-pack-1200000000004", "sound-pack-1200000000005", "sound-pack-1200000000006", "sound-pack-1200000000007", "sound-pack-1200000000008", "sound-pack-1200000000009", "sound-pack-1200000000010", "sound-pack-1200000000011", "sound-pack-1200000000012", "sound-pack-1200000000014", "sound-pack-1203000000018", "sound-pack-1203000000026", "sound-pack-1203000000023", "sound-pack-1203000000024", "sound-pack-1203000000025", "sound-pack-1203000000027", "sound-pack-1203000000028", "sound-pack-1203000000084", "sound-pack-1203000000030", "sound-pack-1203000000031", "sound-pack-1203000000032", "sound-pack-1203000000034", "sound-pack-1203000000058", "sound-pack-1203000000036", "sound-pack-1203000000037", "sound-pack-1203000000038", "sound-pack-1203000000039", "sound-pack-1203000000041", "sound-pack-1203000000081", "sound-pack-1203000000042", "sound-pack-1203000000048", "sound-pack-1203000000044", "sound-pack-1203000000045", "sound-pack-1203000000046", "sound-pack-1203000000021", "sound-pack-1203000000079", "sound-pack-1203000000019", "sound-pack-1203000000051", "sound-pack-1203000000053", "sound-pack-1203000000054", "sound-pack-1203000000091", "sound-pack-1203000000056", "sound-pack-1203000000057", "sound-pack-1203000000016", "sound-pack-1203000000063", "sound-pack-1203000000069", "sound-pack-1203000000074", "sound-pack-1203000000093", "sound-pack-1203000000071", "sound-pack-1203000000073", "sound-pack-1203000000075", "sound-pack-1203000000077", "sound-pack-1203000000076", "sound-pack-1203000000062", "sound-pack-1203000000064", "sound-pack-1203000000072", "sound-pack-1203000000020", "sound-pack-1203000000022", "sound-pack-1203000000082", "sound-pack-1203000000083", "sound-pack-1203000000087", "sound-pack-1203000000088", "sound-pack-1203000000089", "sound-pack-1203000000043", "sound-pack-1203000000090", "sound-pack-1203000000092", "sound-pack-1203000000047", "sound-pack-1203000000080", "sound-pack-1203000000085", "sound-pack-1203000000094", "sound-pack-1203000000055", "sound-pack-1203000000096", "sound-pack-1203000000029", "sound-pack-1203000000052", "sound-pack-1203000000040", "sound-pack-1203000000049", "sound-pack-1203000000050", "sound-pack-1203000000066", "sound-pack-1203000000067", "sound-pack-1203000000068", "sound-pack-1203000000078", "sound-pack-1203000000086", "sound-pack-1203000000033"}
	primaryPurple  = color.NRGBA{R: 0x0C, G: 0x09, B: 0x6D, A: 0xFF}
	accentPink     = color.NRGBA{R: 0xD4, G: 0x69, B: 0xD0, A: 0xFF}
	cyan           = color.NRGBA{R: 0x4B, G: 0xB3, B: 0xC3, A: 0xFF}
	vividBlue      = color.NRGBA{R: 0x5F, G: 0x6F, B: 0xF8, A: 0xFF}
	vividBlueTest  = color.NRGBA{R: 0, G: 0, B: 0, A: 0x5F}
	goldenYellow   = color.NRGBA{R: 0xEC, G: 0xCF, B: 0x29, A: 0xFF}
	pureWhite      = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	gradientCenter = color.NRGBA{R: 0x16, G: 0x11, B: 0x79, A: 0xFF}
	gradientEdge   = color.NRGBA{R: 0x09, G: 0x05, B: 0x68, A: 0xFF}
)

func (CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return accentPink
	case theme.ColorNameButton:
		return vividBlueTest
	case theme.ColorNameForeground:
		return color.White
	case theme.ColorNameSelection:
		return accentPink
	case theme.ColorNameInputBackground:
		return vividBlueTest
	case theme.ColorNamePrimary:
		return cyan

	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func truncateFilename(filename string) string {
	maxLength := 15
	if len(filename) <= maxLength {
		return filename
	}

	extIndex := strings.LastIndex(filename, ".")
	if extIndex == -1 {
		return filename[:maxLength-3] + "..."
	}

	name := filename[:extIndex]
	ext := filename[extIndex:]

	keep := maxLength - len(ext) - 3

	if keep <= 0 {
		return filename[:maxLength-3] + "..."
	}

	truncatedName := name[:keep]
	return truncatedName + "..." + ext
}

func loadLogo() fyne.Resource {
	logoPath := "assets/logo.png"
	logoURI := storage.NewFileURI(logoPath)
	logoRes, _ := storage.LoadResourceFromURI(logoURI)
	return logoRes
}

func main() {

	myApp := app.New()
	myApp.Settings().SetTheme(&CustomTheme{})
	window := myApp.NewWindow("Vibrant Typing Sound Generator")

	// Gradient background
	gradient := canvas.NewRadialGradient(gradientCenter, gradientEdge)
	gradient.CenterOffsetX = 0.5
	gradient.CenterOffsetY = 0.5

	// Custom theme colors
	// textColor := color.White
	dividerColor := color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}

	// Logo (replace with your own image)
	logo := canvas.NewImageFromResource(loadLogo())
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 100))

	// Instruction text
	instruction := widget.NewLabel("Hear clicks and clacks when you type and do your business.")
	instruction.Alignment = fyne.TextAlignCenter
	instruction.TextStyle = fyne.TextStyle{Bold: true}
	// instruction.Wrapping = fyne.TextWrapWord
	instruction.Resize(fyne.NewSize(380, 60))

	// Divider line
	divider := canvas.NewRectangle(dividerColor)
	divider.SetMinSize(fyne.NewSize(380, 2))

	// selected_dropdown_value := ""
	// Soundpack dropdown
	dropdown := widget.NewSelect(soundPacks, func(s string) {

		// selected_dropdown_value = strings.Split(s, "sound-pack-")[1]
	})
	dropdown.PlaceHolder = "Choose Soundpack"

	dropdown.SetSelectedIndex(0)

	browse_path_value := make([]string, 2)

	createFileButton := func(buttonText string, icon fyne.Resource, index int) *widget.Button {
		btn := widget.NewButtonWithIcon(buttonText, icon, nil)
		btn.OnTapped = func() {
			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil || reader == nil {
					return
				}

				browse_path_value[index] = reader.URI().Path()

				// Update button text with filename
				fileName := filepath.Base(reader.URI().Path())
				btn.SetText(truncateFilename(fileName))
				btn.Refresh()
			}, window)
		}
		return btn
	}

	// Create file buttons
	leftFile := createFileButton("Choose Sound File", theme.ContentUndoIcon(), 0)
	rightFile := createFileButton("Choose Sound File", theme.ContentRedoIcon(), 1)

	// Place buttons side-by-side with spacing
	fileButtons := container.NewHBox(
		leftFile,
		rightFile,
	)

	dropdownLabel := widget.NewLabel("Choose a keyboard soundpack from MechVibes")
	dropdownLabel.Alignment = fyne.TextAlignLeading
	dropdownLabel.TextStyle = fyne.TextStyle{Bold: false}

	mouseLabel := widget.NewLabel("Choose Left and Right mouse sounds")
	mouseLabel.Alignment = fyne.TextAlignLeading
	mouseLabel.TextStyle = fyne.TextStyle{Bold: false}

	volumeValue := widget.NewLabel("70%")
	volumeValue.Alignment = fyne.TextAlignCenter
	volumeValue.TextStyle = fyne.TextStyle{Bold: true}

	// Create volume label row
	volumeLabelRow := container.NewHBox(
		widget.NewLabelWithStyle("Volume", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		layout.NewSpacer(),
		volumeValue,
	)

	volumeSlider := widget.NewSlider(0, 100)
	volumeSlider.Value = 70 // Default volume
	volumeSlider.Step = 1
	volumeSlider.Orientation = widget.Horizontal

	volumeSlider.OnChanged = func(value float64) {
		volumeValue.SetText(fmt.Sprintf("%d%%", int(value)))
	}

	textArea := widget.NewMultiLineEntry()
	textArea.SetPlaceHolder("Info From ClickClack is shown here")
	textArea.Wrapping = fyne.TextWrapWord
	textArea.Disable()
	textArea.TextStyle.Monospace = true

	infoContainer := container.NewStack(
		canvas.NewRectangle(vividBlue),
		container.NewVBox(
			widget.NewLabelWithStyle("Info: Stop the program to Apply changes", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
		),
	)

	infoContainer.Hide()

	// Full-width progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()
	progress.Resize(fyne.NewSize(380, 10))

	fileButtons.Layout = layout.NewGridLayout(2)

	startButton := widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func() {

		textArea.Append("Progress Bar shown \n")
		infoContainer.Hide()
		progress.Show()

		// conf, err := sound.InitConfig(volumeSlider.Value/100, selected_dropdown_value, browse_path_value[0], browse_path_value[1])

		// go sound.CreateSound(conf, eventChan)
		go func() {
			time.Sleep(20 * time.Second)
			progress.Hide()

			// sound.Exit <- true

			infoContainer.Show()
			textArea.Append("Progress Bar hidden \n")

		}()
	})

	startButton.Importance = widget.HighImportance

	footerText := widget.NewRichTextFromMarkdown(
		`Created with ❤️ by [Sairash](sairashgautam.com.np)`,
	)
	footerText.Wrapping = fyne.TextWrapWord

	footerText2 := widget.NewRichTextFromMarkdown(
		`Buy me a ☕ [coffee](buymeacoffee.com/sairash)`,
	)
	footerText2.Wrapping = fyne.TextWrapWord

	footer := container.NewHBox(
		// footerText2,
		footerText)

	footer.Layout = layout.NewGridLayout(2)

	// Updated content layout with padding
	content := container.NewVBox(
		container.NewCenter(logo),
		container.NewCenter(instruction),
		container.NewCenter(divider),

		dropdownLabel,
		dropdown,
		mouseLabel,
		fileButtons,
		volumeLabelRow,
		volumeSlider,
		infoContainer,
		layout.NewSpacer(), // Space before start button
		progress,           // Full-width progress (not centered)
		textArea,
		container.NewStack(startButton),
	)

	// Create background rectangle
	mainContainer := container.NewStack(
		gradient,
		container.NewBorder(
			nil,                         // Top
			container.NewPadded(footer), // Bottom
			nil,                         // Left
			nil,                         // Right
			container.NewPadded(content),
		),
	)

	window.SetContent(mainContainer)
	window.Resize(fyne.NewSize(420, 590))
	window.ShowAndRun()
}
