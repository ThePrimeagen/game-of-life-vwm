package main

import (
	"log"
	"time"

	"github.com/theprimeagen/vim-with-me/pkg/testies"
	"github.com/theprimeagen/vim-with-me/pkg/window"
)

var dir = [][]int{
    {-1, -1},
    {0, -1},
    {1, -1},
    {1, 0},
    {1, 1},
    {0, 1},
    {-1, 1},
    {-1, 0},
}

// o(c)
func getNeighborCount(state [][]int, x, y int) int {
    count := 0
    for _, d := range dir {
        row := y + d[0]
        col := x + d[1]

        if row < 0 || row >= len(state) {
            continue
        }

        if col < 0 || col >= len(state[0]) {
            continue
        }

        if state[row][col] == 1 {
            count++
        }
    }
    return count
}

func playRound(state [][]int) [][]int {
    out := make([][]int, len(state))
    for i := range state {
        out[i] = make([]int, len(state[i]))
    }

    /*
    Any live cell with fewer than two live neighbors dies, as if by underpopulation.
    Any live cell with two or three live neighbors lives on to the next generation.
    Any live cell with more than three live neighbors dies, as if by overpopulation.
    Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
    */
    for y, rows := range state {
        for x, _ := range rows {
            out[y][x] = state[y][x]

            count := getNeighborCount(state, x, y)
            if count < 2 {
                out[y][x] = 0
            }
            if count >= 4 {
                out[y][x] = 0
            }
            if count == 3 {
                out[y][x] = 1
            }
        }
    }

    return out
}

func toString(state [][]int) string {
    out := ""
    for _, row := range state {
        for _, cell := range row {
            if cell == 1 {
                out += "x"
            } else {
                out += " "
            }
        }
    }

    return out
}

func main() {
    server, err := testies.CreateServerFromArgs()
    if err != nil {
        log.Fatal("owned by skill issues", err)
    }

    win := window.NewWindow(24, 80)
    server.ToSockets.Welcome(window.OpenCommand(win))
    ticker := time.NewTicker(25 * time.Millisecond)

    state := make([][]int, 24)
    for i := range state {
        state[i] = make([]int, 80)
    }

    // - - x x - -
    // - - x x - -
    // x x x - - -
    // - x - - - -
    state[12][40] = 1
    state[12][41] = 1
    state[13][40] = 1
    state[13][41] = 1
    state[14][38] = 1
    state[14][39] = 1
    state[14][40] = 1
    state[15][39] = 1

    for {
        <-ticker.C

        if server.ToSockets.Len() == 0 {
            continue
        }

        state = playRound(state)
        err := win.SetWindow(toString(state))
        if err != nil {
            log.Fatal("owned by skill issues part deux", err)
        }
        cmds := win.Flush()
        for _, cmd := range cmds {
            server.ToSockets.Spread(cmd.Command())
        }
    }
}

