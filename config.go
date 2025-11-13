package main

import "errors"

type GameMode int

const (
	ModeLine GameMode = iota // 0 -> Line
	ModeWays                 // 1 -> Ways
)

// 輪帶表
var REELSTRIPS = [][]int{
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // 第 1 軸
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // 第 2 軸
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // 第 3 軸
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // 第 4 軸
	{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 2, 3, 4, 5, 6, 7, 8, 9, 10}, // 第 5 軸
}

// 11 個有效符號
var SYMBOLS = []string{"None", "Scatter", "Wild", "H1", "H2", "H3", "H4", "L1", "L2", "L3", "L4", "L5"}

// 20 線路表
var LINES = [][]int{
	{1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0},
	{2, 2, 2, 2, 2},
	{0, 1, 2, 1, 0},
	{2, 1, 0, 1, 2},
	{1, 0, 0, 0, 1},
	{1, 2, 2, 2, 1},
	{0, 0, 1, 2, 2},
	{2, 2, 1, 0, 0},
	{1, 0, 1, 2, 1},
	{1, 2, 1, 0, 1},
	{0, 1, 1, 1, 0},
	{2, 1, 1, 1, 2},
	{0, 1, 0, 1, 0},
	{2, 1, 2, 1, 2},
	{1, 1, 0, 1, 1},
	{1, 1, 2, 1, 1},
	{0, 0, 2, 0, 0},
	{2, 2, 0, 2, 2},
	{0, 2, 2, 2, 0},
}

// 賠率表
var PAYTABLE = [][]int{
	{0, 0, 0, 0, 0},       // Z1
	{0, 0, 0, 0, 0},       // C1 (Scatter)
	{0, 0, 100, 200, 300}, // W1 (Wild)
	{0, 0, 10, 50, 200},   // H1
	{0, 0, 10, 50, 200},   // H2
	{0, 0, 10, 50, 200},   // H3
	{0, 0, 10, 50, 200},   // H4
	{0, 0, 5, 20, 100},    // L1
	{0, 0, 5, 20, 100},    // L2
	{0, 0, 5, 20, 100},    // L3
	{0, 0, 5, 20, 100},    // L4
	{0, 0, 5, 20, 100},    // L5
}

var ROWS, COLS int = 3, 5 // 列數, 行數

type Config struct {
	// 設定檔有的數值
	ReelStrips [][]int  // 輪帶表
	Symbols    []string // 符號清單
	Lines      [][]int  // 線獎組合
	Paytable   [][]int  // 賠率表
	Rows       int      // 列數
	Cols       int      // 軸數
	Mode       GameMode // 算分模式(enum)

	// 設定檔案沒有的靜態數值
	ScreenSize int   // 盤面大小
	ReelLens   []int // 每一軸輪帶長度

	// 初始化狀態
	initFlag bool // 初始化旗標

}

// 初始化方法
func (c *Config) Init() error {
	// 1. 先檢查 initFlag
	if c.initFlag {
		return nil
	}

	// 2. 執行設定檔驗證
	if err := c.validate(); err != nil {
		return err
	}

	// 3. 計算盤面大小、輪帶長度清單
	c.ScreenSize = c.Rows * c.Cols   // 盤面大小
	c.ReelLens = make([]int, c.Cols) // 輪帶長度清單
	for col := 0; col < c.Cols; col++ {
		c.ReelLens[col] = len(c.ReelStrips[col])
	}

	// 4. 更新初始化狀態
	c.initFlag = true
	return nil
}

func (c *Config) validate() error {

	// 1. Rows/Cols
	if c.Rows <= 0 {
		return errors.New("rows must > 0")
	}

	if c.Cols <= 0 {
		return errors.New("cols must > 0")
	}
	// 2. 輪帶
	if len(c.ReelStrips) != c.Cols {
		return errors.New("reelStrips length  must equal Cols")
	}

	// 3. 符號清單，這邊怪怪的感覺有很多例外
	if len(c.Symbols) == 0 {
		return errors.New("symbols must not be emypt")
	}

	// 4. Line Mode
	if c.Mode == ModeLine {
		if len(c.Lines) == 0 {
			return errors.New("line must not be emypt")
		}
		// 檢查 線獎
		// if len(c.Lines) > 0 {
		// 	for i, line :=
		// }
	}

	// 5. PayTable： 每個符號 5 欄（1~5 連）

	// 6. 模式檢查

	return nil
}

// func

// 建構函數: 創建 instance 時調用
func NewConfig(reelStrips [][]int, symbols []string, lines [][]int, payTable [][]int, rows, cols int, mode GameMode) (*Config, error) {
	// 1. 賦值
	cfg := &Config{
		ReelStrips: reelStrips,
		Symbols:    symbols,
		Lines:      lines,
		Paytable:   payTable,
		Rows:       rows,
		Cols:       cols,
		Mode:       mode,
	}

	// 2. 執行初始化
	if err := cfg.Init(); err != nil {
		return nil, err
	}

	return cfg, nil

}
