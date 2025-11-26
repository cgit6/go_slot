#  [程式 / Golang] 11/11 判斷得分

[程式計劃書](https://hackmd.io/@chiSean/rkGvg3Wxbl) 
[PAR sheet](https://docs.google.com/spreadsheets/d/1ZZUYusikSiQDiVpgFvjWn6BTJfpeQmHp8YFo_HBsEkw/edit?usp=sharing)

<!-- ## 執行

### 下載專案

``` shell=

```

### 執行專案

``` shell=

``` -->

## 完成進度

[v]Config struct  
    [v]struct
    [v]Init
    [v]Reset
    [v]validate
    [v]NewConfig

[v]ScreenGenerator struct
    [v]struct
    [v]GenScreen
    [v]NewScreenGenerator

[v]LineResult struct
    [v]struct
[v]ScreenResult struct
    [v]struct
[-]SpinCalculator struct
    [v]struct
    [v]calcFnMap
    [v]initCalcFunc
    [v]CalcScreen
    [v]CalcLineGame
    [x]CalcWaysGame

[v]runner func
    [v]runner

## 錯誤的種類
設定擋案錯誤 -> log.Fatal()


## 搞懂
[] 計算 deriveFilter 函數的運算邏輯
[] 轉換 Config Init LINES -> FlatLines 

## 文件檢查清單
Base Game
[]輪帶表
    []輪帶表內容
    []輪帶長度
    []權重值
    []權重上方的輪帶長度
[]Slice 表
[]C1 符號數量表
[]C1 組合數
    []組合數的排列/數量/組合數
    []組合機率
        []數量/Hit./Prob./pay/RTP
    []驗算
        []每一個數量
        []顆數
[]符號數量
    []ANY 符號 == 輪帶長度
    []下方 符號等於 Slice 出現數量
    []加總等於ANY 數量
    []#符號數量等於 輪帶長度 - 符號數量
    []符號數量+W1 等於上方兩者相加
    []#符號數量+W1 等於 輪帶長度 - (符號數量+W1)
[]組合計算表
    []set 代表
    []組合
    []連線符號
    []連線數量
    []賠率
    []Hit.
    []Prob.
    []RTP

Free Game
[]輪帶表
    []輪帶表內容
    []輪帶長度
    []權重值
    []權重上方的輪帶長度
[]Slice 表
[]C1 符號數量表
[]C1 組合數
    []組合數的排列/數量/組合數
    []組合機率
        []數量/Hit./Prob./pay/RTP
    []驗算
        []每一個數量
        []顆數
[]符號數量
    []ANY 符號 == 輪帶長度
    []下方 符號等於 Slice 出現數量
    []加總等於ANY 數量
    []#符號數量等於 輪帶長度 - 符號數量
    []符號數量+W1 等於上方兩者相加
    []#符號數量+W1 等於 輪帶長度 - (符號數量+W1)
[]組合計算表
    []set 代表
    []組合
    []連線符號
    []連線數量
    []賠率
    []Hit.
    []Prob.
    []RTP
    []W1 乘數

Game Summery
[]加局計算
    []trigger
    []re-reigger
    []RTP
    []平均局數

* 文件那邊我發現如果輪帶太短，平均局數會出現負數的問題

## code 修正

[v] Bet / line 有小數點問題
[v] C1 有賠率要計算
    [v] 檢視 C1 得分
    [v] 檢視 sym 得分
    [v] 計算 C1 rtp
    [v] 計算 sym rtp
    [v] ! 現在發現問題 C1 會在 line 中被計算，不知道哪裡有問題
[] 執行 Free Game
    [v] 判斷觸發條件對應局數 {3,4,5} C1 -> {10,12,15} 寫在設定欄位中
    [v] wild multiple x3 不累積 -> WinDetail 添加一個這條線中獎段中有沒有包含 W1 的 bool 值，用於處理 multiple feature 看來我只能在 算分函數中處理這件事了
    [v] 計算初始 10 局 FG spin
    [v] re-trigger 追加局數(所以不能用 for 的概念要用 while 的寫法)
    [v] 返回結果 rtp，計算統計值

    [v] 如果沒有 W1 的時候數值是否依舊正確?
        [v] 如果 乘數是 1 的情況下數值誤差極小
        [] 如果沒有 W1 數值依然正確 

[] 統計值公式(CV、RTP、得分率、FG 觸發率)
    [v] CV 值計算 樣本標準差 / 樣本平均值
        [v] 計算 樣本標準差，需要樣本數、樣本平均值、樣本值。
        std = sqrt((樣本值 - 樣本平均值)^2 / (樣本數 - 1))
        [v] 計算 rtp 
        [-] only base game cv -> 看 bg 平不平順
        [-] only free game cv -> 看 fg 平不平順
    [v] RTP all Wins / all Bet
    [-] 得分率 1 spin > 0 次數 / 總模擬次數
        [v] 支付 spin 為單位
        [] 獨立一局 spin 為單位
    [v] FG 觸發率 執行 FG 次數/ 總模擬次數
    
[] CV= 9
[] rtp= 96-97%
[] hit rate= 30%-40%
[] trigger= 0.83%