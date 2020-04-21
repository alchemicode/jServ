import 'dart:async';
import 'dart:io';
import 'dart:convert';
import 'DataObject.dart';
import 'Collection.dart';
import 'AttributeContainer.dart';
import 'package:path/path.dart';

Map<String, dynamic> requestTypes = {
  "GET": true,
  "POST": true,
  "PUT": true,
  "HEAD": false,
  "DELETE": false,
  "PATCH": false,
  "OPTIONS": false
};
String ip = "localhost";
InternetAddress ipAddress;
int port = 4040;
List<Collection> dbs = new List<Collection>();

Future main() async {
  await readDatabases();
  await readConfig();
  await createIP();

  var server = await HttpServer.bind(
    ipAddress,
    port,
  );
  print("Listening on $ip:${server.port}");

  await for (HttpRequest request in server) {
    String m = request.method;
    if (requestTypes[m]) {
      handleRequest(request);
    } else {
      request.response.write("Error: Invalid request");
    }
  }
}

bool parseBool(String b) {
  if (b.toLowerCase() == "true") {
    return true;
  } else {
    return false;
  }
}

void handleRequest(HttpRequest r) {
  String m = r.method;
  try {
    switch (m) {
      case ("GET"):
        {
          handleGet(r);
        }
        break;
      case ("POST"):
        {
          handlePost(r);
        }
        break;
      case ("PUT"):
        {
          handlePost(r);
        }
        break;
      case ("HEAD"):
        {
          handlePost(r);
        }
        break;
      case ("DELETE"):
        {
          handlePost(r);
        }
        break;
      case ("PATCH"):
        {
          handlePost(r);
        }
        break;
      case ("OPTIONS"):
        {
          handlePost(r);
        }
        break;
    }
  } catch (e) {
    print("Exception in handleRequest: $e");
  }
}

void handleGet(HttpRequest r) {
  String end = "";
  String path = r.uri.path;

  if (path == "/query") {
    print("Object query from ${r.connectionInfo.remoteAddress}");
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    if (dbs.any((Collection value) => value.name == query)) {
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data = c.dataList.singleWhere((d) => d.id == id);
      end = data.toString();
    }
  }

  if (path == "/query/attribute") {
    print("Attribute query from ${r.connectionInfo.remoteAddress}");
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    String att = r.uri.queryParameters["a"];
    if (dbs.any((Collection value) => value.name == query)) {
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data = c.dataList.singleWhere((d) => d.id == id);
      AttributeContainer attribute =
          new AttributeContainer(att, data.data[att]);
      end = attribute.toJSON();
    }
  }

  var response = r.response;
  response.write(end);
  response.close();
}

void handlePost(HttpRequest r) {
  String end = "";
  String path = r.uri.path;

  if (path == "/add") {
    print("Add request from ${r.connectionInfo.remoteAddress}");
    String add = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    if (dbs.any((Collection value) => value.name == add)) {
      Collection c = dbs.singleWhere((col) => col.name == add);
      c.dataList.add(new DataObject.emptyMap(id));
      c.updateFile();
      var response = r.response;
      response.write("Successfully added object $id to $add");
      response.close();
    }
  } else if (path == "/add/obj") {
    print("Object add request from ${r.connectionInfo.remoteAddress}");
    String add = r.uri.queryParameters["q"];
    Future<String> content = utf8.decodeStream(r);
    if (dbs.any((Collection value) => value.name == add)) {
      Collection c = dbs.singleWhere((col) => col.name == add);
      content.then((result) {
        DataObject d = new DataObject.fromJsonString(result);
        c.dataList.add(d);
        c.updateFile();
        var response = r.response;
        response.write("\nSuccessfully added this object to $add");
        response.close();
      });
    }
  } else if (path == "/add/attribute") {
    String newEnd = "";
    print("Attribute add request from ${r.connectionInfo.remoteAddress}");
    String add = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    String att = r.uri.queryParameters["a"];
    Future<String> content = utf8.decodeStream(r);
    if (dbs.any((Collection value) => value.name == add)) {
      Collection c = dbs.singleWhere((col) => col.name == add);
      content.then((result) {
        AttributeContainer attribute =
            new AttributeContainer(att, json.decode(result)[att]);
        DataObject data = c.dataList.singleWhere((d) => d.id == id);
        print(attribute.key);
        if (!data.data.containsKey(attribute.key)) {
          data.data[att] = attribute.value;
          c.updateFile();
          var response = r.response;
          response.write("Successfully added $att to $id");
          response.close();
        } else {
          var response = r.response;
          response.write("The attribute $att already exists in object $id");
          response.close();
        }
      });
    }
  } else if (path == "/mod") {
    print("Object mod request from ${r.connectionInfo.remoteAddress}");
    String mod = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    int newId = int.parse(r.uri.queryParameters["v"]);
    if (dbs.any((Collection value) => value.name == mod)) {
      Collection c = dbs.singleWhere((col) => col.name == mod);
      DataObject data = c.dataList.singleWhere((d) => d.id == id);
      data.id = newId;
      c.updateFile();
      var response = r.response;
      response.write("Successfully modified $id to $newId");
      response.close();
    }
  } else if (path == "/mod/attribute") {
    String newEnd = "";
    print("Attribute mod request from ${r.connectionInfo.remoteAddress}");
    String mod = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    String att = r.uri.queryParameters["a"];
    Future<String> content = utf8.decodeStream(r);
    if (dbs.any((Collection value) => value.name == mod)) {
      Collection c = dbs.singleWhere((col) => col.name == mod);
      content.then((result) {
        AttributeContainer attribute =
            new AttributeContainer(att, json.decode(result)[att]);

        DataObject data = c.dataList.singleWhere((d) => d.id == id);

        if (data.data.containsKey(att)) {
          data.data[att] = attribute.value;
          c.updateFile();
          var response = r.response;
          response.write("Successfully modified $att of $id in $mod");
          response.close();
        } else {
          var response = r.response;
          response.write("The attribute $att does not exist in object $id");
          response.close();
        }
      });
    }
    end = newEnd;
  } else {
    var response = r.response;
    response.write("Error: Invalid request");
    response.close();
  }
}

void handlePut(HttpRequest r) {}
void handleHead(HttpRequest r) {}
void handleDelete(HttpRequest r) {
  String m = r.method;
  r.response.write("this is a $m");
}

void handlePatch(HttpRequest r) {}
void handleOptions(HttpRequest r) {
  String m = r.method;
  r.response.write("this is a $m");
}

Future readConfig() async {
  File file = new File("config.json");
  String content;

  if (await file.exists()) {
    content = await file.readAsString();
  }
  var map = json.decode(content);
  ip = map["ip"];
  port = map["port"];
  requestTypes = map["Requests"];
}

void readDatabases() {
  Directory dir = new Directory("Databases");
  dir.list(recursive: false).listen((FileSystemEntity e) {
    print("file: ");
    print(e.path);
    String name = basename(e.path).split(".")[0];
    print(name);
    Collection c = new Collection(name);
    dbs.add(c);
  });
}

void createIP() {
  if (ip == "localhost" || ip == "127.0.0.1") {
    ipAddress = InternetAddress.loopbackIPv4;
  } else {
    ipAddress = new InternetAddress(ip);
  }
}
