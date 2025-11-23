#  [程式 / Golang] 11/24 判斷得分

## 規則


[] 掉落消除
[] 任意位置支付

先做 base game 就好。

config struct 需要
符號清單(Symbols): H1-H4,L1-L5,C1,W1
輪帶表(ReelStrips): 
賠率表(Paytable)
盤面大小參數(Rows,Cols,ScreenSize)
算分模式(Mode)
每軸輪帶長度(ReelLens)
C1 代號(C1Id)
W1 代號(W1Id)
最小數量(minCount)

初始化狀態(initFlag)


