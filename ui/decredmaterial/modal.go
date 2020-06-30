package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

type Modal struct {
	titleLabel     Label
	titleSeparator *Line

	overlayColor    color.RGBA
	backgroundColor color.RGBA
}

func (t *Theme) Modal() *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	return &Modal{
		titleLabel:     t.H6(""),
		titleSeparator: t.Line(),

		overlayColor:    overlayColor,
		backgroundColor: t.Color.Surface,
	}
}

func (m *Modal) Layout(gtx *layout.Context, title string, widgets []func()) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillMax(gtx, m.overlayColor)
		}),
		layout.Stacked(func() {
			widgetFuncs := []func(){
				func() {
					m.titleLabel.Text = title
					m.titleLabel.Layout(gtx)
				},
				func() {
					m.titleSeparator.Width = gtx.Constraints.Width.Max
					m.titleSeparator.Layout(gtx)
				},
			}
			widgetFuncs = append(widgetFuncs, widgets...)

			layout.Inset{
				Top:   unit.Dp(modalTopInset),
				Left:  unit.Dp(modalSideInset),
				Right: unit.Dp(modalSideInset),
			}.Layout(gtx, func() {
				(&layout.List{Axis: layout.Vertical, Alignment: layout.Middle}).Layout(gtx, len(widgetFuncs), func(i int) {
					gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
					fillMax(gtx, m.backgroundColor)
					layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
				})
			})
		}),
	)
}
