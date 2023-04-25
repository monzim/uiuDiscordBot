import 'dart:convert';
import 'dart:io';

void main() {
  String filePath = "courses.json"; // Update with your file path
  updateJsonFile(filePath);
  print("JSON file updated successfully.");
}

void updateJsonFile(String filePath) {
  // Load JSON data from file
  File file = File(filePath);
  String jsonString = file.readAsStringSync();
  List<dynamic> data = jsonDecode(jsonString);

  // Loop through each object in the JSON array
  print("data.length: ${data.length}");

  List<String> courses = [];

  for (var obj in data) {
    String room = obj["Course Code"];
    courses.add(room);
  }

  print("Total Courses: ${courses.length}");
  final raw = courses.toSet();
  print("Total Unique Courses: ${raw.length},");
  // print int this format {name: 'couse', value: ''}
  for (var course in raw) {
    print('{name: "${course}", value: "${course}"},');
  }

  // Write the updated JSON data back to file
  // file.writeAsStringSync(jsonEncode(data));
}
