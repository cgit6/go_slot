package game

import (
	"fmt"
	"math/rand"
	"time"
)

func Runner() error {

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
	rounds := 1000_000_0 // 模擬次數
	bet := len(LINES)    // Bet: 一次 spin 下注分數，保持與 paytable 長度一致避免小數點產生
	totalBet := 0        // 總下注
	totalWin := 0        // 累積贏分
	start := time.Now()  // 起始時間

	// ------------------------------------------------------
	totalC1Win := 0  // 總 C1 贏分
	totalSymWin := 0 // 總符號贏分
	// ------------------------------------------------------

	// 5. 執行模擬
	for i := 0; i < rounds; i++ {
		// 執行模擬
		screen := sg.GenScreen()
		baseGameResult := sc.calcFn(sc, screen, bet) // base game result

		// 如果 C1 數量大於 3 顆，執行 15 局免費遊戲
		if baseGameResult.C1Win > 3 {
			// 執行免費遊戲
			for fgRound := 0; fgRound < 15; fgRound++ {
				// 生成免費遊戲盤面
			}
		}

		// 更新狀態
		totalBet += bet // 總下注

		// ------------------------------------------------------
		totalC1Win += baseGameResult.C1Win   // 總 C1 贏分
		totalSymWin += baseGameResult.SymWin // 總符號贏分
		// ------------------------------------------------------

		totalWin += baseGameResult.Win // 總贏分

		//  顯示進度

	}

	if totalBet == 0 {
		return nil
	}

	elapsed := time.Since(start)

	fmt.Printf("Elapsed time: %.6f seconds\n", elapsed.Seconds())

	// 6. 計算統計值

	// ------------------------------------------------------
	fmt.Printf("Total C1 Win=%f\n", float64(totalC1Win)/float64(totalBet)*100)   // C1
	fmt.Printf("Total Sym Win=%f\n", float64(totalSymWin)/float64(totalBet)*100) // 得分符號
	// ------------------------------------------------------

	// rtp
	rtp := float64(totalWin) / float64(totalBet)
	fmt.Printf("TotalBet=%d TotalWin=%d RTP=%.6f\n", totalBet, totalWin, rtp*100)
	return nil

	// cv

	// fg trigger rate

	// hit rate

}
