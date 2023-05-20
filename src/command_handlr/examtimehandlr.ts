import { CacheType, ChatInputCommandInteraction } from "discord.js";

export const examTimeHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  await interaction
    .reply({
      content: `Opps! Currently, no exam is scheduled. It will be updated as soon as the exam schedule is published.`,
    })
    .catch((err) => {
      console.log("🚀 ~ file: examtimehandlr.ts:11 ~ err:", err);
    });
};

// const department: string =
//   interaction.options.getString("department") ?? "";

// const courseCode: string =
//   interaction.options.getString("course") ?? "";

// const section: string = interaction.options.getString("section") ?? "";

// await interaction.deferReply().catch((err) => {
//   console.log("🚀 ~ file: index.ts:199 ~ bot.on ~ err:", err);
//   sendWebhookErrorMessage("index.ts:199", err);

//   interaction.followUp("Unknown command").catch((err) => {
//     console.log("🚀 ~ file: index.ts:204 ~ bot.on ~ err:", err);
//     sendWebhookErrorMessage("index.ts:204", err);
//   });
// });

// const courses = await findCoursesByDepartmentCourseCodeAndSection(
//   department.toUpperCase(),
//   courseCode.toUpperCase(),
//   section.toUpperCase()
// );

// if (courses.length === 0) {
//   console.log("No courses found");
//   await interaction
//     .followUp({
//       content: `${interaction.user}  No courses found for ${courseCode} ${section}`,
//     })
//     .catch((err) => {
//       console.log("🚀 ~ file: index.ts:223 ~ bot.on ~ err:", err);
//       sendWebhookErrorMessage("index.ts:223", err);

//       interaction.followUp("Unknown command").catch((err) => {
//         console.log("🚀 ~ file: index.ts:228 ~ bot.on ~ err:", err);
//         sendWebhookErrorMessage("index.ts:228", err);
//       });
//     });
//   return;
// }

// let listEm: any = [];

// courses.map((course) => {
//   let _selectedSection = course.Section.toLowerCase();
//   _selectedSection = _selectedSection.replace(/\s/g, "");

//   if (_selectedSection === section.toLowerCase()) {
//     listEm.push(
//       new EmbedBuilder()
//         .setTitle(course.CourseCode)
//         .setColor("Random")
//         .setDescription(
//           `Section: ${course.Section}     Faculty: ${course.Teacher}\n**${course.ExamDate} at ${course.ExamTime}**\n${course.Room}\n`
//         )
//     );
//   }
// });

// let title = `In total, the query has **${courses.length}** ${
//   courses.length > 1 ? "courses" : "course"
// }. The following **${listEm.length}** ${
//   listEm.length > 1 ? "courses" : "course"
// } match. :)`;

// await interaction
//   .followUp({
//     content: `${interaction.user} **${courses[0].CourseTitle}**\n ${title}`,
//     embeds: listEm,
//   })
//   .catch((err) => {
//     console.log("🚀 ~ file: index.ts:261 ~ bot.on ~ err:", err);
//     sendWebhookErrorMessage("index.ts:261", err);

//     interaction.followUp("Unknown command").catch((err) => {
//       console.log("🚀 ~ file: index.ts:267 ~ bot.on ~ err:", err);
//       sendWebhookErrorMessage("index.ts:267", err);
//     });
//   });
