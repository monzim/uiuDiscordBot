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

export const generateMarkdownMessageWithCourse = (
  courses: Course[]
): string => {
  const total = courses.length;
  let message = "";
  message += total > 1 ? `**Total ${total} courses found**\n\n` : ``;

  courses.forEach((course, index) => {
    message += `*${course.CourseCode}*\n`;
    message += `Section: ${course.Section}    `;
    message += `Teacher: ${course.Teacher}\n`;
    message += `**${course.ExamDate} at ${course.ExamTime}**\n`;
    message += `${course.Room}\n`;

    message += "\n";
  });
  return message;
};
