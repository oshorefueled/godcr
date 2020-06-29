package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

type Modal struct {
	titleLabel     Label
	titleSeparator *Line

	contentHeight      int
	hasCalculatedeight bool

	overlayColor    color.RGBA
	backgroundColor color.RGBA

	widgetItemPadding float32
}

func (t *Theme) Modal() *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	return &Modal{
		titleLabel:     t.H6(""),
		titleSeparator: t.Line(),

		contentHeight:      0,
		hasCalculatedeight: false,

		overlayColor:    overlayColor,
		backgroundColor: t.Color.Surface,

		widgetItemPadding: 20,
	}
}

func (m *Modal) calculateWidgetHeight(gtx *layout.Context) int {
	return gtx.Dimensions.Size.Y + int(m.widgetItemPadding)
}

func (m *Modal) Layout(gtx *layout.Context, title string, widgets []func()) {
	contentHeight := 0
	maxHeight := gtx.Dimensions.Size.Y

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillMax(gtx, m.overlayColor)
		}),
		layout.Stacked(func() {
			widgetsFuncs := []func(){
				func() {
					m.titleLabel.Text = title
					m.titleLabel.Layout(gtx)
					contentHeight += m.calculateWidgetHeight(gtx)
				},
				func() {
					m.titleSeparator.Width = gtx.Constraints.Width.Max
					m.titleSeparator.Layout(gtx)
					contentHeight += m.calculateWidgetHeight(gtx)
				},
			}

			for i := range widgets {
				index := i
				widgetsFuncs = append(widgetsFuncs, func() {
					widgets[index]()
					contentHeight += m.calculateWidgetHeight(gtx)
				})
			}

			var inset layout.Inset
			if !m.hasCalculatedeight {
				inset = layout.Inset{}
				inset.Layout(gtx, func() {
					fillMax(gtx, m.backgroundColor)
					(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(widgetsFuncs), func(i int) {
						layout.UniformInset(unit.Dp(m.widgetItemPadding/2)).Layout(gtx, widgetsFuncs[i])
					})
				})
				if contentHeight != 0 {
					m.contentHeight = contentHeight
				}
				m.hasCalculatedeight = true
			} else {
				inset = layout.Inset{
					Top:    unit.Dp(modalTopInset),
					Bottom: unit.Dp(float32(maxHeight - m.contentHeight - modalTopInset - 20)),
					Left:   unit.Dp(modalSideInset),
					Right:  unit.Dp(modalSideInset),
				}
				inset.Layout(gtx, func() {
					fillMax(gtx, m.backgroundColor)
					(&layout.List{Axis: layout.Vertical, Alignment: layout.Middle}).Layout(gtx, len(widgetsFuncs), func(i int) {
						layout.UniformInset(unit.Dp(m.widgetItemPadding/2)).Layout(gtx, widgetsFuncs[i])
					})
				})
			}
		}),
	)
}
