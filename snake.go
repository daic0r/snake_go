package main

import (
   "fmt"
   "time"
   "sync"
   "math/rand"
   "slices"
   rl "github.com/lachee/raylib-goplus/raylib"
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

type Tile int
const (
   Empty = Tile(iota)
   Food = Tile(iota)
)

type PlayingField [][]Tile
func NewPlayingField(w, h int) PlayingField {
   ret := make([][]Tile, h)
   for y := 0; y < h; y++ {
      ret[y] = make([]Tile, w)
   }
   return ret
}

type GameState struct {
   alive bool
   duration int
   occupied_fields map[rl.Vector2]Tile
   snake Snake
   score int
}

func (this *GameState) CheckFoodEaten() {
   for pos := range this.occupied_fields {
      last := this.snake[len(this.snake) - 1]
      if pos == last {
         delete(this.occupied_fields, pos)
         append_dir := last.Subtract(this.snake[len(this.snake) - 2])
         fmt.Printf("%v\n", append_dir)
         this.snake = append(this.snake, last.Add(append_dir))
         this.PlaceFood()
         this.score += 100
      }
   }
}

func (this *GameState) PlaceFood() {
   food_y := rand_gen.Intn(PIXELS_Y)
   food_x := rand_gen.Intn(PIXELS_X)
   this.occupied_fields[rl.NewVector2(float32(food_x), float32(food_y))] = Food
}

func NewGame() GameState {
   ret := GameState{ alive: true, duration: 400, occupied_fields: make(map[rl.Vector2]Tile), snake: make(Snake, 0, 50) }
   ret.snake = append(ret.snake, rl.NewVector2(10, 10))
   ret.snake = append(ret.snake, rl.NewVector2(11, 10))
   ret.snake = append(ret.snake, rl.NewVector2(12, 10))
   ret.PlaceFood()
   return ret
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

func DrawPixel(pos rl.Vector2, pw int, ph int, color rl.Color) {
   rl.DrawRectangle(int(pos.X) * pw, int(pos.Y) * ph, pw, ph, color)
}

func DrawSnake(s Snake, pw int, ph int) {
   for _, pos := range s {
      DrawPixel(pos, pw, ph, rl.White)
      //fmt.Printf("%v\n", pos)
   }
}

// @return: true => legal move
//          false => player dies
func (s *Snake) MoveSnake(dir rl.Vector2) (bool) {
   new_pos := (*s)[len(*s) - 1].Add(dir)
   if (new_pos.X < 0 || new_pos.Y < 0 || new_pos.X > PIXELS_X-1 || new_pos.Y > PIXELS_Y-1) {
      return false
   }
   new_point := (*s)[len(*s) - 1].Add(dir)
   if slices.Index(*s, new_point) != -1 {
      return false
   }
   *s = (*s)[1:]
   *s = append(*s, new_point)
   return true
}

var rand_gen = rand.New(rand.NewSource(time.Now().Unix()))

func main() {
   rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Snake")
   defer rl.CloseWindow()

   rl.SetTargetFPS(60)

   pixel_width := int(float32(SCREEN_WIDTH) / float32(PIXELS_X))
   pixel_height := int(float32(SCREEN_HEIGHT) / float32(PIXELS_Y))

   // snake := make(Snake, 0, 50)
   // snake = append(snake, rl.NewVector2(10, 10))
   // snake = append(snake, rl.NewVector2(11, 10))
   // snake = append(snake, rl.NewVector2(12, 10))

   //occupied_fields := make(map[rl.Vector2]Tile)

   state := NewGame()

   dir := Direction{val: rl.NewVector2(1, 0)}

   fmt.Printf("%v %v %v %v\n", Up, Left, Down, Right)

   for !rl.WindowShouldClose() {
      if rl.IsKeyDown(rl.KeyUp) {
         dir.Set(0, -1)
      } else
      if rl.IsKeyDown(rl.KeyDown) {
         dir.Set(0, 1)
      } else
      if rl.IsKeyDown(rl.KeyLeft) {
         dir.Set(-1, 0)
      } else
      if rl.IsKeyDown(rl.KeyRight) {
         dir.Set(1, 0)
      }

      state.alive = state.snake.MoveSnake(dir.Get())
      state.CheckFoodEaten()

      rl.BeginDrawing()

         rl.ClearBackground(rl.Black)

         DrawSnake(state.snake, pixel_width, pixel_height)
         for pos, t := range state.occupied_fields {
            switch t {
            case Food:
               DrawPixel(pos, pixel_width, pixel_height, rl.Green)
            default:
               DrawPixel(pos, pixel_width, pixel_height, rl.Red)
            }
         }

         rl.DrawText(fmt.Sprintf("Score: %d", state.score), 10, 10, 30, rl.Red)

      rl.EndDrawing()

      time.Sleep(time.Duration(state.duration) * time.Millisecond)
   }
}
