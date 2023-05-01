import { CacheType, ChatInputCommandInteraction } from "discord.js";
import { sendWebhookErrorMessage } from "../webhook/send_message";

export const makeupHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  // Jul 24, 2023 Mon Make-up class: Regular Saturday Classes
  // Aug 7, 2023 Mon Make-up class: Regular Tuesday Classes
  // Aug 21, 2023 Mon Make-up class: Regular Saturday Classes

  await interaction
    .reply({
      content:
        `List of makeup classes ☹️. UIU ruined our holidays 😭\n\n` +
        `1. Regular **Saturday** Classes on **${new Date(
          "Jul 24, 2023"
        ).toDateString()}**\n` +
        `2. Regular **Tuesday** Classes on **${new Date(
          "Aug 7, 2023"
        ).toDateString()}**\n` +
        `3. Regular **Saturday** Classes on **${new Date(
          "Aug 21, 2023"
        ).toDateString()}**\n`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: makeuphandlr.ts:26 ~ err:", err);

      sendWebhookErrorMessage("makeupHandlr", err);
      interaction.followUp("Error occurred while executing the command!");
    });
};
