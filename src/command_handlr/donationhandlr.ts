import { CacheType, ChatInputCommandInteraction } from "discord.js";

export const donationHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  await interaction
    .reply({
      //   content: `### **Support the bot's longevity with a donation if you found it helpful. By contributing to our bot's hosting expenses, you directly enable us to keep the service running smoothly, ensuring continuous accessibility for students who rely on it. [Click here for 🎁 Donate](https://monzim.com/support).Thank you!\n**Bkash**\n- **01303382422**\n**Rocket**\n- **018415003604**`,
      content: `Support the bot's longevity with a donation if you found it helpful. By contributing to our bot's hosting expenses, you directly enable us to keep the service running smoothly, ensuring continuous accessibility for students who rely on it. [Click here for 🎁 Donate](https://monzim.com/support).Thank you!`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: versionhandlr.ts:9 ~ err:", err);
    });
};
