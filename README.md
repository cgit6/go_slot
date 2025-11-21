#  [程式 / Golang] 11/11 判斷得分

[程式計劃書](https://hackmd.io/@chiSean/rkGvg3Wxbl) 
[PAR sheet](https://docs.google.com/spreadsheets/d/1WncTL93uOFgXVq_zC1yJQky_ZsAkpolA5LZoVFo4OuE/edit?usp=sharing)

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



## code 修正

[v] Bet / line 有小數點問題
[] C1 有賠率要計算
    [] 檢視 C1 得分
    [] 檢視 sym 得分
    [] 計算 C1 rtp
    [] 計算 sym rtp
    ! 現在發現問題，C1 會在 line 中被計算，不知道哪裡有問題
[] 執行 Free Game
[] 統計值公式(CV、AllRTP、HitRate、FG trigger)

