<div align="center">

<a href="https://t.me/Unrated_Coder">
<img src="https://imgyx.pages.dev/DE8Ni" alt="Telegram MediaInfo Bot Banner" width="100%" style="border-radius: 14px; box-shadow: 0 10px 30px rgba(0,0,0,0.5);" /></a>

  <h1>⚡ Telegram MediaInfo Engine ⚡</h1>

  <a href="https://t.me/Unrated_Coder">
    <img src="https://readme-typing-svg.demolab.com?font=Plus+Jakarta+Sans&weight=700&size=20&duration=2500&pause=1000&color=26A5E4&center=true&vCenter=true&width=500&height=40&lines=%E2%9A%A1+Instant+Range-Based+Streaming;%F0%9F%9A%80+Pure+Go+%2B+MTProto+v2.0;%F0%9F%92%A1+Ultra-Low+Memory+Footprint;%F0%9F%93%84+Sub-Second+Probing+Engine" alt="Typing Animation" />
  </a>

  <p>
    <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" /></a>
    <a href="https://t.me/Unrated_Coder"><img src="https://img.shields.io/badge/Telegram-MTProto%20v2.0-26A5E4?style=for-the-badge&logo=telegram&logoColor=white" alt="Telegram" /></a>
    <a href="https://www.mongodb.com"><img src="https://img.shields.io/badge/Database-MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white" alt="MongoDB" /></a>
    <a href="https://github.com/Unrated-Coder/Mediainfo-Bot/blob/master/LICENSE"> 
    <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" alt="License" /></a>
  </p>

  <h4><i>Experience premium, elite, sub-second technical media extraction without the overhead of full file downloads.</i></h4>

  <p>
    <a href="#-elite-features"><b>Elite Features</b></a> •
    <a href="#%EF%B8%8F-quick-one-click-deployments"><b>Cloud Deploy</b></a> •
    <a href="#-configuration"><b>Configuration</b></a> •
    <a href="#-command-matrix"><b>Commands</b></a>
  </p>

</div>

---

<h2 id="-elite-features">🌟 Elite Features</h2>

- **⚡ Sub-Second Range Streaming:** Employs a custom local streaming server (`pkg/streamer`) translating arbitrary range requests into 4KB aligned chunks fetched via Telegram's `UploadGetFile` MTProto API for instant header probing.
- **🚀 Pure Go MTProto v2.0:** Built completely in raw Go on top of the ultra-high performance [`gotd/td`](https://github.com/gotd/td) library, running with negligible idle memory (~15MB) and zero CPU footprint.
- **📱 Premium Navigation Dashboard:** Completely interactive `/admins` control panel with live-updating inline colored buttons. Instantly delete admin credentials or refresh stats with zero commands.
- **🛠️ Auto-Configuration on Startup:** Auto-registers all bot slash commands dynamically with the Telegram API upon launch, including customized permissions and tags.
- **📝 Clean Small-Caps UI:** All user-facing menus, help screens, and interactive messages are styled elegantly in lowercase small caps.
- **📄 Smart Overflow Fallback:** Dynamically measures the output payload. If the details exceed Telegram's 4,096-character mark, it instantly compiles and delivers a formatted `.txt` report.

---

<h2 id="️-quick-one-click-deployments">🚀 One-Click Cloud Deployments</h2>

Deploy your own private high-speed MediaInfo Engine instantly using any of the premium templates below:

| Cloud Provider | Deploy Template Link | Button |
| :--- | :--- | :---: |
| **Koyeb** | [Deploy to Koyeb](https://deploy.koyeb.com/?type=git&repository=github.com/Unrated-Coder/Mediainfo-Bot&branch=master&name=mediainfo-bot) | [![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://deploy.koyeb.com/?type=git&repository=github.com/Unrated-Coder/Mediainfo-Bot&branch=master&name=mediainfo-bot) |
| **Render** | [Deploy to Render](https://render.com/deploy?repo=https://github.com/Unrated-Coder/Mediainfo-Bot) | [![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/Unrated-Coder/Mediainfo-Bot) |
| **Heroku** | [Deploy to Heroku](https://www.heroku.com/deploy/?template=https://github.com/Unrated-Coder/Mediainfo-Bot) | [![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://www.heroku.com/deploy/?template=https://github.com/Unrated-Coder/Mediainfo-Bot) |
| **Railway** | [Deploy on Railway](https://railway.app/new) | [![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new) |

---

<h2 id="-configuration">⚙️ Environment Configuration</h2>

Configure the following environment variables in your environment or cloud deployment panel:

| Variable | Required | Default | Description |
| :--- | :---: | :---: | :--- |
| `API_ID` | **YES** | - | Your Telegram API ID from [my.telegram.org](https://my.telegram.org) |
| `API_HASH` | **YES** | - | Your Telegram API Hash from [my.telegram.org](https://my.telegram.org) |
| `BOT_TOKEN` | **YES** | - | Your Telegram Bot Token from [@BotFather](https://t.me/BotFather) |
| `MONGO_URI` | **YES** | - | MongoDB Connection URI (`mongodb+srv://...`) |
| `OWNER_ID` | **YES** | - | Telegram User ID of the master Bot Owner |

---

<h2 id="-command-matrix">🎮 Command Matrix</h2>

All command descriptions are beautifully registered in lowercase small-caps with clear usage indicators:

| Command | Permission | Description |
| :--- | :---: | :--- |
| `/start` | 🌐 Everyone | Initialize the session and view welcome status |
| `/help` | 🌐 Everyone | View usage instructions & technical details |
| `/about` | 🌐 Everyone | Read system information and technologies used |
| `/users` | 👑 Admin/Owner | Show total registered users statistics `(admin)` |
| `/ban <user_id>` | 👑 Admin/Owner | Ban a user from using the bot `(admin)` |
| `/unban <user_id>` | 👑 Admin/Owner | Unban a banned user `(admin)` |
| `/admins` | 👑 Admin/Owner | Show interactive color dashboard of admins `(admin)` |
| `/add_admin <user_id>` | ⚡ Owner | Grant administrative privileges `(owner)` |
| `/remadmin <user_id>` | ⚡ Owner | Remove admin privileges `(owner)` |

---

## 🏗️ Self-Hosted Local Setup

If you prefer running the bot locally on your Linux or macOS machine:

### 1. Install System Dependencies
```bash
# Debian / Ubuntu
sudo apt update && sudo apt install -y ffmpeg mediainfo

# macOS Homebrew
brew install ffmpeg mediainfo
```

### 2. Configure Environment & Launch
```bash
export API_ID="your_api_id"
export API_HASH="your_api_hash"
export BOT_TOKEN="your_bot_token"
export MONGO_URI="mongodb://localhost:27017"
export OWNER_ID="your_owner_id"

go run cmd/bot/main.go
```

---

## 👑 Developer & Credits

* **Core Architect & Developer:** [Unrated Coder](https://t.me/Unrated_Coder)
* **Libraries:** [`gotd/td`](https://github.com/gotd/td) • [`MediaInfo`](https://mediaarea.net/en/MediaInfo) • [`MongoDB Go Driver`](https://go.mongodb.org/mongo-driver)

---

<div align="center">
  <p><b>Built with passion and high-performance Go programming.</b></p>
</div>
