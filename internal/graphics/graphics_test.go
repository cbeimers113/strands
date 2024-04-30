package graphics

import (
	"errors"
	"testing"

	"github.com/g3n/engine/texture"
	"github.com/stretchr/testify/assert"
)

func Test_Texture(t *testing.T) {
	mockId := "mock"
	mockTexture := new(texture.Texture2D)
	textures = map[string]*texture.Texture2D{mockId: mockTexture}

	tests := []struct {
		name  string
		texId string
		want  *texture.Texture2D
		err   error
	}{
		{
			name:  "Happy path",
			texId: mockId,
			want:  mockTexture,
		},
		{
			name:  "Sad path - texture doesn't exist",
			texId: "no_tex",
			err:   errors.New("texture does not exist: [no_tex]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Texture(tt.texId)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}
