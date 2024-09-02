package main

import (
  "crypto/sha256"
  "fmt"
)

// Pick color for sensor ID and return the color has a hex-formatted
// string.
//
// The color is chosen by hashing the sensor ID and using the first
// three bytes of digest to pick a value for each color channel from the
// following set of values: [0x00, 0x33, 0x66, 0x99, 0xcc, 0xff].
//
// This is based on the color palette from the 4th column of the
// following swatch:
//
// https://www.selecolor.com/en/recommended-color-palette/
//
// FIXME: There are better ways to pick colors that "match", but this is
// sufficient for now.
func defaultColor(id string) string {
  // get SHA256 digest of sensor ID (32 bytes)
  d := sha256.Sum256([]byte(id))

  return fmt.Sprintf("#%02x%02x%02x", (d[0]%6)*51, (d[1]%6)*51, (d[2]%6)*51)
}
