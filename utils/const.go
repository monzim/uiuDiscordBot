package utils

const (
	SUPPORT_MESSAGE = "**Support the bot's longevity with a donation if you found it helpful. [Click here for Donate 🎁](https://monzim.com/support)**\n**Thank you!**"

	ERROR_MESSAGE = "Oops! Something went wrong. Please try again later."
	BOT_LOGO      = "https://cdn-monzim-com.azureedge.net/public-com/public/fa9b22c0-dc41-11ee-a977-8d8dd057b6ad"
	BUY_ME_COFFEE = "https://res.cloudinary.com/monzim/image/upload/v1688984685/download_kh1syl.png"
)

func Bold(s string) string {
	return "**" + s + "**"
}
