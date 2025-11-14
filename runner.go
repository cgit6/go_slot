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

	fmt.Println("C1 id: ", cfg.C1Id)
	fmt.Println("W1 id: ", cfg.W1Id)

	// 2. 建立亂數生成
	randSeed := rand.NewSource(123456789) // 固定 randSeed
	// randSeed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(randSeed) // 返回 pointer

	// 3. 建立 生成盤面、算分實例
	sg := NewScreenGenerator(cfg, rng)
	r := NewScreenResult()
	sc := NewSpinCalculator(cfg, r)

	// if err != nil {
	// 	return err
	// }

	// 4. 初始化模擬參數
	rounds := 1 // 模擬次數
	bet := 1000 // Bet: 一次 spin 下注分數
	totalBet := 0
	totalWin := 0
	// start := time.Now().

	// 5. 執行模擬
	for i := 0; i < rounds; i++ {
		// 執行模擬
		screen := sg.GenScreen()
		result := CalcScreen(sc, screen, bet)

		// 更新狀態
		totalBet += result.TotalBets // 總下注
		totalWin += result.TotalWins // 總贏分

		// 顯示進度
	}

	if totalBet == 0 {
		return nil
	}

	// end := time.Now() - start

	// 6. 計算統計值
	rtp := float64(totalWin / totalBet)
	fmt.Printf("TotalBet=%T TotalWin=%T RTP=%.6f", totalBet, totalWin, rtp)
	return nil

}
