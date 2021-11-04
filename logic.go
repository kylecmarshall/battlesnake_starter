package main

// This file can be a nice home for your Battlesnake logic and related helper functions.
//
// We have started this for you, with a function to help remove the 'neck' direction
// from the list of possible moves!

import (
	"log"
	"math/rand"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
// TIP: If you open your Battlesnake URL in browser you should see this data.
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "kylecmarshall", // TODO: Your Battlesnake username
		Color:      "#023047",       // TODO: Personalize
		Head:       "pixel",         // TODO: Personalize
		Tail:       "pixel",         // TODO: Personalize
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// Step 0: Don't let your Battlesnake move back in on it's own neck
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of body piece directly behind your head (your "neck")
	if myNeck.X < myHead.X {
		possibleMoves["left"] = false
	} else if myNeck.X > myHead.X {
		possibleMoves["right"] = false
	} else if myNeck.Y < myHead.Y {
		possibleMoves["down"] = false
	} else if myNeck.Y > myHead.Y {
		possibleMoves["up"] = false
	}

	// TODO: Step 1 - Don't hit walls.
	// Use information in GameState to prevent your Battlesnake
	// from moving beyond the boundaries of the board.
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	if myHead.X == 0 {
		possibleMoves["left"] = false
	} else if myHead.X == boardWidth-1 {
		possibleMoves["right"] = false
	}

	if myHead.Y == 0 {
		possibleMoves["down"] = false
	} else if myHead.Y == boardHeight-1 {
		possibleMoves["up"] = false
	}

	// TODO: Step 2 - Don't hit yourself.
	// Use information in GameState to prevent your Battlesnake
	// from colliding with itself.
	myBody := state.You.Body
	for _, move := range [4]string{"up", "down", "left", "right"} {
		possibleCoord := getCoordFromMove(state, move)

		for _, bodyPart := range myBody {
			if bodyPart == possibleCoord {
				log.Printf("%s Body collision prevent moving %s\n",
					state.Game.ID,
					move)
				possibleMoves[move] = false
				break
			}
		}
	}

	// TODO: Step 3 - Don't collide with others.
	// Use information in GameState to prevent your Battlesnake
	// from colliding with others.
	for _, move := range [4]string{"up", "down", "left", "right"} {
		possibleCoord := getCoordFromMove(state, move)

		for _, snake := range state.Board.Snakes {

			for _, bodyPart := range snake.Body {
				if bodyPart == possibleCoord {
					log.Printf("%s Snake collision prevent moving %s\n",
						state.Game.ID,
						move)
					possibleMoves[move] = false
					break
				}
			}
		}
	}

	mode := "default"
	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.
	if state.You.Health < int32(boardHeight-1+boardWidth-1) {
		mode = "starving"
	} /* else if bodyPartsOnDiagonals(state) >= 2 {
		mode = "scared"
	} */

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string

	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
	} else {
		switch mode {
		case "scared":
			nextMove = scared_pickMove(state, safeMoves)
      break
		case "starving":
			nextMove = starving_pickMove(state, safeMoves)
			break
		case "default":
			nextMove = default_pickMove(state, safeMoves)
			break
		default:
			nextMove = default_pickMove(state, safeMoves)
			break
		}
		log.Printf("%s MODE %s MOVE %d: %s\n",
			state.Game.ID,
			mode,
			state.Turn,
			nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}

//pickMove
func default_pickMove(state GameState, safeMoves []string) string {
	// take a random walk, snake!
	return safeMoves[rand.Intn(len(safeMoves))]
}

func starving_pickMove(state GameState, safeMoves []string) string {
	minDistToFood := state.Board.Height + state.Board.Width
	towardsFood := "up"
	for _, move := range safeMoves {
		target := getCoordFromMove(state, move)
		dist := distanceToNearestFood(state, target)
		if dist < minDistToFood {
			towardsFood = move
			minDistToFood = dist
		}
	}

	return towardsFood
}

func scared_pickMove(state GameState, safeMoves []string) string {
	// run away, run away
	return safeMoves[rand.Intn(len(safeMoves))]
}

//==========================================================
//
// Utilities
//
//==========================================================
func getCoordFromMove(state GameState, move string) Coord {
	myHead := state.You.Body[0]
	targetCoord := myHead
	switch move {
	case "up":
		targetCoord.Y += 1
		break
	case "down":
		targetCoord.Y -= 1
		break
	case "left":
		targetCoord.X -= 1
		break
	case "right":
		targetCoord.X += 1
		break
	default:
		log.Printf("%s ERROR: unknown move: %s\n", state.Game.ID, move)
		break
	}
	return targetCoord
}

func bodyPartsOnDiagonals(state GameState) int {
	relDiags := [4]Coord{
		{X: 1, Y: 1},
		{X: -1, Y: 1},
		{X: -1, Y: -1},
		{X: 1, Y: -1},
	}
	res := 0
	for _, diag := range relDiags {
		absDiag := add(diag, state.You.Head)
		if bodyPartOn(state, absDiag) {
			res += 1
		}
	}

	return res
}

func distanceToNearestFood(state GameState, target Coord) int {
	minDistToFood := state.Board.Height + state.Board.Width
	for _, food := range state.Board.Food {
		dist := distance(target, food)
		if dist < minDistToFood {
			minDistToFood = dist
		}
	}

	return minDistToFood
}

//
// uses the Manhattan distance vs the standard
// Euclidean distance since the snake can't travel
// diagonals
//
func distance(s Coord, t Coord) int {
	return int(abs(int64(s.X-t.X)) + abs(int64(s.Y-t.Y)))
}

//
// Shift-XOR version of Abs
// http://cavaliercoder.com/blog/optimized-abs-for-int64-in-go.html
//
func abs(n int64) int64 {
	y := n >> 63       // y ← x ⟫ 63
	return (n ^ y) - y // (x ⨁ y) - y
}

func bodyPartOn(state GameState, target Coord) bool {
	for _, bodyPart := range state.You.Body {
		if bodyPart == target {
			return true
		}
	}
	return false
}
func add(l Coord, r Coord) Coord {
	return Coord{
		X: l.X + r.X,
		Y: l.Y + r.Y,
	}
}
