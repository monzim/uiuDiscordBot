import {
  Client,
  ClientOptions,
  GatewayIntentBits,
  Message,
  EmbedBuilder,
  Interaction,
} from "discord.js";

import { findCoursesByCodeAndSection } from "./cosmos_search";
import { discordToken } from "./config";
import { initializeCommands } from "./commands";

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
  bot.login(discordToken);

  const prefix = "!";

  bot.on("ready", () => {
    console.log("Bot is online!");
    console.log(`Logged in as ${bot.user}`);
  });

  bot.on("messageCreate", async (message: Message) => {
    if (message.author.bot) return;

    console.log("Message received" + message);

    // log author and message
    console.log("Author: " + message.author.id);

    if (!message.content.startsWith(prefix)) {
      console.log("Message does not start with prefix");
      return;
    }

    const args = message.content.slice(prefix.length).trim().split(/ +/);

    const command = args.shift()?.toLowerCase();
    console.log(`Command: ${command}`);

    // Handle specific commands
    if (command === "uiu") {
      const courseCode = args[0];
      const section = args[1];

      const courses = await findCoursesByCodeAndSection(
        courseCode.toUpperCase(),
        section.toUpperCase()
      );

      if (courses.length === 0) {
        message.channel.send({
          content: `${message.author}  No courses found for ${courseCode} ${section}`,
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
            .setColor("Random") // Set the color of the embed
            .setDescription(
              `Section: ${course.Section}     Faculty: ${course.Teacher}\n**${course.ExamDate} at ${course.ExamTime}**\n${course.Room}\n`
            )
        );
      }),
        message.channel.send({
          content: `${message.author} **${courses[0].CourseTitle}\n ${title}`,
          embeds: listEm,
        });
    }
  });

  bot.on("interactionCreate", async (interaction: Interaction) => {
    if (!interaction.isCommand()) return;
    const { commandName } = interaction;

    if (interaction.isChatInputCommand()) {
      if (commandName === "ping") {
        await interaction.reply("Hey there! I'm alive! :D").catch((err) => {
          console.log("🚀 ~ file: index.ts:99 ~ bot.on ~ err:", err);
          interaction.followUp("Unknown command").catch((err) => {
            console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
          });
        });
      } else if (commandName === "uiu") {
        await interaction
          .reply(
            "**UIU** Discord Bot created by  <@669529872644833290>\nThe bot is still in development. Please report any bugs to the developer. [Contact](https://monzim.com/monzim)\nThanks for using the bot!"
          )
          .catch((err) => {
            console.log("🚀 ~ file: index.ts:105 ~ bot.on ~ err:", err);
            interaction.followUp("Unknown command").catch((err) => {
              console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
            });
          });
      } else if (commandName === "exam_time") {
        const courseCode: string =
          interaction.options.getString("course") ?? "";
        const section: string = interaction.options.getString("section") ?? "";

        await interaction.deferReply().catch((err) => {
          console.log("🚀 ~ file: index.ts:109 ~ bot.on ~ err:", err);
          interaction.followUp("Unknown command").catch((err) => {
            console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
          });
        });

        const courses = await findCoursesByCodeAndSection(
          courseCode.toUpperCase(),
          section.toUpperCase()
        );

        if (courses.length === 0) {
          console.log("No courses found");
          await interaction
            .followUp({
              content: `${interaction.user}  No courses found for ${courseCode} ${section}`,
            })
            .catch((err) => {
              console.log("🚀 ~ file: index.ts:124 ~ bot.on ~ err:", err);
              interaction.followUp("Unknown command").catch((err) => {
                console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
              });
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
        }),
          await interaction
            .followUp({
              content: `${interaction.user} **${courses[0].CourseTitle}**\n ${title}`,
              embeds: listEm,
            })
            .catch((err) => {
              console.log("🚀 ~ file: index.ts:154 ~ bot.on ~ err:", err);
              interaction.followUp("Unknown command").catch((err) => {
                console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
              });
            });
      } else {
        await interaction.followUp("Unknown command").catch((err) => {
          console.log("🚀 ~ file: index.ts:160 ~ bot.on ~ err:", err);
          interaction.followUp("Unknown command").catch((err) => {
            console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
          });
        });
      }
    } else {
      await interaction.reply("Unknown command LOL").catch((err) => {
        console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);

        interaction.followUp("Unknown command").catch((err) => {
          console.log("🚀 ~ file: index.ts:165 ~ bot.on ~ err:", err);
        });
      });
    }
  });
})();
