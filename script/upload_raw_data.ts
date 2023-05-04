import { CosmosClient } from "@azure/cosmos";
import * as fs from "fs";

import {
  COSMOS_CONTAINER_ID,
  COSMOS_DATABASE_ID,
  COSMOS_KEY,
  COSMOS_ENDPOINT,
} from "../src/config";

const filesPath: string[] = [
  "data/spring_23_final/bscse_bseee.json",
  "data/spring_23_final/bba.json",
];
const client = new CosmosClient({
  endpoint: COSMOS_ENDPOINT as string,
  key: COSMOS_KEY as string,
});

const insertFileData = async (filePath: string) => {
  // Read the JSON data from file
  const jsonData = JSON.parse(fs.readFileSync(filePath, "utf8"));

  let count = 0;

  try {
    const { database } = await client.databases.createIfNotExists({
      id: COSMOS_DATABASE_ID,
    });
    const { container } = await database.containers.createIfNotExists({
      id: COSMOS_CONTAINER_ID,
    });

    for (const data of jsonData) {
      count++;

      await container.items
        .create(data)
        .then((result) => {
          console.log(
            `Count: ${count} Inserted item with id: ${result.resource.id}`
          );
        })
        .catch((err) => {
          console.log(`Count: ${count} Error inserting data: ${err}`);
        });
    }

    console.log(`Total data inserted: ${count}`);
  } catch (error) {
    console.log("🚀 ~ insertReadDataData ~ error:", error);
  }
};

async function main() {
  for (const filePath of filesPath) {
    console.log(`Inserting data from ${filePath}`);
    await insertFileData(filePath);
  }
}

main();
