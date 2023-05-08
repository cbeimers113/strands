package game

import (
	"fmt"

	"github.com/g3n/engine/texture"
)

var Textures map[string]*texture.Texture2D

// Try to load the texture with a given name. Cache loaded textures into Textures buffer
func Texture(texID string) (tex *texture.Texture2D, ok bool) {
	tex, ok = Textures[texID]

	// If texture not found in buffer, load it from disk into buffer
	if !ok {
		tex, _ = texture.NewTexture2DFromImage(fmt.Sprintf("res/%s.png", texID))

		if tex != nil {
			ok = true
			Textures[texID] = tex
		} else {
			fmt.Printf("Could not load texture \"res/%s.png\"\n", texID)
		}
	}

	return
}
