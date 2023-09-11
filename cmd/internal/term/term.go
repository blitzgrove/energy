package term

import (
	"github.com/energye/energy/v2/cmd/internal/consts"
	"github.com/energye/golcl/lcl/rtl/version"
	"github.com/pterm/pterm"
	"os"
)

const (
	energyCliVersion = "1.0.0"
)

var Logger *pterm.Logger
var TermOut = new(Termout)
var Section *pterm.SectionPrinter

func init() {
	if consts.IsWindows {
		// < windows 10 禁用颜色
		version.VersionInit()
		ov := version.OSVersion
		if ov.Major < 10 {
			pterm.DisableColor()
		}
	}

	// logger
	Logger = pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	TermOut = new(Termout)
	Logger.ShowTime = false
	Logger.ShowCaller = false

	// Section
	Section = &pterm.SectionPrinter{
		Style:           &pterm.Style{pterm.Bold, pterm.FgLightYellow},
		Level:           1,
		IndentCharacter: "$",
	}
}

type Termout struct {
}

func (m *Termout) Write(p []byte) (n int, err error) {
	Section.Println(string(p))
	return len(p), nil
}

func GoENERGY() {
	pterm.Println(pterm.LightBlue("      GO\n") + pterm.LightBlue("    ENERGY"))
}

func BoxPrintln(a ...any) {
	pterm.DefaultBox.Println(a...)
}

func Println(a ...any) {
	pterm.Println(a...)
}

func TextInputWith(text, delimiter string) string {
	tiw := pterm.DefaultInteractiveTextInput.WithMultiLine(false).WithOnInterruptFunc(func() {
		os.Exit(1)
	})
	tiw.DefaultText = text
	tiw.Delimiter = delimiter
	result, err := tiw.Show()
	if err != nil {
		return ""
	}
	return result
}

// WithBoolean helps an option setter (WithXXX(b ...bool) to return true, if no boolean is set, but false if it's explicitly set to false.
func WithBoolean(b []bool) bool {
	if len(b) == 0 {
		b = append(b, true)
	}
	return b[0]
}

// NewCancelationSignal for keeping track of a cancelation
func NewCancelationSignal(interruptFunc func()) (func(), func()) {
	canceled := false

	cancel := func() {
		canceled = true
	}

	exit := func() {
		if canceled && interruptFunc != nil {
			interruptFunc()
		}
	}

	return cancel, exit
}
