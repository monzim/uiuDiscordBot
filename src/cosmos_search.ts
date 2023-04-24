import { CosmosClient } from "@azure/cosmos";
import * as config from "./config.json"; // Import the config file
import Course from "./models/course";

// Read the Cosmos DB account settings from the config file
const { endpoint, key, databaseId, containerId } = config;

// Connect to Cosmos DB account
const cosmosClient = new CosmosClient({ endpoint, key });

export const findCoursesByCodeAndSection = async (
  courseCode: string,
  section: string
): Promise<Course[]> => {
  try {
    const { database } = await cosmosClient.databases.createIfNotExists({
      id: databaseId,
    });
    const { container } = await database.containers.createIfNotExists({
      id: containerId,
    });

    // Query the container for courses with matching roll
    const query = `SELECT * FROM c WHERE CONTAINS(c['Course Code'], '${courseCode}') AND CONTAINS(c['Section'], '${section}')`;

    const { resources: courses } = await container.items
      .query(query)
      .fetchAll();

    // return courses;

    // Map the query result to Course objects
    const mappedCourses: Course[] = courses.map((course: any) => {
      return {
        Dept: course["Dept."],
        CourseCode: course["Course Code"],
        CourseTitle: course["Course Title"],
        Section: course.Section,
        Teacher: course.Teacher,
        ExamDate: course["Exam Date"],
        ExamTime: course["Exam Time"],
        Room: course.Room,
        id: course.id,
        _rid: course._rid,
        _self: course._self,
        _etag: course._etag,
        _attachments: course._attachments,
        _ts: course._ts,
      };
    });

    console.log("Courses found:", mappedCourses);

    return mappedCourses;
  } catch (error) {
    console.error("Error finding courses by roll:", error);
    throw error;
  }
};

// export const findCoursesByRollAsync = async () => {
//   const rl = readline.createInterface({
//     input: process.stdin,
//     output: process.stdout,
//   });

//   // Prompt the user for course code and section
//   rl.question("Enter course code: ", async (code) => {
//     rl.question("Enter section: ", async (section) => {
//       const courses = await findCoursesByCodeAndSection(code, section);
//       console.log("Courses found:", courses);
//       rl.close();
//     });
//   });
// };
