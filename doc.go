/*
Package etk provides an Ebitengine tool kit for creating graphical user interfaces.

# Widgets

Custom widgets may be created entirely from scratch or may be
based on official widgets.

The following official widgets are available:

	Box - Building block for creating other widgets.
	Button - Clickable button.
	Flex - Flexible stack-based layout. Each Flex widget may be oriented horizontally or vertically.
	Frame - Widget container. All child widgets are displayed at once. Child widgets are not repositioned by default.
	Grid - Highly customizable cell-based layout. Widgets added to the Grid may span multiple cells.
	Input - Text input widget. The Input widget is simply a Text widget that also accepts user input.
	List - List of widgets as selectable items.
	Text - Text display widget.
	Window - Widget paging mechanism. Only one widget added to a window is displayed at a time.

# Input Propagation

Mouse events are passed to the topmost widget under the mouse. If a widget
returns a handled value of false, the event continues to propagate down the
stack of widgets under the mouse.

Clicking or tapping on a widget focuses the widget. This is handled by etk
automatically when a widget returns a handled value of true.

Keyboard events are passed to the focused widget.

# Focus Propagation

When attempting to change which widget is focused, etk checks whether the widget
to be focused accepts this focus. If it does, the previously focused widget is
un-focused. If the widget does not accept the focus, the previously focused
widget remains focused.

# Cursor Unification

Input events generated by desktop mice and touch screens are unified in etk.
These input events are simplified into an image.Point specifying the location
of the event and two parameters: clicked and pressed.

Clicked is true the first frame the mouse event or touch screen event is received.
When the mouse click or touch screen tap is released, the widget that was originally
clicked or tapped always receives a final event where clicked and pressed are both false.

# Draw Order

Each time etk draws a widget it subsequently draws all of the widget's children
in the order they are returned.

# Subpackages

There are two subpackages in etk: messeji and kibodo. These are available for
use without requiring etk. Usually you will not reference any subpackages, as
etk wraps them to provide widgets with additional features.
*/
package etk
