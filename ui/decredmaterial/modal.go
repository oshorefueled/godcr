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

	list            *layout.List
	overlayColor    color.RGBA
	backgroundColor color.RGBA
	button          *widget.Button
}

func (t *Theme) Modal(title string) *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	return &Modal{
		titleLabel:     t.H6(title),
		titleSeparator: t.Line(),

		list:            &layout.List{Axis: layout.Vertical, Alignment: layout.Middle},
		overlayColor:    overlayColor,
		backgroundColor: t.Color.Surface,
		button:          new(widget.Button),
	}
}

func (m *Modal) SetTitle(title string) {
	m.titleLabel.Text = title
}

// Layout lays out the widget passed as an argument in a modal using a specified
// margin. Its left and right margin are respect to 3840 resolution.
func (m *Modal) Layout(gtx *layout.Context, widgets []func(), margin int) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillMax(gtx, m.overlayColor)
			m.button.Layout(gtx)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
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
			scaled := 3840 / float32(gtx.Constraints.Width.Max)
			mg := unit.Px(float32(margin) / scaled)

			layout.Center.Layout(gtx, func() {
				layout.Inset{
					Left:  mg,
					Right: mg,
				}.Layout(gtx, func() {
					m.list.Layout(gtx, len(widgetFuncs), func(i int) {
						gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
						fillMax(gtx, m.backgroundColor)
						layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
					})
				})
			})
		}),
	)
}
