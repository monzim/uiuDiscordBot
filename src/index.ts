import {
  Client,
  ClientOptions,
  GatewayIntentBits,
  Message,
  EmbedBuilder,
  Interaction,
} from "discord.js";

import { findCoursesByCodeAndSection } from "./cosmos/cosmos_search";
import { CHANNEL_ID, DISCORD_TOKEN } from "./config";
import { initializeCommands } from "./commands";
import { generateMarkdownMessageWithCourse } from "./models/course";
import {
  sendWebhookErrorMessage,
  sendWebhookInteractionMessage,
  sendWebhookMessage,
  sendWebhookStatusMessage,
} from "./webhook/send_message";
import { installmentHandlr } from "./command_handlr/installment";
import { holidayHandlr } from "./command_handlr/holidayhandlr";
import { makeupHandlr } from "./command_handlr/makeuphandlr";
import { examHandlr } from "./command_handlr/examhandlr";
import { helpHandlr } from "./command_handlr/helphandlr";
import { versionHandlr } from "./command_handlr/versionhandlr";
import { unizimHandlr } from "./command_handlr/unizimhandlr";
import { applinkHandlr } from "./command_handlr/applinkhandlr";
import { examTimeHandlrProd } from "./command_handlr/examtimehandlr";

const botOptions: ClientOptions = {
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.DirectMessages,
    GatewayIntentBits.MessageContent,
  ],
};

const bot = new Client(botOptions);

(async () => {
  await initializeCommands();
  bot.login(DISCORD_TOKEN);

  const prefix = "!";

  bot.on("ready", () => {
    console.log("Bot is online!");
    console.log(`Logged in as ${bot.user}`);

    sendWebhookStatusMessage(
      "Online",
      `Bot is online!. Logged in as ${bot.user}`
    );
    Client;

    // send a message to from bot to the channel
    const channel = CHANNEL_ID as string;
    bot.channels
      .fetch(channel)
      .then((channel) => {
        if (channel != null && channel.isTextBased()) {
          channel
            .send({
              embeds: [
                new EmbedBuilder()
                  .setColor("Green")
                  .setTitle("Wassup. After a new update, I am back online!"),
              ],
            })
            .catch((err) => {
              console.log(
                "🚀 ~ file: index.ts:69 ~ bot.channels.fetch ~ err:",
                err
              );

              sendWebhookErrorMessage("index.ts:69", err);
            });
        }
      })
      .catch((err) => {
        console.log("🚀 ~ file: index.ts:81 ~ bot.on ~ err:", err);

        sendWebhookErrorMessage("index.ts:81", err);
      });
  });

  bot.on("messageCreate", async (message: Message) => {
    if (message.author.bot) return;

    let _logMessage = `id: ${message.id} guildId: <@${message.guildId}> channelId: <#${message.channelId}> : Author: <@${message.author.id}> Message: ${message}`;

    sendWebhookMessage(message);

    if (!message.content.startsWith(prefix)) {
      _logMessage += " Command: " + null;
      console.log(_logMessage);
      return;
    }

    const args = message.content.slice(prefix.length).trim().split(/ +/);

    const command = args.shift()?.toLowerCase();

    _logMessage += " Command: " + command;
    console.log(" 🚀message :", _logMessage);

    // Handle specific commands
    if (command === "uiu") {
      const courseCode = args[0];
      const section = args[1];

      const courses = await findCoursesByCodeAndSection(
        courseCode.toUpperCase(),
        section.toUpperCase()
      );

      if (courses.length === 0) {
        message.channel
          .send({
            content: `${message.author}  No courses found for ${courseCode} ${section}`,
          })
          .catch((err) => {
            console.log("🚀 ~ file: index.ts:82 ~ bot.on ~ err:", err);
            sendWebhookErrorMessage("index.ts:82", err);
          });
        return;
      }

      let title =
        "The following **" + courses.length + "** courses were found\n";
      let listEm: any = [];

      courses.map((course) => {
        listEm.push(
          new EmbedBuilder()
            .setTitle(course.CourseCode)
            .setColor("Random")
            .setDescription(
              `Section: ${course.Section}     Faculty: ${course.Teacher}\n**${course.ExamDate} at ${course.ExamTime}**\n${course.Room}\n`
            )
        );
      });

      if (listEm.length > 10) {
        let msg = generateMarkdownMessageWithCourse(courses);
        let fullMessage = `${message.author} **${courses[0].CourseTitle}\n ${title} ${msg}`;

        const textLimit = 1900;

        if (fullMessage.length > textLimit) {
          const slices: string[] = [];
          while (fullMessage.length > 0) {
            slices.push(fullMessage.slice(0, textLimit));
            fullMessage = fullMessage.slice(textLimit);
          }

          // Send each slice as a separate message with a delay between each message
          const delay = 1000; // milliseconds
          const sendSlices = async () => {
            for (let i = 0; i < slices.length; i++) {
              await new Promise((resolve) => setTimeout(resolve, delay));
              const chunkedMessage = `Message ${i + 1} of ${slices.length}:\n${
                slices[i]
              }`;
              message.channel.send(chunkedMessage).catch((err) => {
                console.log("🚀 ~ file: index.ts:125 ~ sendSlices ~ err:", err);
                sendWebhookErrorMessage("index.ts:125", err);
              });
            }
          };

          sendSlices();
        } else {
          // If the full message is within the character limit, send it as a single message
          message.channel.send(fullMessage).catch((err) => {
            console.log("🚀 ~ file: index.ts:135 ~ bot.on ~ err:", err);
            sendWebhookErrorMessage("index.ts:135", err);
          });
        }
      } else {
        message.channel
          .send({
            content: `${message.author} **${courses[0].CourseTitle}\n ${title}`,
            embeds: listEm,
          })
          .catch((err) => {
            console.log("🚀 ~ file: index.ts:146 ~ bot.on ~ err:", err);
            sendWebhookErrorMessage("index.ts:146", err);
          });
      }
    }
  });

  bot.on("interactionCreate", async (interaction: Interaction) => {
    if (!interaction.isCommand()) return;

    let _logMessage = `id: ${interaction.id} guildId: <@${interaction.guildId}> channelId: <#${interaction.channelId}> : Author: <@${interaction.user}> Message: ${interaction}`;

    console.log(" 🚀interaction: ", _logMessage);
    sendWebhookInteractionMessage(interaction);

    const { commandName } = interaction;

    if (interaction.isChatInputCommand()) {
      if (commandName === "help") {
        await helpHandlr(interaction);
      } else if (commandName === "ping") {
        await interaction.reply("Hey there! I'm alive! :D").catch((err) => {
          sendWebhookErrorMessage("index.ts:167", err);

          interaction.followUp("Unknown command").catch((err) => {
            sendWebhookErrorMessage("index.ts:172", err);
          });
        });
      } else if (commandName === "uiu") {
        await interaction
          .reply(
            "**UIU** Discord Bot created by  <@669529872644833290>\nThe bot is still in development. Please report any bugs to the developer. [Contact](https://monzim.com/monzim)\nThanks for using the bot!"
          )
          .catch((err) => {
            sendWebhookErrorMessage("index.ts:180", err);

            interaction.followUp("Unknown command").catch((err) => {
              sendWebhookErrorMessage("index.ts:184", err);
            });
          });
      } else if (commandName === "exam_time") {
        await examTimeHandlrProd(interaction);
      } else if (commandName === "installment") {
        await installmentHandlr(interaction);
      } else if (commandName === "holiday") {
        await holidayHandlr(interaction);
      } else if (commandName === "makeup") {
        await makeupHandlr(interaction);
      } else if (commandName === "exam") {
        await examHandlr(interaction);
      } else if (commandName === "version") {
        await versionHandlr(interaction);
      } else if (commandName === "unizim") {
        await unizimHandlr(interaction);
      } else if (commandName === "applink") {
        await applinkHandlr(interaction);
      } else {
        await interaction.followUp("Unknown command").catch((err) => {
          console.log("🚀 ~ file: index.ts:274 ~ bot.on ~ err:", err);
          sendWebhookErrorMessage("index.ts:274", err);

          interaction.followUp("Unknown command").catch((err) => {
            console.log("🚀 ~ file: index.ts:278 ~ bot.on ~ err:", err);
            sendWebhookErrorMessage("index.ts:278", err);
          });
        });
      }
    } else {
      await interaction.reply("Unknown command LOL").catch((err) => {
        console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);

        sendWebhookErrorMessage("interaction.isChatInputCommand", err);

        interaction.followUp("Unknown command").catch((err) => {
          console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
        });
      });
    }
  });
})();
