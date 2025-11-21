package game

import (
	"fmt"
	"log"
)

// 細項
type WinDetail struct {
	win    int    // 得分
	symbol uint8  // 得分圖標
	hitmap []bool // 中獎圖

	lineId int // 得分線 (Line)
	length int // 長度 (Line Way)
	comb   int // 組合數 (Way)
	cnt    int // 數量(Cluster Count)
}

// 一次 spin 的結果
type ScreenResult struct {
	C1Count    int         // 盤面中 C1 (scatter) 出現次數
	Win        int         // 累積賠分
	Mode       GameMode    // 算分模式
	WinDetails []WinDetail // 細項

	// 輔助參數
	C1Win  int // C1 賠分總和
	SymWin int // 符號賠分總和
}

// input SpinCalculator、screen 與 1 次 spin 下注分數
type CalcFunc func(*SpinCalculator, []uint8, int) *ScreenResult // 接收 *SpinCalculator

type SpinCalculator struct {
	cfg    *Config       // 匿名嵌入
	sr     *ScreenResult // 結果緩存
	calcFn CalcFunc      // 算分函數

	// 輔助參數
	filter uint64 // 特殊符號
}

// 不計分符號清單
func deriveFilter(pay [][]int, wildID uint8) uint64 {
	out := uint64(0) // 0x00000000000000
	// fmt.Println("-------------------------------------------------")
	// fmt.Println("out:", out)
	for sid, row := range pay {
		// fmt.Println("sid:", sid)
		// fmt.Println("row:", row)
		allZero := true
		for _, p := range row {
			// fmt.Println("p:", p)
			if p != 0 {
				allZero = false
				break
			}
		}
		if allZero && uint8(sid) != wildID {
			out |= 1 << uint64(sid)
		}
	}
	return out
}

// 建構函數: 創建 NewSpinCalculator instance 時調用
func NewSpinCalculator(cfg *Config) *SpinCalculator {
	sc := &SpinCalculator{
		cfg: cfg,
		sr:  &ScreenResult{}, // 結果緩存
	}
	sc.initCalcFn()                                         // 初始化算分方式
	sc.sr.WinDetails = make([]WinDetail, 0, len(cfg.Lines)) // 預先建立走線表空間
	sc.filter = deriveFilter(cfg.Paytable, cfg.W1Id)        // 使用bitmask判斷是否得分符號
	fmt.Println("Non-scoring symbols filter:", sc.filter)
	return sc
}

// 選擇算分方式
func (s *SpinCalculator) initCalcFn() {

	// 選擇算分策略
	if fn, ok := calcFnMap[s.cfg.Mode]; ok {
		s.calcFn = fn // 選擇算分方式存到 s.calcFn

		return // 必要，不然會往外跳執行 log.Fatal("未知 mode")
	}
	log.Fatal("未知 mode")
	// panic 表示還有救，但這個沒救了(設定檔錯誤)，類似 try ... catch ...

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

// 維護一個map註冊表
var calcFnMap = map[GameMode]CalcFunc{
	ModeLines: CalcLinesGame, // lines 算法
	ModeWays:  CalcWaysGame,  // ways 算法

}

// ------- 不同算分模式的內部函數 -------

// lines 算分模式
func CalcLinesGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {

	// 初始化 ScreenResult
	r := s.sr

	r.C1Count, r.SymWin, r.C1Win, r.Win = 0, 0, 0, 0 // 重置結果
	linesLen := len(s.cfg.Lines)                     // 線路數量
	r.WinDetails = r.WinDetails[:0]                  // 清空邏輯長度，保留原指針與空間

	totalLinePay := 0 // 累積線路賠分

	// 計算 C1 出現次數
	r.C1Count = countSymbol(screen, s.cfg.C1Id)

	if r.C1Count >= s.cfg.minLen {
		r.C1Win = s.cfg.Paytable[int(s.cfg.C1Id)][r.C1Count-1] * bet // C1 賠分總和
	}

	// 逐條線計分
	for i := 0; i < linesLen; i++ {
		// 單條線的狀態
		wildCount := 0
		wildContinue := true

		var symId uint8     // 得分符號ID
		symStarted := false // 是否已確定得分符號
		symCount := 0       // 符號連線數量

		// 從左到右掃這條線
		for j := 0; j < s.cfg.Cols; j++ {

			// 1. 獲取該位置符號
			sid := screen[s.cfg.FlatLines[i*s.cfg.Cols+j]] // 平坦化線路清單

			// 2. 連續 Wild 數
			if wildContinue && sid == s.cfg.W1Id {
				wildCount++
			} else {
				wildContinue = false
			}

			// 3. 計算得分符號連線

			// 3.1. 尚未決定得分符號
			if !symStarted {
				if sid == s.cfg.W1Id {
					continue
				}

				// ❗ Scatter 一律不能當線獎符號
				if sid == s.cfg.C1Id {
					break
				}

				// 第一個非 Wild：若是不計分符號（Z1/C1 等），此線只能靠純 Wild
				if s.filter&(1<<uint64(sid)) != 0 {
					break
				}
				// 合法得分符號確立
				symId = sid
				symStarted = true
				symCount = wildCount + 1 // 包含前面的 Wild
				continue
			}

			// 3.2. 已決定得分符號，延伸連線：同符號或 Wild 都可
			if sid == symId || sid == s.cfg.W1Id {
				symCount++
			} else {
				break // 如果開頭直接是 C1 直接結束
			}
		}

		// 4. 未達最小連線長度 → 0 分 該條線沒中
		if symCount < s.cfg.minLen && wildCount < s.cfg.minLen {
			continue
		}

		// 5. 計算兩種賠率

		// 5.1. 獲取得分符號賠率
		symPay := 0                                 // 得分符號賠率
		if symStarted && symCount >= s.cfg.minLen { // 只做「是否該算」的必要判斷
			symPay = s.cfg.Paytable[int(symId)][symCount-1] // 獲取賠率
		}

		// 5.2. 獲取 Wild 賠率
		wildPay := 0 // W1 賠率
		if wildCount >= s.cfg.minLen {
			wildPay = s.cfg.Paytable[int(s.cfg.W1Id)][wildCount-1] // 獲取賠率
		}

		// 6. 取較大者
		winSym := symId    // 得分符號
		winCnt := symCount // 連線數量
		winPay := symPay   // 線獎賠率

		// 6.1. 如果 Wild 賠率較大
		if wildPay > symPay {
			winSym = s.cfg.W1Id
			winCnt = wildCount
			winPay = wildPay
		}

		// 7. 更新結果
		totalLinePay += winPay
		r.WinDetails = append(r.WinDetails, WinDetail{
			symbol: winSym, // 得分符號
			cnt:    winCnt, // 連線數量
			win:    winPay, // 賠分
			lineId: i,      // 線路 ID
		}) // 更新結果
	}

	// 一次 spin 盤面贏分結果
	r.SymWin = totalLinePay * bet / linesLen // 符號賠分總和
	r.Win = r.C1Win + r.SymWin               // 總賠分(sym + c1)
	return r
}

// ways 算分模式
func CalcWaysGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	// 未實做
	return s.sr
}
