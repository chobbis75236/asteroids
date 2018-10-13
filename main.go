package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

// Global constants.
const screenWidth = 1024
const screenHeight = 768

// Global variables.
var (
	windowTitlePrefix	= "Asteroids"
	frames				= 0
	second				= time.Tick(time.Second)
	window				*pixelgl.Window
	frameLength			float64
	shipPic				pixel.Picture
	asteroidPic			pixel.Picture
	projectilePic		pixel.Picture
	es					[]entity // short for entity slice
	lastFire			= time.Now()
	weapon				= Gun
)

/* ENTITY STRUCTURE */

// We can refer to entity type by name rather than just a number for readability.
type etype int
const (
	Ship etype = 0
	Asteroid etype = 1
	Projectile etype = 2
)

// Weapon type
type wtype int
const (
	Gun wtype = 0
	Flames wtype = 1
)

// All the information needed for every entity.
type entity struct {
	etype
	x, y, angle, dx, dy, dangle, radius float64
	sprite *pixel.Sprite // * refers to a pointer to the sprite, not a copy.
}

/* FUNCTIONS */

func distance(e1, e2 entity) float64 {
	return math.Hypot(e2.x - e1.x, e2.y - e1.y)
}

func (e1 entity) intersects(e2 entity) bool {
	return math.Pow(e2.x - e1.x, 2) + math.Pow(e2.y - e1.y, 2) <= math.Pow(e2.radius + e1.radius, 2)
}

func (e entity) velocity() float64 {
	return math.Hypot(e.dx, e.dy)
}

// Returns an image from a path.
func loadImageFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Will close the file once the function returns a value.
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func initiate() {

	var initError error

	// Define the settings for the window.
	cfg := pixelgl.WindowConfig{
		Bounds:	pixel.R(0, 0, screenWidth, screenHeight),
		VSync:	true, // Makes time between frames more even.
	}

	// Make the window.
	window, initError = pixelgl.NewWindow(cfg)
	if initError != nil {
		panic(initError)
	}

	// Load the images used into the game.
	shipImage, initError := loadImageFile("ship.png")
	if initError != nil {
		panic(initError)
	}
	shipPic = pixel.PictureDataFromImage(shipImage)

	asteroidImage, initError := loadImageFile("asteroid.png")
	if initError != nil {
		panic(initError)
	}
	asteroidPic = pixel.PictureDataFromImage(asteroidImage)

	projectileImage, initError := loadImageFile("projectile.png")
	if initError != nil {
		panic(initError)
	}
	projectilePic = pixel.PictureDataFromImage(projectileImage)

	// Initiate entity slice by adding the ship.
	es = []entity{
		{
			etype:	Ship,
			x:		float64(screenWidth / 2),
			y:		float64(screenHeight / 2),
			angle:	0.0,
			dx:		0,
			dy:		0,
			dangle:	0.0,
			radius:	30,
			sprite:	pixel.NewSprite(shipPic, shipPic.Bounds()),
		},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialise random variable.

	for i := 0; i < 20; i++ {
		e := entity{
			etype:	Asteroid,
			x:		r.Float64() * screenWidth,
			y:		r.Float64() * screenHeight,
			angle:	r.Float64() * 2 * math.Pi,
			dx:		r.Float64() * 100 - 50,
			dy:		r.Float64() * 100 - 50,
			dangle:	r.Float64() * 2 - 1,
			radius:	r.Float64() * 20 + 20,
			sprite:	pixel.NewSprite(asteroidPic, asteroidPic.Bounds()),
		}

		es = append(es, e)
	}

}

func game() {

	initiate()

	// MAIN GAME LOOP
	for !window.Closed() {

		frameStart := time.Now()

		// KEY DETECTION
		a := 15.0
		if window.Pressed(pixelgl.KeyLeft) {
			es[0].dangle += 0.1
		}
		if window.Pressed(pixelgl.KeyRight) {
			es[0].dangle -= 0.1
		}
		if window.Pressed(pixelgl.KeyUp) || window.Pressed(pixelgl.KeyW) {
			es[0].dx -= a * math.Sin(es[0].angle)
			es[0].dy += a * math.Cos(es[0].angle)
		}
		if window.Pressed(pixelgl.KeyDown) || window.Pressed(pixelgl.KeyS) {
			es[0].dx += a * math.Sin(es[0].angle)
			es[0].dy -= a * math.Cos(es[0].angle)
		}
		if window.Pressed(pixelgl.KeyA) {
			es[0].dx -= a * math.Cos(es[0].angle)
			es[0].dy -= a * math.Sin(es[0].angle)
		}
		if window.Pressed(pixelgl.KeyD) {
			es[0].dx += a * math.Cos(es[0].angle)
			es[0].dy += a * math.Sin(es[0].angle)
		}
		if window.Pressed(pixelgl.KeySpace) {

			var fireRate float64
			if weapon == Gun {
				fireRate = 0.2
			} else if weapon == Flames {
				fireRate = 0.0
			}

			if time.Since(lastFire).Seconds() > fireRate {
				lastFire = time.Now()

				projDx := -math.Sin(es[0].angle)
				projDy := math.Cos(es[0].angle)

				es = append(es, entity{
					etype:	Projectile,
					x:		es[0].x + es[0].radius*projDx,
					y:		es[0].y + es[0].radius*projDy,
					angle:	es[0].angle,
					dx:		500 * projDx,
					dy:		500 * projDy,
					dangle:	0.0,
					radius: es[0].radius/3,
					sprite:	pixel.NewSprite(projectilePic, projectilePic.Bounds()),
				})
			}
		}
		if window.JustPressed(pixelgl.KeyLeftShift) || window.Pressed(pixelgl.KeyQ) {
			weapon = (weapon - 1) % 2
			fmt.Println(weapon)
			fmt.Println("SWITCHED DOWN")
		}
		if window.JustPressed(pixelgl.KeyLeftControl) || window.Pressed(pixelgl.KeyE) {
			weapon = (weapon + 1) % 2
			fmt.Println(weapon)
			fmt.Println("SWITCHED UP")
		}

		// ENTITY COLLISION HANDLER
		newAsteroids := make([]entity, 0) // Since more than one asteroid may be added.
		for i := 0; i < len(es); {

			removeI := false

			for j := 1; j < len(es); j++ {

				// If colliding with itself or a projectile colliding with ship or asteroid, ignore.
				if i == j || es[j].etype == Projectile && (i == 0 || es[i].etype == Asteroid) {
					continue
				}

				if es[i].intersects(es[j]) {

					if es[i].etype == Projectile && es[j].etype == Asteroid {

						removeI = true

						es[j].radius /= math.Sqrt2

						newAsteroids = append(newAsteroids, entity{
							etype:	Asteroid,
							x:		es[j].x,
							y:		es[j].y,
							angle:	es[j].angle,
							dx:		-es[j].dx,
							dy:		-es[j].dy,
							dangle:	-es[j].dangle,
							radius:	es[j].radius,
							sprite:	pixel.NewSprite(asteroidPic, asteroidPic.Bounds()),
						})

					}

					d := distance(es[i], es[j])
					unitX := (es[j].x - es[i].x)/d
					unitY := (es[j].y - es[i].y)/d

					v1 := es[i].velocity()
					v2 := es[j].velocity()

					es[i].dx = -v2 * unitX
					es[i].dy = -v2 * unitY

					es[j].dx = v1 * unitX
					es[j].dy = v1 * unitY

				}

			}

			if removeI {
				es = append(es[:i], es[i+1:]...) // Removes es[i].
			} else {
				i++
			}
		}

		// Adding new asteroids.
		es = append(es, newAsteroids...)

		// Removing Asteroids which are too small.
		for i := 0; i < len(es); {
			if es[i].etype == Asteroid && es[i].radius < 10 {
				es = append(es[:i], es[i+1:]...)
			} else {
				i++
			}
		}


		// ENTITY POSITION UPDATE LOOP
		for i := range es {

			// Reducing velocity over time if not increasing.
			if es[i].etype == Ship {
				// If not accelerating reduce velocity.
				if !window.Pressed(pixelgl.KeyUp) && !window.Pressed(pixelgl.KeyDown) && !window.Pressed(pixelgl.KeyW) &&
					!window.Pressed(pixelgl.KeyA) && !window.Pressed(pixelgl.KeyS) && !window.Pressed(pixelgl.KeyD){
					es[i].dx *= 1 - frameLength
					es[i].dy *= 1 - frameLength
				}

				// If not angularly accelerating, reduce angular velocity.
				if !window.Pressed(pixelgl.KeyLeft) && !window.Pressed(pixelgl.KeyRight) {
					es[i].dangle *= 1 - frameLength
				}
			}

			es[i].x += es[i].dx * frameLength
			es[i].y += es[i].dy * frameLength
			es[i].angle += es[i].dangle * frameLength

			if es[i].x < -50 {
				es[i].x += screenWidth + 100
			}
			if es[i].y < -50 {
				es[i].y += screenHeight + 100
			}
			if es[i].x > screenWidth + 50 {
				es[i].x -= screenWidth + 100
			}
			if es[i].y > screenHeight + 50 {
				es[i].y -= screenHeight + 100
			}
		}

		window.Clear(colornames.Black) // Fill window with black.

		/* BEGIN DRAW LOOP */

		for i := range es {
			scale := 2*es[i].radius/((es[i].sprite.Picture().Bounds().W()+es[i].sprite.Picture().Bounds().H())/2)
			matrix := pixel.IM.
				Rotated(pixel.ZV, es[i].angle).
				Scaled(pixel.ZV, scale).
				Moved(pixel.Vec{X: es[i].x, Y: es[i].y})
			es[i].sprite.Draw(window, matrix)
		}

		/* END DRAW LOOP */

		window.Update() // Draw contents of window to the screen.

		frames++
		select {
		case <- second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}

		frameLength = time.Since(frameStart).Seconds()

	}
}

func main() {

	pixelgl.Run(game)

}
