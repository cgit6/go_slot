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

// 輸入 盤面與 1 次 spin 下注分數
type CalcFunc func(*SpinCalculator, []uint8, int) *ScreenResult // 接收 *SpinCalculator

type SpinCalculator struct {
	*Config                // 匿名嵌入
	*ScreenResult          // 結果緩存
	calcFn        CalcFunc // 算分函數
}

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
	ModeLines: (*SpinCalculator).calcLinesGame, // lines 算法
	ModeWays:  (*SpinCalculator).calcWaysGame,  // ways 算法
}

// 選擇計算方式
func (s *SpinCalculator) initCalcFn() {

	// 選擇算分策略
	if fn, ok := calcFnMap[s.Mode]; ok {
		s.calcFn = fn
	}
	log.Fatal("未知 mode")
	// panic 表示還有救，類似 try ... catch ...

}

// 對外統一調用 CalcScreen：不管是 Line 還是 Ways，都呼叫 CalcScreen
func CalcScreen(s *SpinCalculator, screen []uint8, bet int) *ScreenResult {
	return s.calcFn(s, screen, bet)
}

// ------- 不同算分模式的內部函數 -------
func (s *SpinCalculator) calcLinesGame(screen []uint8, bet int) *ScreenResult {
	// TODO: 填入 lines 計分，寫入 s.ScreenResult 後回傳
	return s.ScreenResult
}

func (s *SpinCalculator) calcWaysGame(screen []uint8, bet int) *ScreenResult {
	// TODO: 填入 ways 計分，寫入 s.ScreenResult 後回傳
	return s.ScreenResult
}
