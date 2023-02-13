# How to run 

## 1. Install Go

https://go.dev/doc/install 

## 2. Start the program 
```javascript
go run main.go -demo "{C:/users/full/path/to/match.dem}" -suspect {steamID64}
```

# Usage example

```javascript
go run main.go -demo "C:\Users\name\Desktop\DEMO\match730_003593207023228944468_0353560675_181.dem" -suspect 76561198418101905 -prod
```

( `-prod` is optional for production / API calls)

# Result
### Should look somewhat like this:
```javascript
 {
  "match_id": "467b5cececbd3e8decdb56a6ec3198c8b58d6194765eea0c22d78feb214ed62d",
  "steam_id": 76561198418101900,
  "player_id": 13,
  "stat_result": 0,
  "automation_result": 0,
  "data": {
    "kills": 22,
    "hs": 27.272728,
    "wall": 2,
    "smoke": 0,
    "flashed": 0,
    "rank": 14,
    "rounds": 29,
    "duration": 2840.828051456
  }
}
```
