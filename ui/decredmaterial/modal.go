package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"golang.org/x/image/math/fixed"
)

type Modal struct {
	titleLabel     Label
	titleSeparator *Line

	overlayColor    color.RGBA
	backgroundColor color.RGBA

	hasCalculatedWidth bool
}

func (t *Theme) Modal(title string) *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	return &Modal{
		titleLabel:     t.H6(title),
		titleSeparator: t.Line(),

		overlayColor:    overlayColor,
		backgroundColor: t.Color.Surface,

		hasCalculatedWidth: false,
	}
}

func (m *Modal) SetTitle(title string) {
	m.titleLabel.Text = title
}

func (m *Modal) layoutControls(gtx *layout.Context, controlMaterials []Button, controlWidgets []*widget.Button) (func(), int) {
	totalControlWidth := 0
	children := []layout.FlexChild{}
	for i := range controlMaterials {
		index := i
		totalControlWidth += m.calulateButtonWidth(gtx, controlMaterials[index])
		children = append(children, layout.Rigid(func() {
			in := layout.Inset{}
			if index != 0 {
				in.Left = unit.Dp(5)
			}

			in.Layout(gtx, func() {
				controlMaterials[index].Layout(gtx, controlWidgets[index])
			})
		}))
	}

	return func() {
		layout.Center.Layout(gtx, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
		})
	}, totalControlWidth
}

// Layout lays out the modal with specified title and width
// If the passed width is 0, then a default width is used for the modal
// If not, the modal assumes the width passed to it
func (m *Modal) Layout(gtx *layout.Context, widgets []func(), controlMaterials []Button, controlWidgets []*widget.Button) {
	maxWidth := m.calculateTitleWidth(gtx)

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillMax(gtx, m.overlayColor)
			new(widget.Button).Layout(gtx)
		}),
		layout.Stacked(func() {
			widgetFuncs := []func(){
				func() {
					m.titleLabel.Layout(gtx)
				},
				func() {
					m.titleSeparator.Width = gtx.Constraints.Width.Max
					m.titleSeparator.Layout(gtx)
				},
			}

			for i := range widgets {
				index := i
				widgetFuncs = append(widgetFuncs, func() {
					widgets[index]()
				})
			}

			if controlMaterials != nil {
				controlFunc, controlWidth := m.layoutControls(gtx, controlMaterials, controlWidgets)
				if controlWidth > maxWidth {
					maxWidth = controlWidth
				}
				widgetFuncs = append(widgetFuncs, controlFunc)
			}

			if !m.hasCalculatedWidth {
				(&layout.List{Axis: layout.Vertical, Alignment: layout.Middle}).Layout(gtx, len(widgetFuncs), func(i int) {
					fillMax(gtx, m.backgroundColor)
					layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
				})
				m.hasCalculatedWidth = true
				return
			}

			sidePadding := (gtx.Constraints.Width.Max - maxWidth) / 2
			gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
			layout.Center.Layout(gtx, func() {
				layout.Inset{
					Left:  unit.Dp(float32(sidePadding)),
					Right: unit.Dp(float32(sidePadding)),
				}.Layout(gtx, func() {
					(&layout.List{Axis: layout.Vertical, Alignment: layout.Middle}).Layout(gtx, len(widgetFuncs), func(i int) {
						gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
						fillMax(gtx, m.backgroundColor)
						layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
					})
				})
			})
		}),
	)
}

func (m *Modal) calculateTitleWidth(gtx *layout.Context) int {
	textLen := float32(len(m.titleLabel.Text)) * 2.5
	titleLines := m.titleLabel.shaper.LayoutString(m.titleLabel.Font, fixed.I(gtx.Px(unit.Sp(textLen))), gtx.Constraints.Width.Max, m.titleLabel.Text)
	return linesWidth(titleLines)
}

func (m *Modal) calulateButtonWidth(gtx *layout.Context, btn Button) int {
	textLen := float32(len(btn.Text)) * 2
	lines := btn.shaper.LayoutString(btn.Font, fixed.I(gtx.Px(unit.Sp(textLen))), gtx.Constraints.Width.Max, btn.Text)
	return linesWidth(lines) + 10 // 10 is to allow room for button padding
}

func linesWidth(lines []text.Line) int {
	var width fixed.Int26_6
	if len(lines) > 0 {
		for _, l := range lines {
			if l.Width > width {
				width = l.Width
			}
		}
	}
	return width.Ceil()
}
