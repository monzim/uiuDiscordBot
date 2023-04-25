import {
  Client,
  ClientOptions,
  GatewayIntentBits,
  Message,
  EmbedBuilder,
  Interaction,
} from "discord.js";

import { findCoursesByCodeAndSection } from "./cosmos_search";
import { discordToken, clientID } from "./config";
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
        await interaction.reply("Pong!");
      } else if (commandName === "exam") {
        const courseCode: string =
          interaction.options.getString("course") ?? "";
        const section: string = interaction.options.getString("section") ?? "";
        console.log("Course code: " + courseCode);
        console.log("Section: " + section);

        // await interaction.reply("Ole Ole " + courseCode + " " + section);

        await interaction.deferReply();

        const courses = await findCoursesByCodeAndSection(
          courseCode.toUpperCase(),
          section.toUpperCase()
        );
        console.log("🚀 ~ file: index.ts:153 ~ bot.on ~ courses:", courses);

        if (courses.length === 0) {
          await interaction.followUp({
            content: `${interaction.user}  No courses found for ${courseCode} ${section}`,
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
          await interaction.followUp({
            content: `${interaction.user} **${courses[0].CourseTitle}\n ${title}`,
            embeds: listEm,
          });
        // await interaction.editReply("Pong!");
      } else {
        await interaction.followUp("Unknown command");
      }
    } else {
      await interaction.reply("Unknown command LOL");
    }
  });
})();

// const replayWithCourse = async (
//   message: string,
//   interaction: ChatInputCommandInteraction<CacheType>
// ) => {
//   try {
//     await interaction.deferReply();
//     await interaction.followUp(message);
//   } catch (error) {
//     console.error("Failed to acknowledge interaction", error);
//     try {
//       await interaction.followUp("Failed to acknowledge interaction");
//     } catch (error) {
//       console.error("Failed to send message for interaction", error);
//     }
//   }
// };

// bot.on("interactionCreate", async (interaction: Interaction<CacheType>) => {
//   if (!interaction.isCommand()) return;
//   const { commandName } = interaction;

//   if (interaction.isChatInputCommand()) {
//     if (commandName === "ping") {
//       await interaction.reply("Pong!");
//     } else if (commandName === "exam") {
//       // const courseCode = interaction.options.getString("course");
//       // const section = interaction.options.getString("section");
//       // await interaction.reply("Exam command " + courseCode + " " + section);
//     } else {
//       await interaction.reply("Unknown command");
//     }
//   } else if (interaction.isAutocomplete()) {
//     const autoCompleteInt = interaction as AutocompleteInteraction;
//     const focusedOption = autoCompleteInt.options.getFocused(true);
//     let choices: string[] = [
//       "Popular Topics: Threads",
//       "Sharding: Getting started",
//       "Library: Voice Connections",
//       "Interactions: Replying to slash commands",
//       "Popular Topics: Embed preview",
//     ];

//     // if (focusedOption.name === "version") {
//     //   choices = ["v9", "v11", "v12", "v13", "v14"];
//     // }
//     const filtered = choices.filter((choice) =>
//       choice.startsWith(focusedOption.value)
//     );

//     await autoCompleteInt.respond(
//       filtered.map((choice) => ({ name: choice, value: choice }))
//     );
//   }
// });
