package main

import (
   "fmt"
   "time"
   "sync"
   rl "github.com/lachee/raylib-goplus/raylib"
   hook "github.com/robotn/gohook"
)

const SCREEN_WIDTH = 800;
const SCREEN_HEIGHT = 800;

const PIXELS_X = 20;
const PIXELS_Y = 20;

const (
   Left = 65361
   Up = 65361 + iota
   Right = 65361 + iota
   Down = 65361 + iota
   Esc = 65307
)

type Snake []rl.Vector2

type Direction struct {
   mut sync.Mutex
   val rl.Vector2
}

func (this *Direction) Set(x, y float32) {
   this.mut.Lock()
   defer this.mut.Unlock()
   this.val = rl.Vector2{ X: x, Y: y }
}

func (this *Direction) Get() rl.Vector2 {
   this.mut.Lock()
   defer this.mut.Unlock()
   return this.val 
}

func NewDirection(x, y float32) Direction {
   return Direction{val: rl.NewVector2(x, y)}
}

func DrawPixel(pos rl.Vector2, pw int, ph int) {
   rl.DrawRectangle(int(pos.X) * pw, int(pos.Y) * ph, pw, ph, rl.White)
}

func DrawSnake(s Snake, pw int, ph int) {
   for _, pos := range s {
      DrawPixel(pos, pw, ph)
      //fmt.Printf("%v\n", pos)
   }
}

func (s *Snake) MoveSnake(dir rl.Vector2) {
   new_pos := (*s)[len(*s) - 1].Add(dir)
   if (new_pos.X < 0 || new_pos.Y < 0 || new_pos.X > PIXELS_X-1 || new_pos.Y > PIXELS_Y-1) {
      return
   }
   *s = (*s)[1:]
   *s = append(*s, (*s)[len(*s) - 1].Add(dir))
}

func main() {
   rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Snake")
   defer rl.CloseWindow()

   rl.SetTargetFPS(60)

   pixel_width := int(float32(SCREEN_WIDTH) / float32(PIXELS_X))
   pixel_height := int(float32(SCREEN_HEIGHT) / float32(PIXELS_Y))

   snake := make(Snake, 0, 50)
   snake = append(snake, rl.NewVector2(10, 10))
   snake = append(snake, rl.NewVector2(11, 10))
   snake = append(snake, rl.NewVector2(12, 10))

   dir := Direction{val: rl.NewVector2(1, 0)}

   s := hook.Start()
   defer hook.End()

   fmt.Printf("%v %v %v %v\n", Up, Left, Down, Right)

   go func() {
      for input := range s {
         if input.Kind == hook.KeyDown {
            fmt.Printf("%v", input.Rawcode)
            switch input.Rawcode {
            case Up:
               dir.Set(0, -1)
            case Down:
               dir.Set(0, 1)
            case Left:
               dir.Set(-1, 0)
            case Right:
               dir.Set(1, 0)
            }
         }
      }
   }()

   for !rl.WindowShouldClose() {
      rl.BeginDrawing()

      rl.ClearBackground(rl.Black)

      DrawSnake(snake, pixel_width, pixel_height)

      snake.MoveSnake(dir.Get())

      time.Sleep(500 * time.Millisecond)

      rl.EndDrawing()
   }
}
