package main

import "log"

// 這邊就會遇到不同的算分方式就會用到不同的參數有的算分方式不會用到某些屬性

type ScreenResult struct {
	Screen    []uint8 // 原始盤面 (一維 symbolID)
	C1Count   int     // 盤面中 C1 (scatter) 出現次數
	TotalPay  int     // 線獎賠率總和
	TotalWins int     // 最終贏分 (TotalLinePay * lineBet)
	TotalBets int     // 累積總下注
}

// 建構函數: 創建 NewScreenResult instance 時調用
func NewScreenResult() *ScreenResult {
	r := &ScreenResult{
		Screen:    []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 這樣寫容易錯
		C1Count:   0,
		TotalPay:  0,
		TotalWins: 0,
		TotalBets: 0,
	}
	return r
}

// 輸入 盤面與 1 次 spin 下注分數
type CalcFunc func(*SpinCalculator, []uint8, int) *ScreenResult // 接收 *SpinCalculator

type SpinCalculator struct {
	*Config                // 匿名嵌入
	*ScreenResult          // 結果緩存
	calcFn        CalcFunc // 算分函數
}

// 建構函數: 創建 NewSpinCalculator instance 時調用
func NewSpinCalculator(cfg *Config, screenResult *ScreenResult) *SpinCalculator {
	if screenResult == nil {
		screenResult = &ScreenResult{}
	}
	sc := &SpinCalculator{Config: cfg, ScreenResult: screenResult}
	sc.initCalcFn()
	return sc
}

// 維護一個map註冊表
var calcFnMap = map[GameMode]CalcFunc{
	ModeLines: CalcLinesGame, // lines 算法
	ModeWays:  CalcWaysGame,  // ways 算法

	// 特殊符號

}

// 選擇計算方式
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

// 對外統一調用 CalcScreen：不管是 Line 還是 Ways，都呼叫 CalcScreen
func CalcScreen(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	return s.calcFn(s, screen, bet)
}

// ------- 不同算分模式的內部函數 -------
func CalcLinesGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {

	// 1. 計算每一條線的 C1 / wildCount / symId / symCount / pay

	for i := 0; i < s.ScreenSize; i++ {

	}

	// 3.1. TotalWin = C1 賠率 * Bet

	// 3.2. TotalWin = TotalLinePay * Bet / 線數

	return s.ScreenResult
}

func CalcWaysGame(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	// 未實做
	return s.ScreenResult
}
