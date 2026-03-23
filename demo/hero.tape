Output demo/hero.gif

Require ppp

Set Shell zsh
Set FontSize 16
Set Width 800
Set Height 500
Set Theme "Catppuccin Mocha"
Set Padding 20
Set TypingSpeed 30ms

Type `echo '{"name":"Alice","active":true,"age":30,"roles":["admin","user"]}' | ppp`
Enter
Sleep 2s

Type `printf 'name,age,city\nAlice,30,NYC\nBob,25,LA\nCharlie,35,Chicago' | ppp`
Enter
Sleep 2s

Type `printf 'host: localhost\nport: 8080\ndebug: true' | ppp`
Enter
Sleep 2s

Type `echo 'user=john%40example.com&token=abc123&active=true' | ppp`
Enter
Sleep 2s

Type `printf '2024-01-15T10:30:00Z [ERROR] connection refused: db timeout\n2024-01-15T10:30:05Z [INFO] retry succeeded' | ppp`
Enter
Sleep 2s

Sleep 1s
