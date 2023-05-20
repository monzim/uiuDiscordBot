import { CacheType, ChatInputCommandInteraction } from "discord.js";

export const versionHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  await interaction
    .reply({
      content: `UIU Bot **v.1.0.6**\n **Summer 2023** information. \n**Developer** <@669529872644833290>\n\n**Features**\n- Unizim\n- applink\n- Exam Time\n- Installment\n- Holiday\n- Makeup\n- Exam\n- Version`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: versionhandlr.ts:9 ~ err:", err);
    });
};
