import { EmbedBuilder, Interaction, Message, WebhookClient } from "discord.js";
import { WEBHOOK_ID, WEBHOOK_TOKEN } from "../config";

const webhookClient = new WebhookClient({
  id: WEBHOOK_ID as string,
  token: WEBHOOK_TOKEN as string,
});

export const sendWebhookMessage = async (message: Message<boolean>) => {
  let _logMessage = `id: ${message.id}> Message: ${message}`;

  if (_logMessage.length > 1990) {
    console.warn(
      "Message length exceeds Discord's message length limit of 2000 characters. Truncating..."
    );
    _logMessage = _logMessage.slice(0, 1997) + "...";
  }

  webhookClient
    .send({
      content: `Author: <@${message.author.id}> : channelId: <#${message.channelId}> : guildId: <@${message.guildId}> `,
      embeds: [new EmbedBuilder().setColor("Aqua").setDescription(_logMessage)],
    })

    .catch((error) => {
      console.log(
        "🚀 ~ file: send_message.ts:47 ~ sendWebhookMessage ~ error:",
        error
      );
    });
};

export const sendWebhookInteractionMessage = async (
  interaction: Interaction
) => {
  let _logMessage = `id: ${interaction.id} Message: ${interaction}`;

  if (_logMessage.length > 1990) {
    console.warn(
      "Message length exceeds Discord's message length limit of 2000 characters. Truncating..."
    );
    _logMessage = _logMessage.slice(0, 1997) + "...";
  }

  webhookClient
    .send({
      content: `${interaction.user} : channelId: <#${interaction.channelId}> : guildId: <@${interaction.guildId}>`,
      embeds: [
        new EmbedBuilder()
          .setColor("DarkNavy")

          .setDescription(_logMessage),
      ],
    })

    .catch((error) => {
      console.log(
        "🚀 ~ file: send_message.ts:47 ~ sendWebhookMessage ~ error:",
        error
      );
    });
};

export const sendWebhookErrorMessage = async (func: string, err: any) => {
  let errorMessage = err.toString();

  if (errorMessage.length > 1990) {
    console.warn(
      "Error message length exceeds Discord's message length limit of 2000 characters. Truncating..."
    );
    errorMessage = errorMessage.slice(0, 1997) + "...";
  }

  webhookClient
    .send({
      content: `❌ Error Occurred ❌ : ${func}`,
      embeds: [new EmbedBuilder().setColor("Red").setDescription(errorMessage)],
    })
    .catch((error) => {
      console.error("Failed to send error message:", error);
    });
};

export const sendWebhookStatusMessage = async (
  status: string,
  message: string
) => {
  let statusMessage = message.toString();

  if (statusMessage.length > 1990) {
    console.warn(
      "Status message length exceeds Discord's message length limit of 2000 characters. Truncating..."
    );
    statusMessage = statusMessage.slice(0, 1997) + "...";
  }

  webhookClient
    .send({
      content: `Status: ${status}`,
      embeds: [
        new EmbedBuilder().setColor("Green").setDescription(statusMessage),
      ],
    })
    .catch((error) => {
      console.error("Failed to send status message:", error);
    });
};
