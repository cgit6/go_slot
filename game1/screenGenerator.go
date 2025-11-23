package game

import "math/rand"

type ScreenGenerator struct {
	*Config              // 匿名嵌入 Config
	ScreenBuf []uint8    // 盤面緩存
	rng       *rand.Rand // RNG
}

// 建構函數: 創建 ScreenGenerator instance 時調用
func NewScreenGenerator(cfg *Config, rng *rand.Rand) *ScreenGenerator {

	// 創建 ScreenGenerator instance & 賦值
	return &ScreenGenerator{
		Config:    cfg,                           // 嵌入 Config
		ScreenBuf: make([]uint8, cfg.ScreenSize), // 盤面緩存
		rng:       rng,                           // RNG
	}
}

// 盤面生成
func (g *ScreenGenerator) GenScreen() []uint8 {

	// 對每一軸操作
	for i := 0; i < g.Cols; i++ {
		idx := g.rng.Intn(g.ReelLens[i])
		for j := 0; j < g.Rows; j++ {
			g.ScreenBuf[i*g.Rows+j] = g.ReelStrips[i][(idx+j)%g.ReelLens[i]]
		}
	}

	return g.ScreenBuf
}
