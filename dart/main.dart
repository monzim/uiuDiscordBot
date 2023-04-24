import 'dart:convert';
import 'dart:io';

void main() {
  String filePath = "input.json"; // Update with your file path
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
  int count = 0;
  for (var obj in data) {
    count++;

    // if (count == 6) {
    //   break;
    // }
    // Extract the room number from the "Room" field
    String room = obj["Room"];

    print("room: $room");

    RegExp regex = RegExp(r'\((\d+-\d+)\)');
    Iterable<Match>? matches = regex.allMatches(room);

    List<String> rollRanges = [];
    List<int> rolls = [];
    for (Match match in matches) {
      // Extract the roll number range from the match
      String? rollRange = match.group(1);
      if (rollRange == null) {
        print("rollRange is null");
        continue;
      }

      int startRoll = int.parse(rollRange.split("-")[0]);
      int endRoll = int.parse(rollRange.split("-")[1]);

      // Create a list of roll numbers from start to end, inclusive
      List<int> rangeRolls =
          List.generate(endRoll - startRoll + 1, (index) => startRoll + index);

      rolls.addAll(rangeRolls);

      print("Total Roles: ${rangeRolls.length} | rolls: ${rolls.length}");

      // Add the roll number range to the list
      rollRanges.add(rollRange);

      // Update the "Room" field with the roll number range
      obj["Room"] = obj["Room"].replaceFirst(regex, rollRange);
    }
  }

  // Write the updated JSON data back to file
  file.writeAsStringSync(jsonEncode(data));
}
