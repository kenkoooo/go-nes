package ui

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-gl/gl/v2.1/gl"
)

const textureSize = 4096
const textureDim = textureSize / 256
const textureCount = textureDim * textureDim

type Texture struct {
	texture uint32
	lookup  map[string]int
	reverse [textureCount]string
	access  [textureCount]int
	counter int
	ch      chan string
}

func NewTexture() *Texture {
	texture := createTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		textureSize, textureSize,
		0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	t := Texture{}
	t.texture = texture
	t.lookup = make(map[string]int)
	t.ch = make(chan string, 1024)
	return &t
}

func (t *Texture) Purge() {
	for {
		select {
		case path := <-t.ch:
			delete(t.lookup, path)
		default:
			return
		}
	}
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.texture)
}

func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) mark(index int) {
	t.counter++
	t.access[index] = t.counter
}

func (t *Texture) lru() int {
	minIndex := 0
	minValue := t.counter + 1
	for i, n := range t.access {
		if n < minValue {
			minIndex = i
			minValue = n
		}
	}
	return minIndex
}

func (t *Texture) coord(index int) (x, y, dx, dy float32) {
	x = float32(index%textureDim) / textureDim
	y = float32(index/textureDim) / textureDim
	dx = 1.0 / textureDim
	dy = dx * 240 / 256
	return
}

func (t *Texture) downloadThumbnail(romPath, hash string) error {
	url := thumbnailURL(hash)
	filename := thumbnailPath(hash)
	dir, _ := path.Split(filename)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	t.ch <- romPath

	return nil
}
