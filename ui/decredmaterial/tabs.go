package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// DefaultTabSizeVertical is the default flexed size of the tab section in a Tabs when vertically aligned
const DefaultTabSizeVertical = .15
// DefaultTabSizeHorizontal is the default flexed size of the tab section in a Tabs when horizontally aligned
const DefaultTabSizeHorizontal = .08

const (
	Top Position = iota
	Right
	Bottom
	Left
)

type Position int

type TabItem struct {
	Button
	index int
}

func tabIndicatorDimensions(gtx *layout.Context, tabPosition Position) (width, height int) {
	switch tabPosition {
	case Top, Bottom:
		width, height = gtx.Dimensions.Size.X, 4
	default:
		width, height = 5, gtx.Dimensions.Size.Y
	}
	return
}

// tabAlignment determines the alignment of the active tab indicator relative to the tab item
// content. It is determined by the position of the tab.
func tabAlignment(tabPosition Position) layout.Direction {
	switch tabPosition {
	case Top:
		return layout.S
	case Left:
		return layout.E
	case Bottom:
		return layout.N
	case Right:
		return layout.W
	default:
		return layout.E
	}
}

// indicator defines how the active tab indicator is drawn
func indicator(gtx *layout.Context, width, height int) layout.Widget {
	return func() {
		paint.ColorOp{Color: keyblue}.Add(gtx.Ops)
		paint.PaintOp{Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(width),
				Y: float32(height),
			},
		}}.Add(gtx.Ops)
		gtx.Dimensions = layout.Dimensions{
			Size: image.Point{X: width, Y: height},
		}
	}
}

func (t *TabItem) Layout(gtx *layout.Context, selected int, btn *widget.Button, tabPosition Position) {
	var tabWidth, tabHeight int

	layout.Stack{Alignment: tabAlignment(tabPosition)}.Layout(gtx,
		layout.Stacked(func() {
			if tabPosition == Left || tabPosition == Right {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			}
			t.Button.Color = darkblue
			t.Button.Background = color.RGBA{}
			t.Button.Layout(gtx, btn)
			tabWidth, tabHeight = tabIndicatorDimensions(gtx, tabPosition)
		}),
		layout.Stacked(func() {
			if selected != t.index {
				return
			}
			if tabPosition == Left || tabPosition == Right {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(0, indicator(gtx, tabWidth, tabHeight)))
			} else {
				indicator(gtx, tabWidth, tabHeight)()
			}
		}),
	)
}

// Tabs lays out a Flexed(Size) List with Selected as the first element and Item as the rest.
type Tabs struct {
	Flex     layout.Flex
	Size     float32
	items    []TabItem
	Selected int
	changed  bool
	btns     []*widget.Button
	list     *layout.List
	Position Position
}

func NewTabs() *Tabs {
	return &Tabs{
		list:     &layout.List{},
		Position: Left,
		Size:     DefaultTabSizeVertical,
	}
}

// SetTabs creates a button widget for each tab item
func (t *Tabs) SetTabs(tabs []TabItem) {
	t.items = tabs
	if len(t.items) != len(t.btns) {
		t.btns = make([]*widget.Button, len(t.items))
		for i := range t.btns {
			t.btns[i] = new(widget.Button)
		}
	}
}

// contentTabPosition depending on the specified tab position determines the order of the tab and
// the page content.
func (t *Tabs) contentTabPosition(gtx *layout.Context, body layout.Widget) (widgets []layout.FlexChild) {
	var content, tab layout.FlexChild

	widgets = make([]layout.FlexChild, 2)
	content = layout.Flexed(1-t.Size, func() {
		layout.Inset{Left: unit.Dp(5)}.Layout(gtx, body)
	})
	tab = layout.Flexed(t.Size, func() {
		t.list.Layout(gtx, len(t.btns), func(i int) {
			t.items[i].index = i
			t.items[i].Layout(gtx, t.Selected, t.btns[i], t.Position)
			if t.btns[i].Clicked(gtx) {
				t.Selected = i
			}
		})
	})

	switch t.Position {
	case Bottom, Right:
		widgets[0], widgets[1] = content, tab
	default:
		widgets[0], widgets[1] = tab, content
	}
	return widgets
}

// Layout the tabs
func (t *Tabs) Layout(gtx *layout.Context, body layout.Widget) {
	switch t.Position {
	case Top, Bottom:
		if t.Size < DefaultTabSizeHorizontal {
			t.Size = DefaultTabSizeHorizontal
		}
		t.list.Axis = layout.Horizontal
		t.Flex.Axis = layout.Vertical
	default:
		t.list.Axis = layout.Vertical
		t.Flex.Axis = layout.Horizontal
	}

	widgets := t.contentTabPosition(gtx, body)
	t.Flex.Layout(gtx, widgets...)
}
