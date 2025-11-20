package main

import (
	"fmt"
	"math/rand"
	"time"
)

func runner() error {

	// 1. 創建 Config 實例
	cfg, err := NewConfig(REELSTRIPS, SYMBOLS, LINES, PAYTABLE, ROWS, COLS, ModeLines)

	// 錯誤檢查
	if err != nil {
		return err
	}

	// 2. 建立亂數生成
	// randSeed := rand.NewSource(123456789) // 固定 randSeed
	randSeed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(randSeed) // 返回 pointer

	// 3. 建立 生成盤面、算分實例
	sg := NewScreenGenerator(cfg, rng)
	sc := NewSpinCalculator(cfg)

	// 4. 初始化模擬參數
	rounds := 1000_000_00 // 模擬次數
	bet := 10             // Bet: 一次 spin 下注分數
	totalBet := 0
	totalWin := 0
	start := time.Now() // 起始時間

	// 5. 執行模擬
	for i := 0; i < rounds; i++ {
		// 執行模擬
		screen := sg.GenScreen()
		result := sc.calcFn(sc, screen, bet)

		// 更新狀態
		totalBet += bet        // 總下注
		totalWin += result.Win // 總贏分

		//  顯示進度

	}

	if totalBet == 0 {
		return nil
	}

	elapsed := time.Since(start)

	fmt.Printf("Elapsed time: %.6f seconds\n", elapsed.Seconds())

	// 6. 計算統計值
	rtp := float64(totalWin) / float64(totalBet)
	fmt.Printf("TotalBet=%d TotalWin=%d RTP=%.6f\n", totalBet, totalWin, rtp)
	return nil

}
