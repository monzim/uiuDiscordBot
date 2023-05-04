import { CacheType, ChatInputCommandInteraction } from "discord.js";

export const versionHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  await interaction
    .reply({
      content: `UIU Bot **v.1.0.5**\n **Summer 2023** information\n`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: versionhandlr.ts:9 ~ err:", err);
    });
};
