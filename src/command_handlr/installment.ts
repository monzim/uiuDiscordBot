import {
  CacheType,
  ChatInputCommandInteraction,
  ColorResolvable,
  EmbedBuilder,
} from "discord.js";
import { sendWebhookErrorMessage } from "../webhook/send_message";

type Installment = {
  index: number;
  date: Date;
  message: string;
  color: ColorResolvable;
};

export const installmentHandlr = async (
  interaction: ChatInputCommandInteraction<CacheType>
) => {
  const installmentList: Installment[] = [
    {
      index: 1,
      color: "DarkBlue",
      date: new Date("Jun 18, 2023"),
      message:
        "1st installment: A fine of Tk. 1,000/- will be imposed if 40% of Tuition Fee and Trimester Fee is not paid within this date.\nTransportation fee, if applicable, must be paid in full (no installment)",
    },
    {
      index: 2,
      color: "Green",
      date: new Date("Jul 16, 2023"),
      message:
        "2nd installment: A fine of Tk. 1,000/- will be imposed if 70% of Tuition Fee and Trimester Fee is not paid within this date",
    },
    {
      index: 3,
      color: "Red",
      date: new Date("Aug 13, 2023"),
      message:
        "3rd installment: A fine of Tk. 1,000/- will be imposed if 100% of Tuition Fee and Trimester Fee is not paid within this date.",
    },
  ];

  const now = new Date();

  let daysLeft = 0;
  let nextPayment = installmentList[0];
  if (now < installmentList[0].date) {
    daysLeft = Math.ceil(
      (installmentList[0].date.getTime() - now.getTime()) /
        (1000 * 60 * 60 * 24)
    );
  } else if (now < installmentList[1].date) {
    daysLeft = Math.ceil(
      (installmentList[1].date.getTime() - now.getTime()) /
        (1000 * 60 * 60 * 24)
    );
    nextPayment = installmentList[1];
  } else if (now < installmentList[2].date) {
    daysLeft = Math.ceil(
      (installmentList[2].date.getTime() - now.getTime()) /
        (1000 * 60 * 60 * 24)
    );
    nextPayment = installmentList[2];
  }

  let listEm: any = [];

  installmentList.forEach((installment) => {
    listEm.push(
      new EmbedBuilder()
        .setTitle(
          `Payment ${
            installment.index
          } - **${installment.date.toDateString()}**`
        )
        .setColor(installment.color)
        .setDescription(`${installment.message}`)
    );
  });

  await interaction
    .reply({
      content: `Installment Payment Plan. Please note the following payment due dates and payment requirements. **Next payment is due in ${daysLeft} days on ${nextPayment.date.toDateString()}**\n`,
      embeds: listEm,
    })
    .catch((err) => {
      console.log("🚀 ~ file: installment.ts:87 ~ err:", err);
      sendWebhookErrorMessage("installmentHandlr", err);
      interaction.followUp("Error occurred while processing the command");
    });
};
