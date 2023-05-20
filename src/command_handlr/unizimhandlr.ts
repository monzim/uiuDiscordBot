import { CacheType, ChatInputCommandInteraction } from "discord.js";

export const unizimHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  await interaction
    .reply({
      content: `Unizim App Public Beta Testing.\n\nA cross-platform scalable app designed specifically for UIU. Self-managed backend, which would give UIU complete control over the app's data and infrastructure.\n\n Join the beta testing program: https://appdistribution.firebase.dev/i/e0e2dffda4801a88`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: unizimhandlr.ts:11 ~ err:", err);
    });
};
