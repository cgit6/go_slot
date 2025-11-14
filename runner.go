package main

import (
	"fmt"
	"math/rand"
)

func runner() error {

	// 1. 創建 Config 實例
	cfg, err := NewConfig(REELSTRIPS, SYMBOLS, LINES, PAYTABLE, ROWS, COLS, ModeLines)
	// 錯誤檢查
	if err != nil {
		return err
	}

	// 2. 建立亂數生成
	randSeed := rand.NewSource(123456789) // 固定 randSeed
	// randSeed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(randSeed) // 返回 pointer

	// 3. 建立 生成盤面、算分實例
	sg := NewScreenGenerator(cfg, rng)

	// if err != nil {
	// 	return err
	// }

	// 測試
	screen := sg.GenScreen()
	fmt.Println("screen: ", screen)

	// 4. 初始化模擬參數

	// 5. 執行模擬

	// 6. 計算統計值

	// 7. 打印結果

	return nil

}
