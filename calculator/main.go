package main

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type calculator struct {
	display *widget.Entry
	op      string
	val1    float64
	newNum  bool
}

func newCalculator(display *widget.Entry) *calculator {
	return &calculator{
		display: display,
		op:      "",
		val1:    0,
		newNum:  true,
	}
}

func (c *calculator) input(digit string) {
	if c.newNum {
		if digit == "." {
			c.display.SetText("0.")
		} else {
			c.display.SetText(digit)
		}
		c.newNum = false
	} else {
		if digit == "." && strings.Contains(c.display.Text, ".") {
			return
		}
		c.display.SetText(c.display.Text + digit)
	}
}

func (c *calculator) setOp(op string) {
	if !c.newNum && c.op != "" {
		c.compute()
	}
	val, err := strconv.ParseFloat(c.display.Text, 64)
	if err == nil {
		c.val1 = val
	}
	c.op = op
	c.newNum = true
}

func (c *calculator) compute() {
	if c.op == "" {
		return
	}
	val2, err := strconv.ParseFloat(c.display.Text, 64)
	if err != nil {
		return
	}

	var res float64
	switch c.op {
	case "+":
		res = c.val1 + val2
	case "-":
		res = c.val1 - val2
	case "×":
		res = c.val1 * val2
	case "÷":
		if val2 == 0 {
			c.display.SetText("Error: Div by 0")
			c.op = ""
			c.newNum = true
			return
		}
		res = c.val1 / val2
	case "%":
		res = float64(int64(c.val1) % int64(val2))
	}

	resStr := strconv.FormatFloat(res, 'f', -1, 64)
	c.display.SetText(resStr)
	c.op = ""
	c.newNum = true
}

func (c *calculator) clear() {
	c.display.SetText("0")
	c.op = ""
	c.val1 = 0
	c.newNum = true
}

func (c *calculator) toggleSign() {
	val, err := strconv.ParseFloat(c.display.Text, 64)
	if err == nil && val != 0 {
		c.display.SetText(strconv.FormatFloat(-val, 'f', -1, 64))
	}
}

func (c *calculator) backspace() {
	if c.newNum {
		return
	}
	txt := c.display.Text
	if len(txt) <= 1 || (len(txt) == 2 && strings.HasPrefix(txt, "-")) {
		c.display.SetText("0")
		c.newNum = true
	} else {
		c.display.SetText(txt[:len(txt)-1])
	}
}

func makeButton(label string, importance widget.ButtonImportance, action func()) *widget.Button {
	b := widget.NewButton(label, action)
	b.Importance = importance
	return b
}

func main() {
	myApp := app.NewWithID("com.example.gocalculator")
	myApp.Settings().SetTheme(theme.DarkTheme())

	w := myApp.NewWindow("Go GUI Calculator")
	w.Resize(fyne.NewSize(320, 420))

	display := widget.NewEntry()
	display.SetText("0")
	display.Disable()
	display.TextStyle = fyne.TextStyle{Bold: true}

	calc := newCalculator(display)

	numBtn := func(lbl string) *widget.Button {
		return makeButton(lbl, widget.MediumImportance, func() { calc.input(lbl) })
	}
	opBtn := func(lbl string) *widget.Button {
		return makeButton(lbl, widget.HighImportance, func() { calc.setOp(lbl) })
	}
	funcBtn := func(lbl string, act func()) *widget.Button {
		return makeButton(lbl, widget.LowImportance, act)
	}

	grid := container.NewGridWithColumns(4,
		// Row 1
		funcBtn("C", calc.clear),
		funcBtn("⌫", calc.backspace),
		funcBtn("±", calc.toggleSign),
		opBtn("÷"),

		// Row 2
		numBtn("7"),
		numBtn("8"),
		numBtn("9"),
		opBtn("×"),

		// Row 3
		numBtn("4"),
		numBtn("5"),
		numBtn("6"),
		opBtn("-"),

		// Row 4
		numBtn("1"),
		numBtn("2"),
		numBtn("3"),
		opBtn("+"),

		// Row 5
		numBtn("0"),
		numBtn("."),
		funcBtn("%", func() { calc.setOp("%") }),
		makeButton("=", widget.HighImportance, calc.compute),
	)

	content := container.NewBorder(
		container.NewPadded(display),
		nil, nil, nil,
		grid,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
