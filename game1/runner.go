package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Runner() error {

	// 1. 創建 Config 實例
	bgCfg, err := NewConfig(BGREELSTRIPS, SYMBOLS, LINES, PAYTABLE, ROWS, COLS, ModeLines) // base game 設定

	if err != nil {
		return err
	}

	fgCfg, err := NewConfig(FGREELSTRIPS, SYMBOLS, LINES, PAYTABLE, ROWS, COLS, ModeLines) // free game 設定

	if err != nil {
		return err
	}

	// 2. 建立亂數生成
	// randSeed := rand.NewSource(123456789) // 固定 randSeed
	randSeed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(randSeed) // 返回 pointer

	// 3. 建立 生成盤面、算分實例

	// bg ------------------------------------------------------
	bgsg := NewScreenGenerator(bgCfg, rng)
	bgsc := NewSpinCalculator(bgCfg)
	// fg ------------------------------------------------------
	fgsg := NewScreenGenerator(fgCfg, rng)
	fgsc := NewSpinCalculator(fgCfg)

	// 4. 初始化模擬參數
	rounds := 1000_000_000 // 模擬次數
	bet := len(LINES)      // Bet: 一次 spin 下注分數，保持與 paytable 長度一致避免小數點產生
	totalBet := 0          // 總下注
	totalWin := 0          // 累積贏分
	start := time.Now()    // 起始時間

	// 基礎驗證
	// ------------------------------------------------------
	BGC1Win := 0  // 總 C1 贏分
	BGSymWin := 0 // 總符號贏分
	// ------------------------------------------------------
	FGC1Win := 0  // FG C1 贏分
	FGSymWin := 0 // FG 符號贏分
	// ------------------------------------------------------

	// ------------------------------------------------------
	roundWin := 0 // 1 spin win
	mean := 0.0   // 樣本平均值
	m2 := 0.0     // 偏差平方和 (Welford)

	hit := 0     // 得分次數
	trigger := 0 // fg 觸發率

	totalSpins := 0 // 總 spins 次數
	spinHits := 0   // 中獎 spin 次數

	// 5. 執行模擬
	round := 1 // 從第一次開始
	for ; round <= rounds; round++ {
		roundWin = 0 // reset 1 spin

		// 執行 base game 模擬
		screen := bgsg.GenScreen()
		baseGameResult := bgsc.calcFn(bgsc, screen, bet, BaseGame) // 1 spin base game result

		roundWin += baseGameResult.Win

		// 更新 spin 值
		totalSpins++
		if baseGameResult.Win > 0 {
			spinHits++
		}

		// 執行 free game 模擬
		if fgRound, ok := C1ToFG[baseGameResult.C1Count]; ok { // 如果 BG C1 數量 3,4,5 顆，執行 10,10,10 局免費遊戲，我利用 map 限定事件發生範圍，避免意外情況發生。
			fg := 0 // 初始化 FG 局數計數，因為有加局機制，使用 while loop

			trigger++

			for fg < fgRound {
				fgScreen := fgsg.GenScreen()                                 // 生成 free game 盤面
				freeGameResult := fgsc.calcFn(fgsc, fgScreen, bet, FreeGame) // 1 spin free game result

				// 判斷是否有 加局
				if extraSpins, ok := ReTriggerFGSpins[freeGameResult.C1Count]; ok {
					fgRound += extraSpins // 加局
				}

				// 計算 CV 需要的參數
				roundWin += freeGameResult.Win

				// 更新狀態
				// ------------------------------------------------------
				FGC1Win += freeGameResult.C1Win   // 總 C1 贏分
				FGSymWin += freeGameResult.SymWin // 總符號贏分
				// ------------------------------------------------------

				totalWin += freeGameResult.Win // 總贏分

				// 更新 spins
				totalSpins++
				if freeGameResult.Win > 0 {
					spinHits++
				}

				// 更新當前局數
				fg++
			}
		}

		// 更新狀態
		totalBet += bet // 總下注

		// ------------------------------------------------------
		BGC1Win += baseGameResult.C1Win   // 總 C1 贏分
		BGSymWin += baseGameResult.SymWin // 總符號贏分
		// ------------------------------------------------------

		totalWin += baseGameResult.Win // 總贏分

		// 計算樣本標準差 Welford 算法
		x := float64(roundWin) / float64(bet) // 1 pay spin get pay value

		delta := float64(x) - mean
		mean += delta / float64(round)
		delta2 := float64(x) - mean
		m2 += delta * delta2

		// 統計得分次數(1 次 bg spin 是否有得分)
		if roundWin > 0 {
			hit++
		}

	}

	if totalBet == 0 {
		return nil
	}

	elapsed := time.Since(start)

	fmt.Printf("Elapsed time: %.6f seconds\n", elapsed.Seconds())

	// 6. 計算統計值

	// 6.1. 基礎驗證
	// ------------------------------------------------------
	fmt.Printf("BG C1 Win=%f\n", float64(BGC1Win)/float64(totalBet))   // base game C1
	fmt.Printf("BG Sym Win=%f\n", float64(BGSymWin)/float64(totalBet)) // base game 得分符號
	// ------------------------------------------------------
	fmt.Println()
	fmt.Printf("FG C1 Win=%f\n", float64(FGC1Win)/float64(totalBet)) // Free Game C1
	fmt.Printf("FG Sym Win=%f\n", float64(FGSymWin)/float64(totalBet))

	// 6.2. 輸出結果
	if round > 1 {
		// rtp
		rtp := float64(totalWin) / float64(totalBet)
		fmt.Printf("TotalBet=%d TotalWin=%d RTP=%.6f\n", totalBet, totalWin, rtp)

		// cv
		variance := m2 / float64(round-1)
		std := math.Sqrt(variance)

		cv := std / rtp
		fmt.Printf("cv=%f\n", cv)

		// fg trigger rate
		fmt.Printf("fg trigger rate=%f\n", float64(trigger)/float64(rounds))

		// hit rate
		fmt.Printf("hit rate(pay spin)=%f\n", float64(hit)/float64(rounds))
		// fmt.Printf("hit rate(any spin)=%f\n", float64(spinHits)/float64(totalSpins))
	}

	return nil

}
