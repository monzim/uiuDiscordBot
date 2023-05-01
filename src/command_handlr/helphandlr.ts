import {
  CacheType,
  ChatInputCommandInteraction,
  ColorResolvable,
  EmbedBuilder,
} from "discord.js";
import { sendWebhookErrorMessage } from "../webhook/send_message";

type uiuCommand = {
  command: string;
  description: string;
};

export const helpHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  /**
   /ping - Replies with Pong!
   /exam - Replies with upcoming exam!
   /exam_time - Get course exam info like date, time, room, etc.
   /makeup - Get upcoming makeup class info!
   /installment - Get info about upcoming installment!
   /holiday - Replies with upcoming holiday!
   /uiu - About the UIU Discord Bot
   /help - Replies with list of commands!
 */
  const commandList: uiuCommand[] = [
    {
      command: "/ping",
      description: "Replies with Pong!",
    },
    {
      command: "/exam",
      description: "Replies with upcoming exam!",
    },
    {
      command: "/exam_time",
      description: "Get course exam info like date, time, room, etc.",
    },
    {
      command: "/makeup",
      description: "Get upcoming makeup class info!",
    },
    {
      command: "/installment",
      description: "Get info about upcoming installment!",
    },
    {
      command: "/holiday",
      description: "Replies with upcoming holiday!",
    },
    {
      command: "/uiu",
      description: "About the UIU Discord Bot",
    },
    {
      command: "/help",
      description: "Replies with list of commands!",
    },
  ];

  await interaction
    .reply({
      content: "Here is the list of commands:",
      embeds: [
        new EmbedBuilder()
          .setColor("NotQuiteBlack")
          .setTitle("UIU Discord Bot Commands")
          .setDescription(
            commandList
              .map((cmd) => {
                let command = "`";
                command += `${cmd.command}`;
                command += "`";

                return `**${cmd.command}**\nAbout: ${cmd.description}\nUsage: ${command}\n`;
              })
              .join("\n")
          ),
      ],
    })
    .catch((err) => {
      console.log("🚀 ~ file: examhandlr.ts:71 ~ err:", err);

      sendWebhookErrorMessage("examhandlr", err);
      interaction.followUp("Error occurred while executing the command!");
    });
};
