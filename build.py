import os
import shutil

version = "1.0.0"

os.environ["GOOS"] = "linux"
os.environ["GOARCH"] = "amd64"

try:
    os.remove("./build/tg-to-rss-bot")
except Exception as e:
    # There's no such file
    pass

os.system("go build -o tg-to-rss-bot main.go")
shutil.move("./tg-to-rss-bot", "build")

os.chdir("build")
os.system(f"docker build -t kirtis/tg-to-rss-bot:{version} .")
