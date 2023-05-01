import {
  CacheType,
  ChatInputCommandInteraction,
  ColorResolvable,
  EmbedBuilder,
} from "discord.js";
import { sendWebhookErrorMessage } from "../webhook/send_message";

type ExamType = {
  startDate: Date;
  endDate: Date;
  day: string;
  type: string;
  timeLeft?: string;
  color: ColorResolvable;
};

function timeUntilNextExam(exam: ExamType): string {
  const now = new Date();
  const timeLeft = exam.startDate.getTime() - now.getTime();
  const daysLeft = Math.floor(timeLeft / (1000 * 60 * 60 * 24));
  const hoursLeft = Math.floor((timeLeft / (1000 * 60 * 60)) % 24);
  const minutesLeft = Math.floor((timeLeft / 1000 / 60) % 60);
  const secondsLeft = Math.floor((timeLeft / 1000) % 60);

  return `${daysLeft} days ${hoursLeft} hours ${minutesLeft} minutes ${secondsLeft} seconds`;
}

export const examHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  // Jul 15 – 22, 2023 Sat - Sat Mid-Term Exam
  // Sep 3 -12, 2023 Sun - Tue Final Exam

  const examList: ExamType[] = [
    {
      startDate: new Date("Jul 15, 2023"),
      endDate: new Date("Jul 22, 2023"),
      day: "Sat - Sat",
      type: "Mid-Term Exam",
      color: "DarkOrange",
    },
    {
      startDate: new Date("Sep 3, 2023"),
      endDate: new Date("Sep 12, 2023"),
      day: "Sun - Tue",
      type: "Final Exam",
      color: "DarkRed",
    },
  ];

  let listEm: any = [];

  examList.forEach((exam) => {
    listEm.push(
      new EmbedBuilder()
        .setTitle(
          `${exam.type} on ${
            exam.day
          } from ${exam.startDate.toDateString()} to ${exam.endDate.toDateString()}`
        )
        .setColor(exam.color)
        .setFooter({
          text: `Time left: ${timeUntilNextExam(exam)}`,
        })
    );
  });

  await interaction
    .reply({
      content: `Upcoming bamboo 🎋\n`,
      embeds: listEm,
    })
    .catch((err) => {
      console.log("🚀 ~ file: examhandlr.ts:71 ~ err:", err);

      sendWebhookErrorMessage("examhandlr", err);
      interaction.followUp("Error occurred while executing the command!");
    });
};
