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
	C1Win      int          // 盤面中 C1 (scatter) 出現次數
	Win        int          // 累積賠分
	LineResult []LineResult // 線路結果
}

// input SpinCalculator、screen 與 1 次 spin 下注分數
type CalcFunc func(*SpinCalculator, []uint8, int) *ScreenResult // 接收 *SpinCalculator

type SpinCalculator struct {
	*Config                // 匿名嵌入
	*ScreenResult          // 結果緩存
	calcFn        CalcFunc // 算分函數

	// 輔助參數
	filterIds []uint8 // 特殊符號
}

// 建構函數: 創建 NewSpinCalculator instance 時調用
func NewSpinCalculator(cfg *Config) *SpinCalculator {
	sc := &SpinCalculator{
		Config:       cfg,
		ScreenResult: &ScreenResult{},
	}
	sc.initCalcFn()
	sc.filterIds = deriveFilterIDs(cfg.Paytable, cfg.W1Id) // ← 自動算不計分符號
	return sc
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

}

// 不計分符號清單
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

// 計算盤面中特定符號出現次數
func countSymbol(screen []uint8, id uint8) int {
	n := 0
	for _, v := range screen {
		if v == id {
			n++
		}
	}
	return n
}

// 判斷 slice 中是否包含某 uint8 元素
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

// ------- 不同算分模式的內部函數 -------

// lines 算分模式
func CalcLinesGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {

	// 初始化結果
	r := s.ScreenResult

	r.C1Win, r.Win = 0, 0
	linesLen := len(s.Lines)                       // 線路數量
	r.LineResult = make([]LineResult, 0, linesLen) // 清空線路結果

	totalLinePay := 0 // 累積線路賠分

	// 計算 C1 出現次數
	r.C1Win = countSymbol(screen, s.C1Id)

	// 逐條線計分
	for i := 0; i < linesLen; i++ {
		// 單條線的狀態
		wildCount := 0
		wildContinue := true

		var symId uint8     // 得分符號ID
		symStarted := false // 是否已確定得分符號
		symCount := 0       // 符號連線數量
		pendingWilds := 0   // 在算得分符號連線時先算前面累積的 W1 數量，如果後面遇到得分符號就加進去

		// 從左到右掃這條線
		for j := 0; j < s.Cols; j++ {

			// 1. 獲取該位置符號
			sid := screen[j*s.Rows+s.Lines[i][j]]

			// 2. 開頭連續 Wild 數
			if wildContinue && sid == s.W1Id {
				wildCount++
			} else {
				wildContinue = false
			}

			// 3. 計算得分符號連線

			// 3.1. 尚未決定得分符號
			if !symStarted {
				if sid == s.W1Id {
					pendingWilds++ // 預先累積 Wild ，後面可能會替代為得分符號
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

			// 3.2. 已決定得分符號，延伸連線：同符號或 Wild 都可
			if sid == symId || sid == s.W1Id {
				symCount++
			} else {
				break // 如果開頭直接是 C1 直接結束
			}
		}

		// 4. 未達最小連線長度 → 0 分 該條線沒中
		if symCount < s.minLen && wildCount < s.minLen {
			r.LineResult = append(r.LineResult, LineResult{
				sym:    0, // 無得分
				cnt:    0,
				win:    0,
				lineId: i,
			}) // 更新結果
			continue
		}

		// 5. 計算兩種賠率

		// 5.1. 得分符號賠率
		symPay := 0
		if symStarted && symCount >= s.minLen { // 只做「是否該算」的必要判斷
			symPay = s.Paytable[int(symId)][symCount-1]
		}

		// 5.2.Wild 賠率

		wildPay := 0 // W1 賠率
		if wildCount >= s.minLen {
			wildPay = s.Paytable[int(s.W1Id)][wildCount-1]
		}

		// 6. 取較大者
		winSym := symId
		winCnt := symCount
		winPay := symPay

		if wildPay > symPay {
			winSym = s.W1Id
			winCnt = wildCount
			winPay = wildPay
		}

		// 7. 更新結果
		totalLinePay += winPay
		r.LineResult = append(r.LineResult, LineResult{
			sym:    winSym, // 得分符號
			cnt:    winCnt, // 連線數量
			win:    winPay, // 賠分
			lineId: i,      // 線路 ID
		}) // 更新結果
	}

	// 一次 spin 盤面得結果
	r.Win = totalLinePay * bet / linesLen // 總賠分
	return r
}

// ways 算分模式
func CalcWaysGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	// 未實做
	return s.ScreenResult
}
