import {
  CacheType,
  ChatInputCommandInteraction,
  EmbedBuilder,
} from "discord.js";
import { sendWebhookErrorMessage } from "../webhook/send_message";
import { findCoursesByDepartmentCourseCodeAndSection } from "../cosmos/cosmos_search";

// export const examTimeHandlr = async (
//   interaction: ChatInputCommandInteraction<CacheType>
// ) => {
//   await interaction
//     .reply({
//       content: `Opps! Currently, no exam is scheduled. It will be updated as soon as the exam schedule is published.`,
//     })
//     .catch((err) => {
//       console.log("🚀 ~ file: examtimehandlr.ts:11 ~ err:", err);
//     });
// };

export const examTimeHandlrProd = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  const department: string = interaction.options.getString("department") ?? "";
  const courseCode: string = interaction.options.getString("course") ?? "";
  const section: string = interaction.options.getString("section") ?? "";

  await interaction.deferReply().catch((err) => {
    sendWebhookErrorMessage("index.ts:199", err);

    interaction.followUp("Unknown command").catch((err) => {
      sendWebhookErrorMessage("index.ts:204", err);
    });
  });

  const courses = await findCoursesByDepartmentCourseCodeAndSection(
    department.toUpperCase(),
    courseCode.toUpperCase(),
    section.toUpperCase()
  );

  if (courses.length === 0) {
    console.log("No courses found");
    await interaction
      .followUp({
        content: `### ${
          interaction.user
        } **NO COURSE FOUND** for ${courseCode.toUpperCase()} ${section.toUpperCase()}\n### **Support the bot's longevity with a donation if you found it helpful. [Click here for 🎁 Donate](https://monzim.com/support). Thank you!**  `,
      })
      .catch((err) => {
        sendWebhookErrorMessage("index.ts:223", err);

        interaction.followUp("Unknown command").catch((err) => {
          sendWebhookErrorMessage("index.ts:228", err);
        });
      });
    return;
  }

  let listEm: EmbedBuilder[] = [];

  courses.map((course) => {
    let _selectedSection = course.Section.toLowerCase();
    _selectedSection = _selectedSection.replace(/\s/g, "");

    if (_selectedSection === section.toLowerCase()) {
      listEm.push(
        new EmbedBuilder()
          .setTitle(course.CourseCode)
          .setColor("Random")
          .setDescription(
            `Section: ${course.Section}     Faculty: ${course.Teacher}\n**${course.ExamDate} at ${course.ExamTime}**\n${course.Room}\n`
          )
      );
    }
  });

  // set last embed footer
  listEm[listEm.length - 1]
    // .setThumbnail("https://source.unsplash.com/user/c_v_r/300x00")
    .setFooter({
      text: "Help Us Make a Difference",
      iconURL:
        "https://res.cloudinary.com/monzim/image/upload/v1688984685/download_kh1syl.png",
    });

  const courseCount = courses.length;
  const matchCount = listEm.length;
  const isMultipleCourses = courseCount > 1;
  const isMultipleMatches = matchCount > 1;

  const coursePlural = isMultipleCourses ? "courses" : "course";
  const matchPlural = isMultipleMatches ? "courses" : "course";

  const title = `In total, the query has **${courseCount}** ${coursePlural}. The following **${matchCount}** ${matchPlural} match. :)`;

  try {
    // await interaction.followUp({
    //   content: `## ${interaction.user} **${courses[0].CourseTitle}**\n### **Support the bot's longevity with a donation if you found it helpful. [Click here for 🎁 Donate](https://monzim.com/support). Thank you!** `,
    //   embeds: listEm,
    // });

    await interaction.followUp({
      content: `### ${interaction.user} **Support the bot's longevity with a donation if you found it helpful. [Click here for 🎁 Donate](https://monzim.com/support). Thank you!** `,
      embeds: listEm,
    });
  } catch (err) {
    sendWebhookErrorMessage("index.ts:261", err);

    try {
      await interaction.followUp("Unknown command");
    } catch (err) {
      sendWebhookErrorMessage("index.ts:267", err);
    }
  }

  // // send a message for donation link
  // await interaction.user
  //   .send({
  //     content: `Support the bot's longevity with a donation if you found it helpful. Thank you!`,
  //     embeds: [
  //       new EmbedBuilder()
  //         .setTitle("Donate")
  //         .setColor("Random")
  //         .setDescription(
  //           `Support the bot's longevity with a donation if you found it helpful. Thank you!`
  //         )
  //         .addFields({
  //           name: "PayPal",
  //           value: `[Donate](https://www.paypal.com/donate?hosted_button_id=ZQZQZQZQZQZQZ)`,
  //         }),
  //     ],
  //   })
  //   .catch((err) => {
  //     console.log("🚀 ~ file: examtimehandlr.ts:111 ~ err:", err);
  //   });
  // .followUp({
  //   content: `Support the bot's longevity with a donation if you found it helpful. Thank you!`,
  //   embeds: [
  //     new EmbedBuilder()
  //       .setTitle("Donate")
  //       .setColor("Random")
  //       .setDescription(
  //         `Support the bot's longevity with a donation if you found it helpful. Thank you!`
  //       )
  //       .addFields({
  //         name: "PayPal",
  //         value: `[Donate](https://www.paypal.com/donate?hosted_button_id=ZQZQZQZQZQZQZ)`,
  //       }),
  //   ],
  // })
  // .catch((err) => {
  //   sendWebhookErrorMessage("index.ts:297", err);

  //   interaction.followUp("Unknown command").catch((err) => {
  //     sendWebhookErrorMessage("index.ts:302", err);
  //   });
  // });
};
