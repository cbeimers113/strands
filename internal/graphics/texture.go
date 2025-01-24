package graphics

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/draw"

	"github.com/g3n/engine/texture"
)

const (
	TexMenuLogo = "menuLogo"
	TexCursor   = "cursor"
	TexDirt     = "dirt"
	TexGrass    = "grass"
	TexSand     = "sand"
	TexSeed     = "seed"
	TexStalk    = "stalk"
	TexStone    = "stone"
	TexWater    = "water"
)

var (
	//go:embed textures/menuLogo.png
	bytesMenuLogo []byte
	//go:embed textures/cursor.png
	bytesCursor []byte
	//go:embed textures/dirt.png
	bytesDirt []byte
	//go:embed textures/grass.png
	bytesGrass []byte
	//go:embed textures/sand.png
	bytesSand []byte
	//go:embed textures/seed.png
	bytesSeed []byte
	//go:embed textures/stalk.png
	bytesStalk []byte
	//go:embed textures/stone.png
	bytesStone []byte
	//go:embed textures/water.png
	bytesWater []byte
)

var Textures map[string]*texture.Texture2D

func init() {
	Textures = make(map[string]*texture.Texture2D)

	for _, texLoader := range []struct {
		id   string
		data []byte
	}{
		{id: TexMenuLogo, data: bytesMenuLogo},
		{id: TexCursor, data: bytesCursor},
		{id: TexDirt, data: bytesDirt},
		{id: TexGrass, data: bytesGrass},
		{id: TexSand, data: bytesSand},
		{id: TexSeed, data: bytesSeed},
		{id: TexStalk, data: bytesStalk},
		{id: TexStone, data: bytesStone},
		{id: TexWater, data: bytesWater},
	} {
		var (
			tex *texture.Texture2D
			err error
		)

		if tex, err = decode(texLoader.data); err == nil {
			Textures[texLoader.id] = tex
		} else {
			panic(fmt.Errorf("error loading texture [%s]: %w", texLoader.id, err))
		}
	}
}

func decode(data []byte) (tex *texture.Texture2D, err error) {
	var img image.Image
	if img, _, err = image.Decode(bytes.NewReader(data)); err != nil {
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		err = errors.New("unsupported stride")
		return
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	tex = texture.NewTexture2DFromRGBA(rgba)
	return
}

// Texture returns the texture for a given key if it exists, and errors if it doesn't
func Texture(texId string) (tex *texture.Texture2D, err error) {
	var ok bool
	if tex, ok = Textures[texId]; !ok {
		err = fmt.Errorf("texture does not exist: [%s]", texId)
	}

	return
}
