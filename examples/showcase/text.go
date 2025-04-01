//go:build example

package main

import (
	"codeberg.org/tslocum/etk"
)

func newTextExample() (etk.Widget, etk.Widget) {
	text := etk.NewText(loremIpsum)
	text.SetPadding(etk.Scale(10))
	text.SetFollow(false)

	return text, nil
}

const loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam pellentesque lorem eu mauris feugiat, vel posuere nibh lobortis. Nam eget elit vitae arcu maximus fringilla at ut nisl. Curabitur rutrum a est ac cursus. Quisque sed sodales libero, ut faucibus augue. Fusce eros magna, porttitor maximus ante at, vestibulum consectetur turpis. Vivamus placerat purus sit amet vestibulum sodales. Vivamus enim lacus, ultricies pharetra venenatis venenatis, volutpat vitae magna. Nullam aliquam orci at ipsum accumsan hendrerit. Duis tincidunt aliquet augue, sed efficitur ex. Cras euismod eget nisi sit amet molestie. Praesent ligula erat, rutrum quis ultrices at, mattis vitae diam. Praesent porta felis justo, et varius lacus blandit vel. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos.

Quisque mattis pulvinar porttitor. Nulla pulvinar justo sed sapien interdum posuere. Morbi id dapibus massa. Aenean non sapien egestas, malesuada velit finibus, tristique mi. Donec facilisis erat urna, in eleifend odio tincidunt nec. Aenean egestas ante mauris. In id euismod turpis, ac scelerisque ipsum. Pellentesque dapibus eget est at viverra. Aliquam volutpat semper magna. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Curabitur efficitur blandit felis at lacinia. Duis luctus felis id lectus suscipit, in faucibus enim volutpat. Curabitur et ipsum eu enim sagittis vulputate. Phasellus ac feugiat lacus.

Mauris non lacinia nibh, vitae rutrum est. Integer auctor eleifend nulla in laoreet. Duis non tellus cursus, ullamcorper tellus a, sollicitudin felis. Duis porta libero nec congue dictum. Aenean felis dui, pharetra in justo vel, laoreet interdum risus. Donec ultrices posuere sapien, eu malesuada lorem vulputate vel. Quisque vehicula tortor mattis, molestie mauris ac, euismod urna. Curabitur maximus velit blandit, finibus felis ac, congue nunc. Cras sapien tortor, scelerisque non tellus sed, faucibus dictum libero. Donec vitae porttitor magna. Curabitur at lobortis arcu, et pellentesque ex. Nunc eu metus sit amet erat sollicitudin ultricies. Curabitur quam diam, vulputate et lorem a, auctor sodales enim. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Suspendisse finibus erat eu odio placerat hendrerit.`
