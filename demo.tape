# https://github.com/charmbracelet/vhs
Output demo.gif

Require echo

Set Shell "bash"
Set FontSize 12
Set Width 800
Set Height 600
Set Padding 20

Type `echo '{"q":"select 1"}' > data.jsonl` Sleep 100ms Enter
Type `echo '{"q":"select 2"}' >> data.jsonl` Sleep 100ms Enter
Type `echo '{"q":"select 3"}' >> data.jsonl` Sleep 100ms Enter
Type "jq -c . data.jsonl" Sleep 500ms Enter

Sleep 2s

Type "qube -d 'root@tcp(127.0.0.1:13306)/' -t 3s -n 3 -r 100 -f data.jsonl --noop -C | head -n 25" Sleep 500ms Enter

Sleep 5s

Type "qube -d 'root@tcp(127.0.0.1:13306)/' -t 10s -n 3 -r 100 -f data.jsonl -C | head -n 25" Sleep 500ms Enter

Sleep 12s

Type "qube -d 'root@tcp(127.0.0.1:13306)/' -t 10s -n 10 -r 1000 -f data.jsonl -C | head -n 25" Sleep 500ms Enter

Sleep 12s
