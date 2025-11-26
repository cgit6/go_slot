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
	rounds := 1000_000_0 // 模擬次數
	bet := len(LINES)    // Bet: 一次 spin 下注分數，保持與 paytable 長度一致避免小數點產生
	// 統計值需要的變數值
	totalBet := 0 // 總下注
	totalWin := 0 // 累積贏分
	// ------------------------------------------------------
	BGC1Win := 0  // 總 C1 贏分
	BGSymWin := 0 // 總符號贏分
	FGC1Win := 0  // FG C1 贏分
	FGSymWin := 0 // FG 符號贏分
	// ------------------------------------------------------
	roundWin := 0 // 1 spin wins

	// Total Game 用
	meanTotal := 0.0 // 總遊戲的樣本平均
	m2Total := 0.0   // 總遊戲的偏差平方和
	nTotal := 0      // 總遊戲樣本數

	// Base Game 專用 (以付費 spin 為單位)
	meanBG := 0.0
	m2BG := 0.0
	nBG := 0

	// Free Game 專用 (以「每次 FG session」為單位)
	meanFG := 0.0
	m2FG := 0.0
	nFG := 0

	hit := 0        // 得分次數
	trigger := 0    // fg 觸發率
	spinHits := 0   // 中獎 spin 次數
	totalSpins := 0 // 總 spins 次數

	// 5. 執行模擬
	start := time.Now() // 起始時間
	round := 1          // 從第一次開始
	for ; round <= rounds; round++ {
		roundWin = 0 // reset 1 spin

		// 執行 base game 模擬
		screen := bgsg.GenScreen()
		baseGameResult := bgsc.calcFn(bgsc, screen, bet, BaseGame) // 1 spin base game result

		roundWin += baseGameResult.Win

		// ---- BG CV 樣本：每局付費 spin 的 BG 贏分倍數 ----
		xBG := float64(baseGameResult.Win) / float64(bet)
		nBG++
		deltaBG := xBG - meanBG
		meanBG += deltaBG / float64(nBG)
		deltaBG2 := xBG - meanBG
		m2BG += deltaBG * deltaBG2

		// 更新 spin 值
		totalSpins++
		if baseGameResult.Win > 0 {
			spinHits++
		}

		// 執行 free game 模擬
		if fgRound, ok := C1ToFG[baseGameResult.C1Count]; ok { // 如果 BG C1 數量 3,4,5 顆，執行 10,10,10 局免費遊戲，我利用 map 限定事件發生範圍，避免意外情況發生。
			fg := 0 // 初始化 FG 局數計數，因為有加局機制，使用 while loop

			trigger++
			sessionWin := 0 // 這一包 FG（含 retrigger）的總贏分

			for fg < fgRound {
				fgScreen := fgsg.GenScreen()                                 // 生成 free game 盤面
				freeGameResult := fgsc.calcFn(fgsc, fgScreen, bet, FreeGame) // 1 spin free game result

				// 判斷是否有 加局
				if extraSpins, ok := ReTriggerFGSpins[freeGameResult.C1Count]; ok {
					fgRound += extraSpins // 加局
				}

				// 計算 CV 需要的參數
				roundWin += freeGameResult.Win
				sessionWin += freeGameResult.Win //

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
			// ---- FG CV 樣本：每次 FG session 的總贏分倍數 ----
			xFG := float64(sessionWin) / float64(bet)
			nFG++
			deltaFG := xFG - meanFG
			meanFG += deltaFG / float64(nFG)
			deltaFG2 := xFG - meanFG
			m2FG += deltaFG * deltaFG2
		}

		// 更新狀態
		// ------------------------------------------------------
		totalBet += bet                   // 總下注
		BGC1Win += baseGameResult.C1Win   // 總 C1 贏分
		BGSymWin += baseGameResult.SymWin // 總符號贏分
		totalWin += baseGameResult.Win    // 總贏分

		// ---- Total Game CV 樣本：每局付費 spin 的「BG+FG 總倍數」 ----
		xTotal := float64(roundWin) / float64(bet)
		nTotal++
		deltaTotal := xTotal - meanTotal
		meanTotal += deltaTotal / float64(nTotal)
		deltaTotal2 := xTotal - meanTotal
		m2Total += deltaTotal * deltaTotal2

		// 統計得分次數(1 次 bg spin 是否有得分)
		if roundWin > 0 {
			hit++
		}

	}

	if totalBet == 0 {
		return nil
	}

	elapsed := time.Since(start)

	// 6. 結果
	if round > 1 {

		// 6.1 計算統計值
		BGC1RTP := float64(BGC1Win) / float64(totalBet)   // base game C1 rtp
		BGSymRTP := float64(BGSymWin) / float64(totalBet) // base game 得分符號 rtp
		// FGC1RTP := float64(FGC1Win) / float64(totalBet)   // free game C1 rtp
		// FGSymRTP := float64(FGSymWin) / float64(totalBet) // free game 得分符號 rtp
		// ---------------------------
		// Total Game
		rtp := float64(totalWin) / float64(totalBet) // RTP (倍數)
		varTotal := m2Total / float64(nTotal-1)
		stdTotal := math.Sqrt(varTotal)
		cvTotal := stdTotal / rtp // 這裡用 rtp 跟 meanTotal 差異只會在數值精度

		// // Base Game （每個付費 spin 的 BG 贏分）
		// varBG := m2BG / float64(nBG-1)
		// stdBG := math.Sqrt(varBG)
		// // BG 的平均倍數就是 meanBG （也可以用 BGC1RTP + BGSymRTP 驗證）
		// cvBG := stdBG / meanBG

		// // Free Game （每次 FG session 的總贏分）
		// varFG := m2FG / float64(nFG-1)
		// stdFG := math.Sqrt(varFG)
		// cvFG := stdFG / meanFG
		FGTriggerRate := float64(trigger) / float64(rounds) // Free Game 觸發率
		hitProd := float64(hit) / float64(rounds)           // game 得分率(一局 pay spin 的 win > 0 次數 / 所有 pay spin 次數)

		// 6.2. 輸出結果
		fmt.Printf("Elapsed time: %.4f seconds\n", elapsed.Seconds())
		fmt.Println("----------------------- detail ----------------------------")
		fmt.Printf("BG RTP=%4f\n", BGC1RTP+BGSymRTP) // base game rtp
		// fmt.Printf("BG Sym Win=%f\n", BGSymRTP)    // base game 得分符號 rtp
		// fmt.Printf("FG C1 Win=%f\n", FGC1RTP)      // Free Game C1 rtp
		// fmt.Printf("FG Sym Win=%f\n", FGSymRTP)    // Free Game 得分符號 rtp
		// fmt.Printf("cv_bg_per_spin=%f\n", cvBG)    // bg cv
		// fmt.Printf("cv_fg_per_session=%f\n", cvFG) // fg cv

		fmt.Println("----------------------- all ----------------------------")
		fmt.Printf("TotalBet=%d\n", totalBet)       // 總下注
		fmt.Printf("TotalWin=%d\n", totalWin)       // 總贏
		fmt.Printf("RTP=%.4f\n", rtp)               // rtp
		fmt.Printf("cv=%4f\n", cvTotal)             // cv
		fmt.Printf("trigger=%4f \n", FGTriggerRate) // Free Game 觸發率
		fmt.Printf("hit rate=%4f \n", hitProd)      // 得分率
		// fmt.Printf("hit rate(any spin)=%f\n", float64(spinHits)/float64(totalSpins)) // 每一局獨立看
	}

	return nil

}
