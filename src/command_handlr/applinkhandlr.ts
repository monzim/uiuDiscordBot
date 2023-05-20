import { CacheType, ChatInputCommandInteraction } from "discord.js";

export const applinkHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  await interaction
    .reply({
      content: `**Unizim App Public Beta Testing is Open Now.**\n\nJoin the beta testing program: https://appdistribution.firebase.dev/i/e0e2dffda4801a88 \nParticipate in the survey: https://forms.gle/w9Ec3drbNw5nQy7n9`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: applinkhandlr.ts:11 ~ err:", err);
    });
};
