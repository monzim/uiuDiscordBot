import { ButtonBuilder, EmbedBuilder } from "discord.js";

export default interface Course {
  Dept: string;
  CourseCode: string;
  CourseTitle: string;
  Section: string;
  Teacher: string;
  ExamDate: string;
  ExamTime: string;
  Room: string;
  id: string;
  _rid: string;
  _self: string;
  _etag: string;
  _attachments: string;
  _ts: number;
}

export const generateMarkdownMessage = (courses: Course[]): string => {
  const total = courses.length;
  let message = "";
  message += `**${courses[0].CourseTitle}**\n`;
  message += total > 1 ? `**Total ${total} courses found**\n\n` : ``;

  courses.forEach((course, index) => {
    // message += `**Course ${index + 1}**\n`;
    // message += `Dept.: ${course.Dept}\n`; // Accessing property with dot notation
    // message += "```dart\n";

    // message += `${course.CourseCode}\n`;
    // message += `${course.CourseTitle}\n`;
    // message += `Section: ${course.Section}    `;
    // message += `Teacher: ${course.Teacher}\n`;
    // message += `${course.ExamDate} at ${course.ExamTime}\n`;
    // message += `${course.Room}\n`;

    // message += "```ts\n";
    message += `*${course.CourseCode}*\n`;
    // message += `**${course.CourseTitle}**\n`;
    message += `Section: ${course.Section}    `;
    message += `Teacher: ${course.Teacher}\n`;
    message += `**${course.ExamDate} at ${course.ExamTime}**\n`;
    message += `${course.Room}\n`;
    // message += `id: ${course.id}\n`;
    // message += `_rid: ${course._rid}\n`;
    // message += `_self: ${course._self}\n`;
    // message += `_etag: ${course._etag}\n`;
    // message += `_attachments: ${course._attachments}\n`;
    // message += `_ts: ${course._ts}\n`;
    // message += "\n";
    // message += "```"; // End of the code block
    message += "\n";
  });
  return message;
};
