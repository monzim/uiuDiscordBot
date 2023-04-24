import {
  Client,
  ClientOptions,
  GatewayIntentBits,
  Message,
  ButtonStyle,
  EmbedBuilder,
  ButtonBuilder,
} from "discord.js";
import dotenv from "dotenv";

import { findCoursesByCodeAndSection } from "./cosmos_search";

dotenv.config();

const botOptions: ClientOptions = {
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.DirectMessages,
    GatewayIntentBits.MessageContent,
  ],
};

const bot = new Client(botOptions);
bot.login(process.env.BOT_TOKEN);

const prefix = "!"; // Command prefix

bot.on("ready", () => {
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

    let title = "The following **" + courses.length + "** courses were found\n";
    let listEm: any = [];
    // let colors: ColorResolvable = ["Random"];

    courses.map((course) => {
      listEm.push(
        new EmbedBuilder()
          .setTitle(course.CourseCode)
          .setColor("Random") // Set the color of the embed
          .setDescription(
            `Section: ${course.Section}     Faculty: ${course.Teacher}\n**${course.ExamDate} at ${course.ExamTime}**\n${course.Room}\n`
          )
      );

      //   embed.addFields({
      //     name: course.CourseCode,
      //     value: `\n**Section: ${course.Section}**     Faculty: ${course.Teacher}\n**${course.ExamDate} at ${course.ExamTime}**\n${course.Room}\n`,
      //   });
    }),
      message.channel.send({
        content: `${message.author} **${courses[0].CourseTitle}\n ${title}`,
        embeds: listEm,
      });
  }
});
