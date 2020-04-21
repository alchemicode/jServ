import 'dart:io';
import 'dart:convert';
import 'DataObject.dart';
import 'Collection.dart';

Map<String, dynamic> requestTypes = {"GET":true,"POST":true,"PUT":true,"HEAD":false,"DELETE":false,"PATCH":false,"OPTIONS":false};
String ip = "localhost";
int port = 4040;
List<Collection> dbs;

Future main() async {
  
  dbs = new List<Collection>();
  dbs.add(new Collection("db"));
  
  await readConfig();
  
  var server = await HttpServer.bind(
    InternetAddress.loopbackIPv4,
    port,
  );
  var obj = new DataObject.withMap(10,{"e": 1, "yee":"ooo"});
  print("Listening on localhost:${server.port}");

  await for (HttpRequest request in server) {
    String m = request.method;
    if(requestTypes[m]){
      handleRequest(request);
    }else{
      request.response.write("Error: Invalid request");
    }
    
  }

  
}
void handleRequest(HttpRequest r){
    String m = r.method;
    try{
      switch(m){
        case("GET"):{
          handleGet(r);
        }
        break;
        case("POST"):{
          handlePost(r);
        }
        break;
        case("PUT"):{
          handlePost(r);
        }
        break;
        case("HEAD"):{
          handlePost(r);
        }
        break;
        case("DELETE"):{
          handlePost(r);
        }
        break;
        case("PATCH"):{
          handlePost(r);
        }
        break;
        case("OPTIONS"):{
          handlePost(r);
        }
        break;
      }
    }catch(e){
      print("Exception in handleRequest: $e");
    }
}
void handleGet(HttpRequest r){
  String end = "";
  String path = r.uri.path;
  
  if(path == "/query"){
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    if(dbs.any((Collection value) => value.name == query)){
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data = c.dataList.singleWhere((d) => d.id == id);
      end = data.toString();
    }
  }

  if(path == "/query/attribute"){
    String query = r.uri.queryParameters["q"];
    int id = int.parse(r.uri.queryParameters["id"]);
    String att = r.uri.queryParameters["a"];
    if(dbs.any((Collection value) => value.name == query)){
      Collection c = dbs.singleWhere((col) => col.name == query);
      DataObject data = c.dataList.singleWhere((d) => d.id == id);
      end = data.data[att];
    }
  }

  var response = r.response;
  response.write(end);
  response.close();

}
void handlePost(HttpRequest r){
  
}
void handlePut(HttpRequest r){
  
}
void handleHead(HttpRequest r){
  
}
void handleDelete(HttpRequest r){
  String m = r.method;
  r.response.write("this is a $m");
}
void handlePatch(HttpRequest r){
  
}
void handleOptions(HttpRequest r){
  String m = r.method;
  r.response.write("this is a $m");
}

Future readConfig() async {
  File file = new File("config.json");
  String content;

  if(await file.exists()){
    content = await file.readAsString();
  }
  var map = json.decode(content);
  ip = map["ip"];
  port = map["port"];
  requestTypes = map["Requests"];
}





