package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Modal struct {
	titleLabel     Label
	titleSeparator *Line

	overlayColor    color.RGBA
	backgroundColor color.RGBA
}

func (t *Theme) Modal(title string) *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	return &Modal{
		titleLabel:     t.H6(title),
		titleSeparator: t.Line(),

		overlayColor:    overlayColor,
		backgroundColor: t.Color.Surface,
	}
}

func (m *Modal) SetTitle(title string) {
	m.titleLabel.Text = title
}

// Layout lays out the modal with specified title and width
// If the passed width is 0, then a default width is used for the modal
// If not, the modal assumes the width passed to it
func (m *Modal) Layout(gtx *layout.Context, widgets []func(), maxWidth int) {
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
			widgetFuncs = append(widgetFuncs, widgets...)

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
