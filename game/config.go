package game

import (
	"errors"
	"log"
)

// type SymbolID uint8 // 之後可以統一使用
type GameMode int

const (
	ModeUnknown GameMode = iota // 0 -> unknown
	ModeLines                   // 1 -> Line
	ModeWays                    // 2 -> Ways
)

// 輪帶表
var REELSTRIPS = [][]uint8{
	{3, 4, 5, 6, 7, 8, 9, 10, 1, 2}, // 第 1 軸
	{3, 4, 5, 6, 7, 8, 9, 10, 1, 2}, // 第 2 軸
	{3, 4, 5, 6, 7, 8, 9, 10, 1, 2}, // 第 3 軸
	{3, 4, 5, 6, 7, 8, 9, 10, 1, 2}, // 第 4 軸
	{3, 4, 5, 6, 7, 8, 9, 10, 1, 2}, // 第 5 軸
}

// 11 個有效符號           0      1     2     3     4     5     6     7     8     9     10
var SYMBOLS = []string{"None", "C1", "W1", "H1", "H2", "H3", "H4", "L1", "L2", "L3", "L4"}

// 20 線路表
var LINES = [][]int{
	{1, 1, 1, 1, 1}, // 線路 1
	{0, 0, 0, 0, 0}, // 線路 2
	{2, 2, 2, 2, 2}, // ...
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
	{0, 2, 2, 2, 0}, // 線路 20
}

// 賠率表
var PAYTABLE = [][]int{
	{0, 0, 0, 0, 0},       // Z1
	{0, 0, 10, 12, 14},    // C1 (Scatter)
	{0, 0, 10, 150, 6667}, // W1 (Wild)
	{0, 0, 20, 50, 200},   // H1
	{0, 0, 8, 50, 200},    // H2
	{0, 0, 8, 13, 100},    // H3
	{0, 0, 8, 13, 100},    // H4
	{0, 0, 8, 12, 50},     // L1
	{0, 0, 8, 12, 50},     // L2
	{0, 0, 8, 12, 50},     // L3
	{0, 0, 8, 12, 50},     // L4
}

var ROWS, COLS int = 3, 5 // 列數, 行數

type Config struct {
	// 設定檔的數值
	ReelStrips [][]uint8 // 輪帶表
	Symbols    []string  // 符號清單
	Lines      [][]int   // 線獎組合
	Paytable   [][]int   // 賠率表
	Rows       int       // 列數
	Cols       int       // 軸數
	Mode       GameMode  // 算分模式(enum)

	// 輔助的數值
	ScreenSize int   // 盤面大小
	ReelLens   []int // 每一軸輪帶長度
	C1Id       uint8 // scatter 索引值
	W1Id       uint8 // wild 索引值
	minLen     int   // 最小連線長度
	FlatLines  []int // 平坦化線路清單

	// 初始化狀態
	initFlag bool // 初始化旗標

}

// 建構函數: 創建 instance 時調用
func NewConfig(reelStrips [][]uint8, symbols []string, lines [][]int, payTable [][]int, rows int, cols int, mode GameMode) (*Config, error) {
	// 1. 創建 Config instance & 賦值
	cfg := &Config{
		ReelStrips: reelStrips, // 輪帶表
		Symbols:    symbols,    // 符號清單
		Lines:      lines,      // 線路清單
		Paytable:   payTable,   // 賠率表
		Rows:       rows,       // 列數
		Cols:       cols,       // 行數
		minLen:     3,          // 最小連線長度(改成用公式判斷)
		Mode:       mode,       // 算分模式
	}

	// 2. 執行初始化
	if err := cfg.Init(); err != nil {
		return nil, err
	}
	// 3. 返回值, 錯誤訊息
	return cfg, nil

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

	// 4. 找到特殊符號索引
	var i uint8 = 0
	for i < uint8(len(c.Symbols)) {
		if c.Symbols[i] == "C1" {
			c.C1Id = i
		}

		if c.Symbols[i] == "W1" {
			c.W1Id = i
		}

		i++
	}

	// 5. 平坦化線路清單
	c.FlatLines = make([]int, len(c.Lines)*c.Cols) // 預先分配空間
	for i, line := range c.Lines {
		for j, pos := range line {
			c.FlatLines[i*c.Cols+j] = j*c.Rows + pos // j*c.Rows + pos -> 索引值
		}
	}

	// 6. 更新初始化狀態
	c.initFlag = true
	return nil
}

func (c *Config) Reset() error {

	// 檢查 initFlag 狀態
	if !c.initFlag {
		return errors.New("not yet init")
	}

	// 重新初始化
	c.initFlag = false
	c.Init()

	return nil
}

func (c *Config) validate() error {

	// 1. Rows/Cols
	if c.Rows <= 0 {
		log.Fatal("rows must > 0")
	}

	if c.Cols <= 0 {
		log.Fatal("cols must > 0")
	}
	// 2. 輪帶
	if len(c.ReelStrips) != c.Cols {
		log.Fatal("reelStrips length  must equal Cols")
	}

	// 3. 符號清單，這邊怪怪的感覺有很多例外
	symLen := len(c.Symbols)
	if symLen == 0 {
		log.Fatal("symbols must not be empty")
	}

	// 4. Line Mode
	if c.Mode == ModeLines {
		if len(c.Lines) == 0 {
			log.Fatal("line must not be emypt")
		}
	}

	if c.Mode == ModeWays {
		log.Fatal("未實作")
	}

	// 5. PayTable： 每個符號 5 欄（1~5 連）
	if len(c.Paytable) != symLen {
		log.Fatal("paytable size not correct")
	}

	// 6. 模式檢查
	// 這邊應該改成不存在於 GameMode enum 清單中，或是 =0
	if c.Mode != ModeLines && c.Mode != ModeWays {
		log.Fatal("invalid mode")
	}

	return nil
}
