import 'dart:async';
import 'dart:io';
import 'dart:convert';
import 'dart:math';
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

Map<String, bool> requestPermissions = {
  //False denotes that user API Keys can request, true denotes that admin API Key is required
  "QObject": false, 
  "QAttribute": false, 
  "QAllAttributes": false, 
  "QByAttribute": false,
  "QnewId": true, 
  "AEmpty": false, 
  "AObject": false, 
  "AAttribute": false, 
  "MObject": false, 
  "MAttribute": false, 
  "DObject": false, 
  "DAttribute": false 
};
String ip = "localhost";
InternetAddress ipAddress;
int port = 4040;
List<Collection> dbs = new List<Collection>();
String adminApiKey;
List<String> userApiKeys = new List<String>();

Future main() async {
  bool running = true;
  await startSequence();
  var server = await HttpServer.bind(
    ipAddress,
    port,
  );
  print(" * Server bound to $ip:${server.port}");

  while (running) {
    await for (HttpRequest request in server) {
      String m = request.method;
      if (requestTypes[m]) {
        print("\n");
        handleRequest(request, request.headers.value("x-api-key"));
      } else {
        print(
            "Error: Invalid Request from ${request.connectionInfo.remoteAddress}");
        request.response.write("Error: Invalid request");
      }
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

void handleRequest(HttpRequest r, key) {
  String m = r.method;
  try {
    switch (m) {
      case ("GET"):
        {
          handleGet(r, key);
        }
        break;
      case ("POST"):
        {
          handlePost(r, key);
        }
        break;
      case ("PUT"):
        {
          handlePost(r, key);
        }
        break;
      case ("HEAD"):
        {
          handlePost(r, key);
        }
        break;
      case ("DELETE"):
        {
          handleDelete(r, key);
        }
        break;
      case ("PATCH"):
        {
          handlePost(r, key);
        }
        break;
      case ("OPTIONS"):
        {
          handlePost(r, key);
        }
        break;
    }
  } catch (e) {
    print("Exception in handleRequest: $e");
  }
}

void handleGet(HttpRequest r, String key) {
  String path = r.uri.path;
  switch (path) {
    case ("/query"):
      {
        if(checkApiKey(key, requestPermissions['QObject'])){
          print("Object query from ${r.connectionInfo.remoteAddress}");
          String query = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
        
      }
      break;
    case ("/query/attribute"):
      {
        if(checkApiKey(key, requestPermissions['QAttribute'])){
          print("Attribute query from ${r.connectionInfo.remoteAddress}");
          String query = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    case ("/query/allAttributes"):
      {
        if(checkApiKey(key, requestPermissions['QAllAttributes'])){
          print("All attributes query from ${r.connectionInfo.remoteAddress}");
          String query = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
        
      }
      break;
    case ("/query/byAttribute"):
      {
        if(checkApiKey(key, requestPermissions['QByAttribute'])){
          print("By attributes query from ${r.connectionInfo.remoteAddress}");
          String query = r.uri.queryParameters["db"];
          String att = r.uri.queryParameters["a"];
          Future<String> content = utf8.decodeStream(r);
          if(content == null){
            String end = "Request body is empty";
                print(end);
                var response = r.response;
                response.write(end);
                response.close();
          }else{
            content.then((result) {
            dynamic v;
            AttributeContainer ac;
            Map<String, dynamic> attData = json.decode(result);
            ac = new AttributeContainer(att, attData[att]);
            v = ac.value;
            print("Queried $query for attribute $att of value $v");
            if (dbs.any((Collection value) => value.name == query)) {
              Collection c = dbs.singleWhere((col) => col.name == query);
              List<Map<String, dynamic>> objs = new List<Map<String, dynamic>>();
              try { 
                c.dataList.forEach((DataObject d) {
                  if (d.data.containsKey(att)) {
                    AttributeContainer atc = new AttributeContainer(att, d.data[att]);
                    if(atc.value == v){
                      d.data.forEach((k,v) => print(k.toString() + " " + v.toString()));
                      objs.add({d.id.toString(): d.data});
                    }
                  }
                });
                String end = json.encode(objs);
                var response = r.response;
                response.write(end);
                response.close();
                print(end);  
              } catch (e) {
                String end = "Could not parse AttributeContainer from request body\n" + e;
                print(end);
                var response = r.response;
                response.write(end);
                response.close();
              }
            }else {
              String end = "Could not find collection $query";
              print(end);
              var response = r.response;
              response.write(end);
              response.close();
            }
                  
            }); 
        }
      }else {
        String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
        print(end);
        var response = r.response;
        response.write(end);
        response.close();
      }
        
    }
      break;
    case ("/query/newId"):
      {
        if(checkApiKey(key, requestPermissions['QnewId'])){
          print("New ID query from ${r.connectionInfo.remoteAddress}");
          String query = r.uri.queryParameters["db"];
          print("Queried $query for new id ");
          if (dbs.any((Collection value) => value.name == query)) {
            Collection c = dbs.singleWhere((col) => col.name == query);
            int maxId = 0;
            c.dataList.forEach((DataObject d) {
              if (d.id > maxId) {
                maxId = d.id;
              }
            });
            String end = json.encode(maxId + 1);
            print(end);
            var response = r.response;
            response.write(end);
            response.close();
          } else {
            String end = "Could not find collection $query";
            print(end);
            var response = r.response;
            response.write(end);
            response.close();
          }
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    default:
      {
        print(
            "Error: Invalid GET request from ${r.connectionInfo.remoteAddress}");
        print("\'" + r.uri.path + "\'");
        var response = r.response;
        response.write("Error: Invalid request");
        response.close();
      }
  }
}

void handlePost(HttpRequest r, String key) {
  String path = r.uri.path;
  switch (path) {
    case ("/add"):
      {
        if(checkApiKey(key, requestPermissions['AEmpty'])){
          print(
              "Empty object add request from ${r.connectionInfo.remoteAddress}");
          String add = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    case ("/add/object"):
      {
        if(checkApiKey(key, requestPermissions['AObject'])){
          print("Object add request from ${r.connectionInfo.remoteAddress}");
          String add = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    case ("/add/attribute"):
      {
        if(checkApiKey(key, requestPermissions['AAttribute'])){
          print("Attribute add request from ${r.connectionInfo.remoteAddress}");
          String add = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    case ("/mod/object"):
      {
        if(checkApiKey(key, requestPermissions['MObject'])){
          print("Object mod request from ${r.connectionInfo.remoteAddress}");
          String mod = r.uri.queryParameters["db"];
          int id = int.parse(r.uri.queryParameters["id"]);
          int newId = int.parse(r.uri.queryParameters["n"]);
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    case ("/mod/attribute"):
      {
        if(checkApiKey(key, requestPermissions['MAttribute'])){
          print("Attribute mod request from ${r.connectionInfo.remoteAddress}");
          String mod = r.uri.queryParameters["db"];
          int id = int.parse(r.uri.queryParameters["id"]);
          String att = r.uri.queryParameters["a"];
          Future<String> content = utf8.decodeStream(r);
          if (dbs.any((Collection value) => value.name == mod)) {
            Collection c = dbs.singleWhere((col) => col.name == mod);
            content.then((result) {
              AttributeContainer attribute =
                  new AttributeContainer(att, json.decode(result)[att]);
              print(
                  "Requested to modify attribute ${attribute.key} of $id in $mod");
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    default:
      {
        print(
            "Error: Invalid POST request from ${r.connectionInfo.remoteAddress}");
        print("\'" + r.uri.path + "\'");
        var response = r.response;
        response.write("Error: Invalid request");
        response.close();
      }
  }
}

void handlePut(HttpRequest r, String key) {}
void handleHead(HttpRequest r, String key) {}
void handleDelete(HttpRequest r, String key) {
  String path = r.uri.path;
  switch (path) {
    case ("/delete/object"):
      {
        if(checkApiKey(key, requestPermissions['DObject'])){
          print("Object delete request from ${r.connectionInfo.remoteAddress}");
          String del = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    case ("/delete/attribute"):
      {
        if(checkApiKey(key, requestPermissions['DAttribute'])){
          print(
              "Attribute delete request from ${r.connectionInfo.remoteAddress}");
          String del = r.uri.queryParameters["db"];
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
        }else {
          String end = "Error: Unauthorized Request from ${r.connectionInfo.remoteAddress}";
          print(end);
          var response = r.response;
          response.write(end);
          response.close();
        }
      }
      break;
    default:
      {
        print(
            "Error: Invalid DELETE request from ${r.connectionInfo.remoteAddress}");
        print("\'" + r.uri.path + "\'");
        var response = r.response;
        response.write("Error: Invalid request");
        response.close();
      }
  }
}

void handlePatch(HttpRequest r, String key) {}
void handleOptions(HttpRequest r, String key) {
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
  Map<String, dynamic> requestPerms = map["Permissions"];
  requestPerms.forEach((k,v) => requestPermissions[k] = parsePermission(v));
}

bool parsePermission(String perm){
  if(perm == "admin"){
    return true;
  }else{
    return false;
  }
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

Future<bool> GenerateAdminApiKey() async {
  File file = new File("data.jserv");
  List<String> data = await file.readAsLines();
  if (data.length > 1) {
    if (data.elementAt(1) == "new" || data.elementAt(1) == "" || data.elementAt(1) == " ") {
      Random rand = new Random();
      var keyList = new List.generate(32, (index) {
        return rand.nextInt(26) + 97;
      });
      String keyString = new String.fromCharCodes(keyList);
      file.writeAsString(data.elementAt(0) + "\n" + keyString);
      return true;
    }else{
      return true;
    }
  }else{
    print("Failed to detect API Key. Type \"new\" on the second line of data.jserv to generate an Admin API Key.");
    return false;
  }
}

Future<void> GenerateUserApiKeys() async {
  File file = new File("keys.jserv");
  String wholeFile = await file.readAsString();
  List<String> data = await file.readAsLines();
  if(wholeFile.isEmpty){
    Random rand = new Random();
    var keyList = new List.generate(32, (index) {
      return rand.nextInt(26) + 97;
    });
    String keyString = new String.fromCharCodes(keyList);
    file.writeAsString(keyString);
  }else{
    if (data.length >= 0) {
      if (data.elementAt(0).isEmpty || data.elementAt(0) == "") {
        Random rand = new Random();
        var keyList = new List.generate(32, (index) {
          return rand.nextInt(26) + 97;
        });
        String keyString = new String.fromCharCodes(keyList);
        file.writeAsString(data.elementAt(0) + keyString);
      }
    }
  }
  
}

Future<void> startSequence() async {
  print(" * Starting...");
  if(await GenerateAdminApiKey() == false){
    print("API Key failure. \nPress enter to exit.");
    stdin.readLineSync();
    exit(0);
  }else{
    await GenerateUserApiKeys();
    String version = "0.1.1";
    List<String> data = await File("data.jserv").readAsLines();
    String implementation = data.elementAt(0);
    adminApiKey = data.elementAt(1);
    List<String> data2 = await File("keys.jserv").readAsLines();
    for(String s in data2){
      print(s);
      userApiKeys.add(s);
    }
    userApiKeys.add(data2.elementAt(0));
    print(" * jServ v$version implemented for $implementation");
    print(
        " * Admin API Key for this instance of jServ is $adminApiKey. Please put this key in the headers of your requests");
    print(" -----------------------");
    await readDatabases();
    print(" * Loading databases...");
    await readConfig();
    print(" * Loading config...");
    await createIP();
    print(" * Binding server...");
    print(" * Done!");
  }
}
bool checkApiKey(String key, bool requireAdmin){
  if(requireAdmin == true){
    if (key == adminApiKey) {
      return true;
    }else{
      return false;
    }
  }else{
    if(userApiKeys.contains(key) || key == adminApiKey){
      return true;
    }else{
      return false;
    }
  }
    
}

  

