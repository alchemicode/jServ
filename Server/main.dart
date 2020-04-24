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
  "PUT": false,
  "HEAD": false,
  "DELETE": true,
  "PATCH": false,
  "OPTIONS": false
};
String ip = "localhost";
InternetAddress ipAddress;
int port = 4040;
List<Collection> dbs = new List<Collection>();

Future main() async {
  await startSequence();
  var server = await HttpServer.bind(
    ipAddress,
    port,
  );
  print(" * Server bound to $ip:${server.port}");

  await for (HttpRequest request in server) {
    String m = request.method;
    if (requestTypes[m]) {
      print("\n");
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
          handleDelete(r);
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
  String path = r.uri.path;

  if (path == "/query") {
    print("Object query from ${r.connectionInfo.remoteAddress}");
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    print("Queried $query for object $id");
    if (dbs.any((Collection value) => value.name == query)) {
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data =
          c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
      if (data != null) {
        String end = data.toString();
        var response = r.response;
        response.write(end);
        response.close();
        print(end);
      } else {
        String end = "Object $id could not be found in $query";
        var response = r.response;
        response.write(end);
        response.close();
        print(end);
      }
    } else {
      String end = "Could not find collection $query";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/query/attribute") {
    print("Attribute query from ${r.connectionInfo.remoteAddress}");
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    String att = r.uri.queryParameters["a"];
    print("Queried $query for attribute $att of object $id");
    if (dbs.any((Collection value) => value.name == query)) {
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data = c.dataList.singleWhere((d) => d.id == id);
      if (data.data.containsKey(att)) {
        AttributeContainer attribute =
            new AttributeContainer(att, data.data[att]);
        String end = attribute.toJSON();
        var response = r.response;
        response.write(end);
        response.close();
        print(end);
      } else {
        String end =
            "Attribute $att could not be found in object $id in $query";
        var response = r.response;
        response.write(end);
        response.close();
        print(end);
      }
    } else {
      String end = "Could not find collection $query";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/query/allAttributes") {
    print("All attributes query from ${r.connectionInfo.remoteAddress}");
    String query = r.uri.queryParameters["q"];
    String att = r.uri.queryParameters["a"];
    print("Queried $query for attributes $att");
    if (dbs.any((Collection value) => value.name == query)) {
      Collection c = dbs.singleWhere((col) => col.name == query);
      List<Map<String, dynamic>> atts = new List<Map<String, dynamic>>();
      c.dataList.forEach((DataObject d) {
        if (d.data.containsKey(att)) {
          AttributeContainer atc = new AttributeContainer(att, d.data[att]);
          atts.add({d.id.toString(): atc.value});
        }
      });
      String end = json.encode(atts);
      var response = r.response;
      response.write(end);
      response.close();
      print(end);
    } else {
      String end = "Could not find collection $query";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/query/newId") {
    print("New ID query from ${r.connectionInfo.remoteAddress}");
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    print("Queried $query for new id $id");
    if (dbs.any((Collection value) => value.name == query)) {
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data =
          c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
      if (data == null) {
        String end = json.encode(true);
        var response = r.response;
        response.write(end);
        response.close();
        print(end);
      } else {
        String end = "Object $id is not available in $query";
        var response = r.response;
        response.write(end);
        response.close();
        print(end);
      }
    } else {
      String end = "Could not find collection $query";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else {
    print("Error: Invalid GET request from ${r.connectionInfo.remoteAddress}");
    print("\'" + r.uri.path + "\'");
    var response = r.response;
    response.write("Error: Invalid request");
    response.close();
  }
}

void handlePost(HttpRequest r) {
  String path = r.uri.path;

  if (path == "/add") {
    print("Empty object add request from ${r.connectionInfo.remoteAddress}");
    String add = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    print("Requested to add object $id to $add");
    if (dbs.any((Collection value) => value.name == add)) {
      Collection c = dbs.singleWhere((col) => col.name == add);
      c.dataList.add(new DataObject.emptyMap(id));
      c.updateFile();
      String end = "Successfully added object $id to $add";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    } else {
      String end = "Could not find collection $add";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/add/object") {
    print("Object add request from ${r.connectionInfo.remoteAddress}");
    String add = r.uri.queryParameters["q"];
    Future<String> content = utf8.decodeStream(r);
    if (dbs.any((Collection value) => value.name == add)) {
      Collection c = dbs.singleWhere((col) => col.name == add);
      DataObject d;
      content.then((result) {
        try {
          d = new DataObject.fromJsonString(result);
          print("Requested to add object ${d.id} to $add");
          c.dataList.add(d);
          c.updateFile();
          String end = "\nSuccessfully added this object to $add";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        } catch (e) {
          String end = "Could not parse DataObject from request body\n" + e;
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      });
    } else {
      String end = "Could not find collection $add";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/add/attribute") {
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
        DataObject data =
            c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
        print("Requested to add attribute ${attribute.key} to $id in $add");
        if (data != null) {
          if (!data.data.containsKey(attribute.key)) {
            data.data[att] = attribute.value;
            c.updateFile();
            String end = "Successfully added $att to $id";
            print(end);
            var response = r.response;
            response.write(end);
            response.close();
          } else {
            String end = "The attribute $att already exists in object $id";
            print(end);
            var response = r.response;
            response.write(end);
            response.close();
          }
        } else {
          String end = "Could not find object $id in $add";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      });
    } else {
      String end = "Could not find collection $add";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/mod/object") {
    print("Object mod request from ${r.connectionInfo.remoteAddress}");
    String mod = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    int newId = int.parse(r.uri.queryParameters["v"]);
    print("Requested to change $id to $newId in $mod");
    if (dbs.any((Collection value) => value.name == mod)) {
      Collection c = dbs.singleWhere((col) => col.name == mod);
      DataObject data =
          c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
      if (data != null) {
        data.id = newId;
        c.updateFile();
        String end = "Successfully modified $id to $newId";
        print(end);
        var response = r.response;
        response.write(end);
        response.close();
      } else {
        String end = "Could not find Object $id in $mod";
        print(end);
        var response = r.response;
        response.write(end);
        response.close();
      }
    } else {
      String end = "Could not find collection $mod";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/mod/attribute") {
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
        print("Requested to modify attribute ${attribute.key} of $id in $mod");
        DataObject data =
            c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
        if (data != null) {
          if (data.data.containsKey(att)) {
            data.data[att] = attribute.value;
            c.updateFile();
            String end = "Successfully modified $att of $id in $mod";
            print(end);
            var response = r.response;
            response.write(end);
            response.close();
          } else {
            String end = "Could not find attribute $att in object $id";
            print(end);
            var response = r.response;
            response.write(end);
            response.close();
          }
        } else {
          String end = "Could not find object $id";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      });
    } else {
      String end = "Could not find collection $mod";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else {
    print("Error: Invalid POST request from ${r.connectionInfo.remoteAddress}");
    print("\'" + r.uri.path + "\'");
    var response = r.response;
    response.write("Error: Invalid request");
    response.close();
  }
}

void handlePut(HttpRequest r) {}
void handleHead(HttpRequest r) {}
void handleDelete(HttpRequest r) {
  String path = r.uri.path;
  if (path == "/delete/object") {
    print("Object delete request from ${r.connectionInfo.remoteAddress}");
    String del = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    if (dbs.any((Collection value) => value.name == del)) {
      Collection c = dbs.singleWhere((col) => col.name == del);
      print("Requested to delete object $id in $del");
      DataObject data =
          c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
      if (data != null) {
        c.dataList.remove(data);
        c.updateFile();
        String end = "Successfully deleted object $id from $del";
        print(end);
        var response = r.response;
        response.write(end);
        response.close();
      } else {
        String end = "Could not find object $id in $del";
        print(end);
        var response = r.response;
        response.write(end);
        response.close();
      }
    } else {
      String end = "Could not find collection $del";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else if (path == "/delete/attribute") {
    print("Attribute delete request from ${r.connectionInfo.remoteAddress}");
    String del = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    String att = r.uri.queryParameters["a"];
    if (dbs.any((Collection value) => value.name == del)) {
      Collection c = dbs.singleWhere((col) => col.name == del);
      print("Requested to delete object $id in $del");
      DataObject data =
          c.dataList.singleWhere((d) => d.id == id, orElse: () => null);
      if (data != null) {
        if (data.data.containsKey(att)) {
          data.data.remove(att);
          c.updateFile();
          String end = "Successfully deleted attribute $att from $id";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        } else {
          String end = "Could not find attribute $att in $id";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      } else {
        String end = "Could not find object $id in $del";
        print(end);
        var response = r.response;
        response.write(end);
        response.close();
      }
    } else {
      String end = "Could not find collection $del";
      print(end);
      var response = r.response;
      response.write(end);
      response.close();
    }
  } else {
    print(
        "Error: Invalid DELETE request from ${r.connectionInfo.remoteAddress}");
    print("\'" + r.uri.path + "\'");
    var response = r.response;
    response.write("Error: Invalid request");
    response.close();
  }
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
    print("--Loaded database: ${basename(e.path)}");
    String name = basename(e.path).split(".")[0];
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

Future<void> startSequence() async {
  String version = await File("version.jserv").readAsString();
  print(" * Starting jServ v$version");
  print(" -----------------------");
  await readDatabases();
  print(" * Loading databases...");
  await readConfig();
  print(" * Loading config...");
  await createIP();
  print(" * Binding server...");
  print(" * Done!");
}
