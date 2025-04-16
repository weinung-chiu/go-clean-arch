# Go Clean Architecture
# 目標
* 社群友善：以實務為出發點，避免過多獨創的設計
* 在一定的限制下，各團隊可以依實際需求 (API、Game) 設計並實作需要的服務
* 開發者友善：只需要照著範例寫即可寫出 Golang 風格、易維護的 Code
* 明確、可追蹤的程式碼


# Design
## Layout 
依 standard project layout

##  clean arch
MUST - 正確的依賴方向

從 request 進入的方向來看， router -> application -> adapter ，4(or 3) layer，依各團隊需求而定
application 是業務核心


## `log/slog` (WIP)

# Get Started (WIP)
WIP
- generate template

# How To
## Logging