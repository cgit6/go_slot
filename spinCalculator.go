package main

import (
	"log"
)

type LineResult struct {
	sym    uint8 // 符號 ID
	cnt    int   // 連線數量
	win    int   // 賠分
	lineId int   // 線路 ID
}

// 一次 spin 的結果
type ScreenResult struct {
	C1Win      int           // 盤面中 C1 (scatter) 出現次數
	Win        int           // 累積賠分)
	LineResult *[]LineResult // 線路結果
}

// 輸入 盤面與 1 次 spin 下注分數
type CalcFunc func(*SpinCalculator, []uint8, int) *ScreenResult // 接收 *SpinCalculator

type SpinCalculator struct {
	*Config                // 匿名嵌入
	*ScreenResult          // 結果緩存
	minLen        int      // 最小連線
	calcFn        CalcFunc // 算分函數

	// 輔助參數
	filterIds []uint8 // 特殊符號索引
}

// 建構函數: 創建 NewSpinCalculator instance 時調用
func NewSpinCalculator(cfg *Config) *SpinCalculator {
	sc := &SpinCalculator{
		Config:       cfg,
		ScreenResult: &ScreenResult{},
		minLen:       3, // 若你有 cfg.MinLen 改用它
	}
	sc.initCalcFn()
	sc.filterIds = deriveFilterIDs(cfg.Paytable, cfg.W1Id) // ← 自動算不計分符號
	return sc
}

func deriveFilterIDs(pay [][]int, wildID uint8) []uint8 {
	out := make([]uint8, 0, len(pay))
	for sid, row := range pay {
		allZero := true
		for _, p := range row {
			if p != 0 {
				allZero = false
				break
			}
		}
		if allZero && uint8(sid) != wildID {
			out = append(out, uint8(sid))
		}
	}
	return out
}

func countSymbol(screen []uint8, id uint8) int {
	// 給你保留一個「無效」的用法；可刪
	if id == 0xFF {
		return 0
	}
	n := 0
	for _, v := range screen {
		if v == id {
			n++
		}
	}
	return n
}

func containsU8(arr []uint8, x uint8) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

// 維護一個map註冊表
var calcFnMap = map[GameMode]CalcFunc{
	ModeLines: CalcLinesGame, // lines 算法
	ModeWays:  CalcWaysGame,  // ways 算法

}

// 選擇算分方式
func (s *SpinCalculator) initCalcFn() {

	// 選擇算分策略
	if fn, ok := calcFnMap[s.Mode]; ok {
		s.calcFn = fn // 選擇算分方式存到 s.calcFn

		return // 必要，不然會往外跳執行 log.Fatal("未知 mode")
	}
	log.Fatal("未知 mode")
	// panic 表示還有救，但這個沒救了(設定檔錯誤)，類似 try ... catch ...
	// 這邊可以用 errors.New() 嗎?

}

// ------- 不同算分模式的內部函數 -------

// lines 算分模式
func CalcLinesGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	r := s.ScreenResult

	// 重置此次結果
	r.C1Win = 0
	r.Win = 0

	// 統計 C1（scatter）次數
	r.C1Win = countSymbol(screen, s.C1Id)

	rows, cols := s.Rows, s.Cols
	wildID := s.W1Id // uint8
	minLen := s.minLen
	linesLen := len(s.Lines)

	totalLinePay := 0
	lineResults := make([]LineResult, 0, linesLen)

	// 逐條線計分
	for i := 0; i < linesLen; i++ {
		// 單條線的狀態
		wildCount := 0
		wildContinue := true

		var symId uint8     // 用 uint8 保存符號ID
		symStarted := false // 取代原本用 -1 當哨兵
		symCount := 0
		pendingWilds := 0

		// 從左到右掃這條線
		for j := 0; j < cols; j++ {
			rowIndex := s.Lines[i][j]
			idx := j*rows + rowIndex
			sid := screen[idx] // uint8

			// 開頭連續 Wild 數
			if wildContinue && sid == wildID {
				wildCount++
			} else {
				wildContinue = false
			}

			// 尚未決定得分符號
			if !symStarted {
				if sid == wildID {
					// 前置 Wild 先累計，遇到第一個可計分符號時再併入
					pendingWilds++
					continue
				}
				// 第一個非 Wild：若是不計分符號（Z1/C1 等），此線只能靠純 Wild
				if containsU8(s.filterIds, sid) {
					break
				}
				// 合法得分符號確立
				symId = sid
				symStarted = true
				symCount = pendingWilds + 1
				continue
			}

			// 已決定得分符號，延伸連線：同符號或 Wild 都可
			if sid == symId || sid == wildID {
				symCount++
			} else {
				break
			}
		}

		// 未達最小連線長度 → 0 分（仍可記錄 line 結果，win=0）
		if symCount < minLen && wildCount < minLen {
			lineResults = append(lineResults, LineResult{
				sym:    0, // 無得分；若想避免混淆可自訂常數 255 表示「無」
				cnt:    0,
				win:    0,
				lineId: i,
			})
			continue
		}

		// 計算兩種賠率（只在這裡把 uint8 轉成 int 當索引）
		symPay := 0
		if symStarted { // 已有合法得分符號
			row := s.Paytable[int(symId)]
			k := symCount - 1
			if k >= 0 && k < len(row) {
				symPay = row[k]
			}
		}
		wildRow := s.Paytable[int(wildID)]
		wildPay := 0
		if wc := wildCount - 1; wc >= 0 && wc < len(wildRow) {
			wildPay = wildRow[wc]
		}

		// 取較大者
		winSym := symId
		winCnt := symCount
		winPay := symPay
		if wildPay > symPay {
			winSym = wildID
			winCnt = wildCount
			winPay = wildPay
		}

		totalLinePay += winPay
		lineResults = append(lineResults, LineResult{
			sym:    winSym, // uint8
			cnt:    winCnt,
			win:    winPay,
			lineId: i,
		})
	}

	// 本把贏分（以 lineBet = bet 為基數）
	r.Win = totalLinePay * bet
	r.LineResult = &lineResults
	return r
}

// ways 算分模式
func CalcWaysGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	// 未實做
	return s.ScreenResult
}
