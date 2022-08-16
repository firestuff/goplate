package main

import (
  "fmt"
  "image"
  "image/color"
  "image/png"
  "math"
	"net/http"
  "os"

  _ "image/jpeg"
)

func max16(a, b uint16) uint16 {
  if a > b {
    return a
  } else {
    return b
  }
}

func max32(a, b uint32) uint32 {
  if a > b {
    return a
  } else {
    return b
  }
}

func maxDiff(a, d1, d2 uint32) uint32 {
  var d uint32
  if d1 > d2 {
    d = d1 - d2
  } else {
    d = d2 - d1
  }

  return max32(a, d)
}

func detectEdges(src image.Image) *image.Gray16 {
  bounds := src.Bounds()
  ret := image.NewGray16(bounds)

  for y := bounds.Min.Y + 1; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
      inr, ing, inb, _ := src.At(x, y).RGBA()
      tempr, tempg, tempb, _ := src.At(x - 1, y - 1).RGBA()

      var c uint32
      c = maxDiff(c, inr, tempr)
      c = maxDiff(c, ing, tempg)
      c = maxDiff(c, inb, tempb)

      outc := color.Gray16{Y: uint16(c)}
      ret.Set(x, y, outc)
    }
  }

  return ret
}

func loadImage(path string) (image.Image, error) {
  fh, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer fh.Close()

  img, _, err := image.Decode(fh)
  if err != nil {
    return nil, err
  }

  return img, nil
}

func normalize(src *image.Gray16) *image.Gray16 {
  bounds := src.Bounds()
  var m uint16

  for y := bounds.Min.Y + 1; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
      m = max16(m, src.Gray16At(x, y).Y)
    }
  }

  scale := math.MaxUint16 / m
  ret := image.NewGray16(bounds)

  for y := bounds.Min.Y + 1; y < bounds.Max.Y; y++ {
    for x := bounds.Min.X + 1; x < bounds.Max.X; x++ {
      c := src.Gray16At(x, y)
      c.Y *= scale
      ret.Set(x, y, c)
    }
  }

  return ret
}

func find(shape, img *image.Gray16) {
  sb := shape.Bounds()
  imgb := img.Bounds()

  for y := imgb.Min.Y; y < imgb.Max.Y - (sb.Max.Y - sb.Min.Y); y++ {
    for x := imgb.Min.X; x < imgb.Max.X - (sb.Max.X - sb.Min.X); x++ {
      score(shape, img, x, y)
    }
  }
}

func score(shape, img *image.Gray16, ox, oy int) {
  sb := shape.Bounds()
  imgb := img.Bounds()

  for y := sb.Min.Y; y < sb.Max.Y; y++ {
    for x := sb.Min.X; x < sb.Max.X; x++ {
      // imgy := oy + y
      //

    }
  }
}

func test(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "image/png")

  t, err := loadImage("shapes/T.png")
  if err != nil {
    panic(err)
  }
  te := detectEdges(t)

  src, err := loadImage("samples/big.png")
  if err != nil {
    panic(err)
  }
  srce := detectEdges(src)

  find(te, srce)

  err = png.Encode(w, srce)
  if err != nil {
    panic(err)
  }
}

func main() {
  fmt.Printf("starting...\n")
	http.HandleFunc("/test", test)
	http.ListenAndServe(":8090", nil)
}
