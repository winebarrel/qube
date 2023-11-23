Output demo.gif

Require echo

Set Shell "bash"
Set FontSize 12
Set Width 800
Set Height 600
Set Padding 20

Type "jq -c . data.jsonl" Sleep 500ms Enter

Sleep 2s

Type "qube -d 'root@tcp(127.0.0.1:13306)/' -t 3s -n 3 -r 100 -f data.jsonl --noop | jq -C | head -n 25" Sleep 500ms Enter

Sleep 5s

Type "qube -d 'root@tcp(127.0.0.1:13306)/' -t 10s -n 3 -r 100 -f data.jsonl | jq -C | head -n 25" Sleep 500ms Enter

Sleep 12s

Type "qube -d 'root@tcp(127.0.0.1:13306)/' -t 10s -n 10 -r 1000 -f data.jsonl | jq -C | head -n 25" Sleep 500ms Enter

Sleep 12s
