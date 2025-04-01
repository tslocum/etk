# etk - Ebitengine Tool Kit
[![GoDoc](https://codeberg.org/tslocum/godoc-static/raw/branch/main/badge.svg)](https://docs.rocket9labs.com/codeberg.org/tslocum/etk)
[![Donate via LiberaPay](https://img.shields.io/liberapay/receives/rocket9labs.com.svg?logo=liberapay)](https://liberapay.com/rocket9labs.com)

[Ebitengine](https://github.com/hajimehoshi/ebiten) tool kit for creating graphical user interfaces

**Note:** This library is still in development. Breaking changes may be made until v1.0 is released.

## Features

- Simplifies GUI development:
  - Propagates layout changes.
  - Propagates user input.
  - Propagates focus.
- Extensible by design:
  - The Box widget is provided as a building block for custom widgets.
  - Widgets may be nested within each other efficiently.
- Tools in the kit:
  - Box: Building block for creating custom widgets.
  - Button: Clickable button.
  - Flex: Flexible stack-based layout. Each Flex widget may be oriented horizontally or vertically.
  - Frame: Widget container. All child widgets are displayed at once. Child widgets are not repositioned by default.
  - Grid: Highly customizable cell-based layout. Each widget added to the Grid may span multiple cells.
  - Input: Text input widget. The Input widget is simply a Text widget that also accepts user input.
  - Keyboard: On-screen keyboard.
  - List: List of widgets as selectable items.
  - Select: Dropdown selection widget.
  - Sprite: Resizable image.
  - Text: Text display widget.
  - Window: Widget paging mechanism. Only one widget added to a window is displayed at a time.

## Demo

Browse the [widget showcase](https://rocketnine.itch.io/etk?secret=etk) using your browser. 

[Boxcars](https://codeberg.org/tslocum/boxcars) uses etk extensively and is available at https://play.bgammon.org

[![Screenshot](https://codeberg.org/tslocum/boxcars/raw/branch/main/screenshot.png)](https://codeberg.org/tslocum/boxcars/src/branch/main/screenshot.png)

## Examples

See the [examples](https://codeberg.org/tslocum/etk/src/branch/main/examples) folder.

## Documentation

Documentation is available via [godoc](https://docs.rocket9labs.com/codeberg.org/tslocum/etk).

## Support

Please share issues and suggestions [here](https://codeberg.org/tslocum/etk/issues).
